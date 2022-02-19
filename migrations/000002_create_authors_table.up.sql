CREATE TABLE IF NOT EXISTS authors(
    id bigserial primary key,
    name text not null,
    surname text not null
);