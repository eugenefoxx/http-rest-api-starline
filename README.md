# http-rest-api-starline
pgweb --host localhost --user postgres --db starline

select * from shipmentbysap where lastname = 'Коновалов' and shipment_date between '2020-02-15' and '2020-05-01';


migrate create -ext sql -dir migrations create_users 
migrate create -ext sql -dir migrations create_shipmentbysap

migrate -path migrations -database "postgres://localhost/starline?sslmode=disable" up
migrate -path migrations -database "postgres://localhost/starline?sslmode=disable" down


CREATE TABLE shipmentbysap (
    material integer not null,
    qty integer not null,
    shipment_date date not null default current_date,
    shipment_time time not null default current_time,
    id bigint not null,
    lastname varchar not null,
    comment varchar null
);

DATE DEFAULT CURRENT_DATE

CREATE TABLE shipmentbysap (
    material integer NOT NULL,
    qty integer NOT NULL,
    shipment_date DATE DEFAULT CURRENT_DATE,
    shipment_time TIME DEFAULT CURRENT_TIME,
    id bigint NOT NULL,
    lastname varchar NOT NULL,
    comment varchar NOT NULL
);

DROP TABLE shipmentbysap;



CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null,
    firstname varchar not null,
    lastname varchar not null
);

DROP TABLE users;
