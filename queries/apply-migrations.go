package queries

const GetMigrations = `
select * from _dbmap_migrations where migration_status = 'PENDING'::_dbmap_migration_status order by migration_time asc;
`

const UpdateMigrationStatus = `
update _dbmap_migrations set migration_status = 'APPLIED'::_dbmap_migration_status where migration_id = $1;
`
const InsertMigrationQuery = `
insert into _dbmap_migration_queries (migration_id, migration_query, migration_type) values ($1, $2, 'UP'::_dbmap_migration_type);
`
