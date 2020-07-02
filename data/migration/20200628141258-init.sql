
-- +migrate Up
create function gen_random_uuid() returns uuid
    language c
    as '$libdir/pgcrypto', 'pg_random_uuid';

create table users (
	id uuid not null unique primary key default gen_random_uuid(),
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	name varchar(255),
	email varchar(255) unique not null,
	avatar text
);

create table conferences (
	id uuid not null unique primary key default gen_random_uuid(),
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	topic varchar(255) not null,
	description text,
	host_id uuid references users (id) not null,
	is_active boolean not null,
	password text
);

create table conference_users (
	user_id uuid references users (id) on delete cascade not null,
	conference_id uuid references conferences (id) on delete cascade not null,
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	primary key (user_id, conference_id)
);

-- +migrate Down
drop table conference_users;
drop table conferences;
drop table users;
drop function gen_random_uuid;