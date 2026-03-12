package relationalpostgres

import (
	"context"
	"fmt"

	"github.com/BHRK-codelabs/corekit/configkit"
	"github.com/BHRK-codelabs/relationalkit"
)

type Module struct {
	cfg   *configkit.Config
	store *Store
}

func NewModule(cfg *configkit.Config) (*Module, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	store, err := Open(cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	return &Module{cfg: cfg, store: store}, nil
}

func (m *Module) Name() string {
	return "postgres-relational"
}

func (m *Module) Store() *Store {
	return m.store
}

func (m *Module) Start(ctx context.Context) error {
	if err := m.store.Ping(ctx); err != nil {
		return err
	}
	var err error
	_, err = m.store.Exec(ctx, relationalkit.Command{
		Text: `CREATE SCHEMA IF NOT EXISTS platform`,
	})
	if err != nil {
		return err
	}
	_, err = m.store.Exec(ctx, relationalkit.Command{
		Text: `CREATE TABLE IF NOT EXISTS platform.file_mananger (
			id UUID PRIMARY KEY,
			filename TEXT NOT NULL,
			kind TEXT NOT NULL,
			download_linkl TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			mime_type TEXT NOT NULL,
			size INTEGER NOT NULL,
			tenant_id UUID NOT NULL
		)`,
	})
	return err
}

func (m *Module) Stop(context.Context) error {
	return m.store.Close()
}
