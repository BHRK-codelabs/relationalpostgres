package relationalpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/BHRK-codelabs/relationalkit"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

func Open(databaseURL string) (*Store, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database url is required")
	}
	db, err := sql.Open("pgx", normalizeDatabaseURL(databaseURL))
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Query(ctx context.Context, query relationalkit.Query) ([]relationalkit.Row, error) {
	rows, err := s.db.QueryContext(ctx, query.Text, query.Args...)
	if err != nil {
		return nil, relationalkit.NewUnavailableError("postgres query failed", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

func (s *Store) Exec(ctx context.Context, command relationalkit.Command) (relationalkit.Result, error) {
	result, err := s.db.ExecContext(ctx, command.Text, command.Args...)
	if err != nil {
		return relationalkit.Result{}, relationalkit.NewUnavailableError("postgres exec failed", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return relationalkit.Result{RowsAffected: rowsAffected}, nil
}

func (s *Store) Begin(ctx context.Context) (relationalkit.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, relationalkit.NewUnavailableError("postgres begin failed", err)
	}
	return &Tx{tx: tx}, nil
}

type Tx struct {
	tx *sql.Tx
}

func (t *Tx) Query(ctx context.Context, query relationalkit.Query) ([]relationalkit.Row, error) {
	rows, err := t.tx.QueryContext(ctx, query.Text, query.Args...)
	if err != nil {
		return nil, relationalkit.NewUnavailableError("postgres tx query failed", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

func (t *Tx) Exec(ctx context.Context, command relationalkit.Command) (relationalkit.Result, error) {
	result, err := t.tx.ExecContext(ctx, command.Text, command.Args...)
	if err != nil {
		return relationalkit.Result{}, relationalkit.NewUnavailableError("postgres tx exec failed", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return relationalkit.Result{RowsAffected: rowsAffected}, nil
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

func scanRows(rows *sql.Rows) ([]relationalkit.Row, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := make([]relationalkit.Row, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		row := relationalkit.Row{}
		for i, column := range columns {
			switch v := values[i].(type) {
			case []byte:
				row[column] = string(v)
			default:
				row[column] = v
			}
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func normalizeDatabaseURL(databaseURL string) string {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return databaseURL
	}

	query := parsed.Query()
	if query.Get("default_query_exec_mode") == "" {
		query.Set("default_query_exec_mode", "simple_protocol")
	}
	parsed.RawQuery = query.Encode()

	normalized := parsed.String()
	if strings.TrimSpace(normalized) == "" {
		return databaseURL
	}
	return normalized
}
