package queries

const CreateMigrationQuery = `
insert into _dbmap_migrations (migration_name) values ($1);
`
