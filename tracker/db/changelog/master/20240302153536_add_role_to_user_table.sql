-- +goose Up
create type "user_role" as enum ('employee', 'manager', 'accountant', 'admin');
alter table "user"
    add column role user_role not null;

-- +goose Down
