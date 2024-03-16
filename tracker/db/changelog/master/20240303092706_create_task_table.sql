-- +goose Up
create type task_status as enum ('new', 'done');

create table if not exists task
(
    id             bigserial primary key,
    public_id      uuid                                   not null,
    user_public_id uuid                                   not null,
    --
    status         task_status                            not null,
    description    varchar                                not null,
    --
    created_at     timestamp with time zone default now() not null,
    updated_at     timestamp with time zone default now() not null
);

create unique index on task (public_id);

-- +goose Down