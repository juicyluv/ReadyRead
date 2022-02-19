CREATE TABLE IF NOT EXISTS books(
    id bigserial primary key,
    title text not null,
    description text not null,
    year smallint,
    price decimal(10,2),
    page_count smallint,
    count int not null default 1,
    author_id bigint,
    genre_id smallint,
    language_id smallint,

    foreign key(author_id) references authors(id) on delete cascade,
    foreign key(genre_id) references genres(id) on delete cascade,
    foreign key(language_id) references languages(id) on delete cascade
);