# http-rest-api-starline

https://codesource.io/build-a-crud-application-in-golang-with-postgresql/

git add .
git commit -a -m "v1.1.979 modify style"
git push
git tag v1.1.979
git push -q origin v1.1.979

pgweb --host localhost --user postgres --db starline
/\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***\*\*\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***/
deploy systemd
vim /etc/systemd/system/appstarlineprod.service
[Unit]
Description=App

[Service]
ExecStart=/home/eugenearch/Code/github.com/eugenefoxx/apiserver
WorkingDirectory=/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/

[Install]
WantedBy=multi-user.target

systemctl status appstarlineprod.service
systemctl start appstarlineprod.service
systemctl start appstarlineprod.service
/\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***\*\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***/

/\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***\*\*\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***
// удаление дублирующихся ЕО, до присвоения статус OK or NG
delete from transfer a using transfer b where a.id < b.id and a.idmaterial = b.idmaterial
and a.status is null;
//\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***\*\*\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***

/\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***\*\*\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***
// создать пакет utils
https://www.youtube.com/watch?v=mkFkWTuDIVU
func ToJSON(w http.ResponseWriter, data interface{}, statusCode int) {
w.Header().Set("Content-Type", "application/json; charset-UTF8")
w.WriteHeader(statusCode)
err := json.NewEncoder(w).Encode(data)
CheckError(err)
}

func CheckError(err error) {

    if err != nil {
        log.Fatal(err)

}
}
//\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***\*\*\***\*\*\*\*\*\*\***\*\*\*\*\***\*\*\*\*\*\*\***

P2013002LK340621858R1000425231Q03000D00000000
P2013002LK340621858R1000425299Q03000D00000000

select \* from shipmentbysap where lastname = 'Коновалов' and shipment_date between '2020-02-15' and '2020-05-01';
select material, qty, to_char (shipment_date, 'YYYY-MON-DD') shipment_date2, to_char (shipment_time, 'HH24:MI:SS') shipment_time2, lastname from shipmentbysap shipment_date;select material, qty, to_char (shipment_date, 'YYYY-MON-DD') shipment_date2, to_char (shipment_time, 'HH24:MI:SS') shipment_time2, lastname from shipmentbysap shipment_date;

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

// выборка для представления
select id_return.material, sap_stock_import.material_description, id_return.id_roll, id_return.lot, id_sap_import.spp_element
from id_return  
left outer join id_sap_import on (id_return.id_roll = id_sap_import.id)
left outer join sap_stock_import on (id_return.material = sap_stock_import.material);

select id_return.id_roll, id_return.material, sap_description.material_description, id_return.lot, id_sap_import.spp_element
from id_return  
left outer join id_sap_import on (id_return.id_roll = id_sap_import.id)
left outer join sap_description on (id_return.material = sap_description.material);

создать представление для выборки

CREATE VIEW sap_description AS
SELECT DISTINCT sap_stock_import.material,
sap_stock_import.material_description
FROM sap_stock_import
ORDER BY sap_stock_import.material;

CREATE VIEW sap_description AS
SELECT DISTINCT sap_stock_import.material,
sap_stock_import.material_description
FROM sap_stock_import
ORDER BY sap_stock_import.material;

select id_return.id_roll, id_return.material, sap_description.material_description, id_return.lot, id_sap_import.spp_element
from id_return  
left outer join id_sap_import on (id_return.id_roll = id_sap_import.id)
left outer join sap_description on (id_return.material = sap_description.material);

DROP VIEW sap_description;

CREATE VIEW stock AS
select sap_stock_import.material, sap_stock_import.warehouse, sap_stock_import.qty_free
from sap_stock_import
where sap_stock_import.warehouse='7813';

// создаем сток по итоговой сумме на 7813
create view stock_sum as
select material, sum(qty_free)
from stock
group by material;

// общий код
CREATE VIEW sap_description AS
SELECT DISTINCT sap_stock_import.material,
sap_stock_import.material_description
FROM sap_stock_import
ORDER BY sap_stock_import.material;

CREATE VIEW stock7813 AS
select sap_stock_import.material, sap_stock_import.warehouse, sap_stock_import.qty_free
from sap_stock_import
where sap_stock_import.warehouse='7813';

create view stock_sum7813 as
select material, sum(qty_free)
from stock7813
group by material;

select id_return.id_roll, id_return.material, sap_description.material_description, id_return.lot,
id_return.qty_fact, id_return.qty_sap, id_return.qty_panacim,
id_sap_import.spp_element, stock_sum7813.sum, id_sap_import.warehouse,
panacim_stock.estimated_current_quantity, panacim_stock.storage_unit, panacim_stock.storage_sub_unit, TO_CHAR(id_return.shipment_date, 'YYYY-MM-DD') shipment_date2,
TO_CHAR(id_return.shipment_time, 'HH24:MI:SS') shipment_time2, id_return.lastname
from id_return
left outer join id_sap_import on (id_return.id_roll = id_sap_import.id)
left outer join sap_description on (id_return.material = sap_description.material)
left outer join stock_sum7813 on(id_return.material = stock_sum7813.material)
left outer join panacim_stock on(id_return.id_roll = panacim_stock.material_barcode);

DROP VIEW sap_description;
DROP VIEW stock_sum7813;
DROP VIEW stock7813;

const btn = document.querySelector('sendupdate');

function sendupdate(data) {
const form = document.querySelector('form[name="valform"]'),
const status = form.elements['status'].value,
const note = form.elements['note'].value,

}

btn.addEventListener( 'click', function() {

    sendupdate()

// let elements = document
//sendData( {test:'ok'} );
})

console.log(status);
