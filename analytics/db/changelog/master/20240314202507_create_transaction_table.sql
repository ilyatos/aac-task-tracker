-- +goose Up
create type transaction_type as enum ('task_assigned', 'task_completed', 'payout');

create table if not exists transaction
(
    id             bigserial primary key,
    public_id      uuid                                      not null,
    user_public_id uuid                                      not null,

    type           transaction_type                          not null,
    credit         int                         default 0     not null,
    debit          int                         default 0     not null,

    created_at     timestamp without time zone default now() not null
);

create unique index on transaction (public_id);

-- +goose Down
