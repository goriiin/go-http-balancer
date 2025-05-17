
create table if not exists some_data (
    id int generated always as identity primary key ,
    name varchar(255) not null,
    description TEXT,

    created_at TIMESTAMPTZ not null default NOW()
);


CREATE INDEX idx_some_data_name ON some_data(id);