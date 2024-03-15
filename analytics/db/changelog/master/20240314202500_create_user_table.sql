-- +goose Up
create table if not exists "user"
(
    id         bigserial primary key,
    public_id  uuid                                      not null,

    name       varchar                                   not null,

    created_at timestamp without time zone default now() not null
);

create unique index on "user" (public_id);

-- +goose Down
