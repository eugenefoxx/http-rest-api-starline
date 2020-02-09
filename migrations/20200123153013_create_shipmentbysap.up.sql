CREATE TABLE shipmentbysap (
    material integer not null,
    qty integer not null,
    shipment_date date not null default current_date,
    shipment_time time not null default localtime,
    id bigint not null,
    lastname varchar not null,
    comment varchar null
);