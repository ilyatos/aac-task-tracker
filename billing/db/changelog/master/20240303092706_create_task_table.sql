-- +goose Up
create table if not exists task
(
    id             bigserial primary key,
    public_id      uuid                                      not null,
    user_public_id uuid                                      not null,

    description    varchar                                   not null,
    credit_price   integer check (credit_price >= 0)         not null,
    debit_price    integer check (debit_price >= 0)          not null,

    created_at     timestamp without time zone default now() not null
);

create unique index on task (public_id);

-- +goose Down