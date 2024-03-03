-- +goose Up
create extension if not exists pgcrypto;
create extension if not exists tsm_system_rows;

create table if not exists "user"
(
    id         bigserial primary key,
    public_id  uuid                                   not null,
    --
    name       varchar                                not null,
    email      varchar                                not null,
    --
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
);

create unique index on "user" (public_id);

-- +goose Down
