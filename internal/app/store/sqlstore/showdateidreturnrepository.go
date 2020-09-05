package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
)

// ShowdateidreturnRepository ...
type ShowdateidreturnRepository struct {
	store *Store
}

// ShowDataByDate ...
func (r *ShowdateidreturnRepository) ShowDataByDate(shipmentDate2 string, shipmentDate3 string) (*model.Showdateidreturns, error) {

	selectDate := `select id_return.id_roll, id_return.material, sap_description.material_description, id_return.lot,
		id_return.qty_fact, id_return.qty_sap, id_return.qty_panacim,
		Coalesce(id_sap_import.spp_element, '') spp, Coalesce (stock_sum7813.sum, 0) summa, Coalesce (id_sap_import.warehouse, '') warehouse, 
		Coalesce (id_sap_import.qty, 0) qty_id_sap, 
		Coalesce (panacim_stock.estimated_current_quantity, 0) estimatedcurrentquantity, 
		Coalesce (panacim_stock.storage_unit, '') storageunit, Coalesce (panacim_stock.storage_sub_unit, '') storagesubunit, 
		TO_CHAR(id_return.shipment_date, 'YYYY-MM-DD') shipment_date2,
		TO_CHAR(id_return.shipment_time, 'HH24:MI:SS') shipment_time2, id_return.lastname
		from id_return
		left outer join id_sap_import on (id_return.id_roll = id_sap_import.id)
		left outer join sap_description on (id_return.material = sap_description.material)
		left outer join stock_sum7813 on(id_return.material = stock_sum7813.material)
		left outer join panacim_stock on(id_return.id_roll = panacim_stock.material_barcode)

		WHERE id_return.shipment_date BETWEEN $1 AND $2;`

	sapDescriptionview := `CREATE VIEW sap_description AS
		SELECT DISTINCT sap_stock_import.material,
    	sap_stock_import.material_description
   		FROM sap_stock_import
	  	ORDER BY sap_stock_import.material;`

	stock7813view := `CREATE VIEW stock7813 AS
		select sap_stock_import.material, sap_stock_import.warehouse, sap_stock_import.qty_free 
		from sap_stock_import
		where sap_stock_import.warehouse='7813';`

	stockSum7813view := `create view stock_sum7813 as
		select material, sum(qty_free)
		from stock7813 
		group by material;`

	dropsapDescriptionview := `DROP VIEW sap_description;`
	dropstock7813view := `DROP VIEW stock_sum7813;`
	dropstockSum7813view := `DROP VIEW stock7813;`

	idreturn := model.Showdateidreturn{}
	idreturnList := make(model.Showdateidreturns, 0)

	_, errsapDescriptionview := r.store.db.Exec(sapDescriptionview)
	if errsapDescriptionview != nil {
		//fmt.Println("ошибка в select")
		panic(errsapDescriptionview)
		//	return nil, errsapDescriptionview
	}

	_, errstock7813view := r.store.db.Exec(stock7813view)
	if errstock7813view != nil {
		//fmt.Println("ошибка в select")
		panic(errstock7813view)
		//	return nil, errstock7813view
	}

	_, errstockSum7813view := r.store.db.Exec(stockSum7813view)
	if errstockSum7813view != nil {
		//fmt.Println("ошибка в select")
		panic(errstockSum7813view)
		//	return nil, errstockSum7813view
	}

	rows, err := r.store.db.Query(
		//	"SELECT id_return.id_roll, id_return.material, sap_description.material_description, id_return.lot, id_return.qty_fact, id_return.qty_sap, id_return.qty_panacim, id_sap_import.spp_element, stock_sum7813.sum, id_sap_import.warehouse, id_sap_import.qty, panacim_stock.estimated_current_quantity, panacim_stock.storage_unit, panacim_stock.storage_sub_unit, TO_CHAR(id_return.shipment_date, 'YYYY-MM-DD') shipment_date2, TO_CHAR(id_return.shipment_time, 'HH24:MI:SS') shipment_time2, id_return.lastname from id_return left outer join id_sap_import on (id_return.id_roll = id_sap_import.id) left outer join sap_description on (id_return.material = sap_description.material) left outer join stock_sum7813 on(id_return.material = stock_sum7813.material) left outer join panacim_stock on(id_return.id_roll = panacim_stock.material_barcode) WHERE id_return.shipment_date BETWEEN $1 AND $2",
		selectDate, shipmentDate2, shipmentDate3)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&idreturn.IDRoll,
			&idreturn.Material,
			&idreturn.MaterialDescription,
			&idreturn.Lot,
			&idreturn.QtyFact,
			&idreturn.QtySAP,
			&idreturn.QtyPanacim,
			&idreturn.SPP,
			&idreturn.Summa,
			&idreturn.Warehouse,
			&idreturn.QtyIDSAP,
			&idreturn.EstimatedCurrentQuantity,
			&idreturn.StorageUnit,
			&idreturn.StorageSubUnit,
			&idreturn.ShipmentDate2,
			&idreturn.ShipmentTime2,
			&idreturn.Lastname,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("ошибка в rows.Scan")
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		idreturnList = append(idreturnList, idreturn)

	}

	_, errdropsapDescriptionview := r.store.db.Exec(dropsapDescriptionview)
	if errdropsapDescriptionview != nil {
		//fmt.Println("ошибка в select")
		panic(errdropsapDescriptionview)
		//	return nil, errdropsapDescriptionview
	}

	_, errdropstock7813view := r.store.db.Exec(dropstock7813view)
	if errdropstock7813view != nil {
		//fmt.Println("ошибка в select")
		panic(errdropstock7813view)
		//	return nil, errdropstock7813view
	}

	_, errdropstockSum7813view := r.store.db.Exec(dropstockSum7813view)
	if errdropstockSum7813view != nil {
		//fmt.Println("ошибка в select")
		panic(errdropstockSum7813view)
		//	return nil, errdropstockSum7813view
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("ошибка в rows.Err")
		return nil, err
	}

	fmt.Println(idreturnList)

	return &idreturnList, nil

}
