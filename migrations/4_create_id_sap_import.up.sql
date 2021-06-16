CREATE TABLE id_sap_import
(
    id integer NOT NULL,
    package integer NOT NULL,
    plant varchar NULL,
    warehouse varchar NULL,
    material varchar NULL,
    qty real NOT NULL,
    storage_unit varchar NULL,
    lot varchar NOT NULL,
    t1 varchar NULL,
    t2 varchar NULL,
    spp_element varchar NULL
);