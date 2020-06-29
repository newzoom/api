
-- +migrate Up
create table users (
	id serial not null unique primary key,
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	name varchar(200),
	email varchar(200)
);

create table conferences (
	id serial not null unique primary key,
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	name varchar(200)
);

create table conference_users (
	user_id integer references users (id) on delete cascade not null,
	conference_id integer references conferences (id) on delete cascade not null,
	created_at timestamp(6) with time zone default now(),
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	primary key (user_id, conference_id)
);

-- +migrate Down
drop table conference_users;
drop table conferences;
drop table users;