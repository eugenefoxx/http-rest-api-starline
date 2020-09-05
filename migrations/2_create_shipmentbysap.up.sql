CREATE TABLE shipmentbysap (
    material integer NOT NULL,
    qty integer NOT NULL,
    shipment_date DATE DEFAULT CURRENT_DATE,
    shipment_time TIME DEFAULT CURRENT_TIME,
    id bigint NOT NULL,
    lastname varchar NOT NULL,
    comment varchar NOT NULL
);