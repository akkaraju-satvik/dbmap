package queries

const CreateMigrationQuery = `
insert into _pgo_migrations (migration_name) values ($1);
`
