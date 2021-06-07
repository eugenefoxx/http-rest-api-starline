CREATE TABLE users
(
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null,
    firstname varchar not null,
    lastname varchar not null,
    role varchar null,
    groups varchar null,
    tabel varchar null
);