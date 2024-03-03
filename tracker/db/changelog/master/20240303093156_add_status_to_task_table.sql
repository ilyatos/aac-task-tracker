-- +goose Up
create type task_status as enum ('new', 'done');
alter table task
    add column status task_status not null;

-- +goose Down
