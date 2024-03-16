-- +goose Up
create table if not exists billing_cycle
(
    id             bigserial primary key,
    user_public_id uuid                                      not null,

    start_at       timestamp without time zone               not null,
    end_at         timestamp without time zone               not null,
    is_closed      boolean                     default false not null,

    created_at     timestamp without time zone default now() not null
);

create unique index on billing_cycle (user_public_id, start_at, end_at);

-- +goose Down
