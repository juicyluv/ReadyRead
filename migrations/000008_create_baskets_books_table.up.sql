CREATE TABLE IF NOT EXISTS baskets_books(
    basket_id bigint not null,
    book_id bigint not null,
    count int not null default 0,

    foreign key(basket_id) references baskets(id) on delete cascade,
    foreign key(book_id) references books(id) on delete cascade
);