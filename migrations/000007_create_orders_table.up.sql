CREATE TABLE IF NOT EXISTS orders(
    id bigint primary key,
    date timestamptz not null default now(),
    total_price decimal(10,2),
    user_id bigint not null,
    basket_id bigint not null,

    foreign key(user_id) references users(id) on delete cascade,
    foreign key(basket_id) references baskets(id) on delete cascade
);