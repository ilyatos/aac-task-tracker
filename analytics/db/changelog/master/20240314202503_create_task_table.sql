-- +goose Up
create table if not exists task
(
    id           bigserial primary key,
    public_id    uuid                                      not null,

    description  varchar                                   not null,
    credit_price integer                                   null,
    debit_price  integer                                   null,

    completed_at timestamp without time zone default null,
    created_at   timestamp without time zone default now() not null
);

create unique index on task (public_id);

-- +goose Down
