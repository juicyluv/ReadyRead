CREATE TABLE IF NOT EXISTS users(
    id bigserial primary key,
    username text not null,
    email text not null,
    password text not null,
    verified bool not null default false,
    address text,
    phone_number text
);