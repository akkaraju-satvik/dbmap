package queries

const InitQuery = `
drop type if exists _dbmigo_migration_status cascade;
create type _dbmigo_migration_status as enum ('APPLIED', 'PENDING');

create table if not exists _dbmigo_migrations (
  migration_id uuid primary key default gen_random_uuid(),
  migration_name text not null,
  migration_time timestamp with time zone not null default now(),
  migration_status _dbmigo_migration_status not null default 'PENDING'  
);

create table if not exists _dbmigo_migration_queries (
  migration_query_id uuid primary key default gen_random_uuid(),
  migration_id uuid not null,
  migration_query text not null,
  query_time timestamp with time zone not null default now()
);

alter table _dbmigo_migration_queries drop constraint if exists fk_migration_id;
alter table _dbmigo_migration_queries add constraint fk_migration_id foreign key (migration_id) references _dbmigo_migrations(migration_id);

`
