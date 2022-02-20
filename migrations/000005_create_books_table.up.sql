CREATE TABLE IF NOT EXISTS books(
    id bigserial primary key,
    title text not null,
    description text not null,
    year smallint,
    price decimal(10,2) not null,
    page_count smallint,
    count int not null,
    author_id bigint not null,
    genre_id smallint not null,
    language_id smallint not null,

    foreign key(author_id) references authors(id) on delete cascade,
    foreign key(genre_id) references genres(id) on delete cascade,
    foreign key(language_id) references languages(id) on delete cascade
);