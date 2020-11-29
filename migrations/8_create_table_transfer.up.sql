CREATE TABLE transfer
(
    id bigserial not null primary key,
    idmaterial VARCHAR NOT NULL,
    sap INTEGER NOT NULL,
    lot VARCHAR NOT NULL,
    idroll INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    productiondate VARCHAR NOT NULL,
    numbervendor VARCHAR NOT NULL,
    vendor VARCHAR NULL,
    note text NULL,
    location text NOT NULL,
    status VARCHAR NULL,
    date DATE DEFAULT CURRENT_DATE,
    time TIME DEFAULT CURRENT_TIME,
    lastname VARCHAR NOT NULL,
    update VARCHAR NULL,
    dateupdate DATE
    DEFAULT CURRENT_DATE,
    timeupdate TIME DEFAULT CURRENT_TIME,
    lastnameaccept VARCHAR NULL,
    dateaccept DATE DEFAULT CURRENT_DATE,
    timeaccept TIME DEFAULT CURRENT_TIME
);