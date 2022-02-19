CREATE TABLE IF NOT EXISTS baskets(
    id bigserial primary key,
    user_id bigint not null,

    foreign key(user_id) references users(id) on delete cascade
);