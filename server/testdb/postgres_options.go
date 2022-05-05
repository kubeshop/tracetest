package testdb

type PostgresOption func(*postgresDB) error

func WithMigrations(migrationFolder string) PostgresOption {
	return func(pd *postgresDB) error {
		pd.migrationsFolder = migrationFolder

		return nil
	}
}
