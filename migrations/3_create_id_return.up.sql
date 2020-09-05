CREATE TABLE id_return (
    material integer NOT NULL,
    id_roll integer NOT NULL primary key,
    lot varchar NOT NULL,
    qty_fact integer NOT NULL,
    qty_sap integer NOT NULL,
    qty_panacim integer NOT NULL,
    shipment_date DATE DEFAULT CURRENT_DATE,
    shipment_time TIME DEFAULT CURRENT_TIME,
    id bigint NOT NULL,
    lastname varchar NOT NULL
);