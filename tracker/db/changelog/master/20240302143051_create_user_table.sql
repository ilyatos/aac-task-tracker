-- +goose Up
create type "user_role" as enum ('employee', 'manager', 'accountant', 'admin');

create table if not exists "user"
(
    id         bigserial primary key,
    public_id  uuid                                   not null,
    --
    role       user_role                              not null,
    name       varchar                                not null,
    email      varchar                                not null,
    --
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
);

create unique index on "user" (public_id);

-- +goose Down
