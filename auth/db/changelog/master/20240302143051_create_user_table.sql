-- +goose Up
create extension if not exists pgcrypto;

create table if not exists "user"
(
    id         bigserial primary key,
    public_id  uuid                     default gen_random_uuid() not null,
    --
    name       varchar                                            not null,
    email      varchar                                            not null,
    password   varchar                                            not null,
    --
    created_at timestamp with time zone default now()             not null,
    updated_at timestamp with time zone default now()             not null
);

create unique index on "user" (public_id);
create unique index on "user" (email);

-- +goose Down
