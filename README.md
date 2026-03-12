# relationalpostgres

`relationalpostgres` is the PostgreSQL adapter for `relationalkit`.

It provides:
- a `database/sql` backed `relationalkit` store
- a kernel module that opens the connection, validates it and creates the `orders` table used by `baseline-ms`

## Usage

Use `DATABASE_URL` in standard Postgres DSN form and wire the module from application bootstrap.

In `baseline-ms`, this adapter is selected with:

```powershell
$env:DEMO_REPOSITORY_DRIVER="external"
$env:DATABASE_URL="postgres://user:password@host:5432/dbname?sslmode=require"
```

For PgBouncer or Supabase poolers, this adapter automatically adds `default_query_exec_mode=simple_protocol` when the DSN does not define it.

## Notes

- this adapter is for real Postgres-backed persistence
- keep credentials in env vars for local development
- do not hardcode DSNs in code or docs intended for publication
