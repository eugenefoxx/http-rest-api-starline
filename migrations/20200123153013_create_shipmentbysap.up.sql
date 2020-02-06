CREATE TABLE shipmentbysap (
    material smallint not null primary key,
    qty integer not null,
    shipment_date date not null default current_date,
    id bigint not null,
    comment varchar not null
);