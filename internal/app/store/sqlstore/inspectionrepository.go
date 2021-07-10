package sqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
)

type InspectionRepository struct {
	store *Store
}

func (r *InspectionRepository) InInspection(s *model.Inspection) error {
	_, err := r.store.db.Exec(
		"INSERT INTO transfer (idmaterial, sap, lot, idroll, qty, productiondate, numbervendor, location, lastname) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		s.IdMaterial,
		s.SAP,
		s.Lot,
		s.IdRoll,
		s.Qty,
		s.ProductionDate,
		s.NumberVendor,
		s.Location,
		s.Lastname,
	)

	if err != nil {
		panic(err)
	}

	return nil
}

func (r *InspectionRepository) InInspectionP5(s *model.Inspection) error {
	_, err := r.store.db.Exec(
		"INSERT INTO transferp5 (idmaterial, sap, lot, idroll, qty, productiondate, numbervendor, location, lastname) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		s.IdMaterial,
		s.SAP,
		s.Lot,
		s.IdRoll,
		s.Qty,
		s.ProductionDate,
		s.NumberVendor,
		s.Location,
		s.Lastname,
	)

	if err != nil {
		panic(err)
	}

	return nil
}

func (r *InspectionRepository) ListInspection() (*model.Inspections, error) {
	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)

	deleteDuble := `DELETE FROM transfer a USING transfer b WHERE a.id < b.id AND a.idmaterial = b.idmaterial
	AND a.status IS NULL;`

	_, err := r.store.db.Exec(
		deleteDuble,
	)
	if err != nil {
		panic(err)
	}

	rows, err := r.store.db.Query(
		"SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, transfer.qty, transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2 FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) WHERE transfer.location ='отгружено на ВК'",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&listInspection.ID,
			&listInspection.IdMaterial,
			&listInspection.SAP,
			&listInspection.Lot,
			&listInspection.IdRoll,
			&listInspection.Qty,
			&listInspection.ProductionDate,
			&listInspection.NameDebitor,
			&listInspection.Location,
			&listInspection.Lastname,
			&listInspection.Status,
			&listInspection.Note,
			&listInspection.Update,
			&listInspection.Dateupdate2,
			&listInspection.Timeupdate2,
			&listInspection.Date2,
			&listInspection.Time2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &listInspectionList, nil
}

func (r *InspectionRepository) ListInspectionP5() (*model.Inspections, error) {
	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)

	deleteDuble := `DELETE FROM transferp5 a USING transferp5 b WHERE a.id < b.id AND a.idmaterial = b.idmaterial
	AND a.status IS NULL;`

	_, err := r.store.db.Exec(
		deleteDuble,
	)
	if err != nil {
		panic(err)
	}

	rows, err := r.store.db.Query(
		"SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll, transferp5.qty, transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''), Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''), Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2, TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2 FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) WHERE transferp5.location ='отгружено на ВК'",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&listInspection.ID,
			&listInspection.IdMaterial,
			&listInspection.SAP,
			&listInspection.Lot,
			&listInspection.IdRoll,
			&listInspection.Qty,
			&listInspection.ProductionDate,
			&listInspection.NameDebitor,
			&listInspection.Location,
			&listInspection.Lastname,
			&listInspection.Status,
			&listInspection.Note,
			&listInspection.Update,
			&listInspection.Dateupdate2,
			&listInspection.Timeupdate2,
			&listInspection.Date2,
			&listInspection.Time2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &listInspectionList, nil
}

func (r *InspectionRepository) CountTotalInspection() (int, error) {

	//	count := &model.Inspection{}
	var count int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transfer.id) FROM transfer WHERE transfer.location ='отгружено на ВК'",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}

func (r *InspectionRepository) CountTotalInspectionP5() (int, error) {

	//	count := &model.Inspection{}
	var count int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transferp5.id) FROM transferp5 WHERE transferp5.location ='отгружено на ВК'",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}

func (r *InspectionRepository) HoldInspection() (int, error) {

	//	count := &model.Inspection{}
	var count int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transfer.id) FROM transfer WHERE transfer.status='NG' AND transfer.location ='отгружено на ВК'",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}

func (r *InspectionRepository) HoldInspectionP5() (int, error) {

	//	count := &model.Inspection{}
	var count int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transferp5.id) FROM transferp5 WHERE transferp5.status='NG' AND transferp5.location ='отгружено на ВК'",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}

func (r *InspectionRepository) NotVerifyComponents() (int, error) {

	//	count := &model.Inspection{}
	var count int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transfer.id) FROM transfer WHERE transfer.status IS NULL AND transfer.location ='отгружено на ВК'",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}

func (r *InspectionRepository) NotVerifyComponentsP5() (int, error) {

	//	count := &model.Inspection{}
	var count int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transferp5.id) FROM transferp5 WHERE transferp5.status IS NULL AND transferp5.location ='отгружено на ВК'",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}

func (r *InspectionRepository) CountDebitor() (*model.Inspections, error) {

	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)
	rows, err := r.store.db.Query(
		"SELECT COUNT (transfer.id), Coalesce (vendor.name_debitor, '') FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) WHERE transfer.location ='отгружено на ВК' GROUP BY vendor.name_debitor ORDER BY COUNT(transfer.id) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&listInspection.Count,
			&listInspection.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &listInspectionList, nil
}

func (r *InspectionRepository) CountDebitorP5() (*model.Inspections, error) {

	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)
	rows, err := r.store.db.Query(
		"SELECT COUNT (transferp5.id), Coalesce (vendor.name_debitor, '') FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) WHERE transferp5.location ='отгружено на ВК' GROUP BY vendor.name_debitor ORDER BY COUNT(transferp5.id) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&listInspection.Count,
			&listInspection.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &listInspectionList, nil
}

func (r *InspectionRepository) HoldCountDebitor() (*model.Inspections, error) {

	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)
	rows, err := r.store.db.Query(
		"SELECT COUNT (transfer.status), Coalesce (vendor.name_debitor, '') FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) WHERE transfer.location ='отгружено на ВК' AND transfer.status = 'NG' GROUP BY vendor.name_debitor ORDER BY COUNT(transfer.id) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&listInspection.Count,
			&listInspection.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &listInspectionList, nil
}

func (r *InspectionRepository) HoldCountDebitorP5() (*model.Inspections, error) {

	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)
	rows, err := r.store.db.Query(
		"SELECT COUNT (transferp5.status), Coalesce (vendor.name_debitor, '') FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) WHERE transferp5.location ='отгружено на ВК' AND transferp5.status = 'NG' GROUP BY vendor.name_debitor ORDER BY COUNT(transferp5.id) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&listInspection.Count,
			&listInspection.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &listInspectionList, nil
}

// NotVerifyDebitor Подсчет еще не проверенных дебиторов
func (r *InspectionRepository) NotVerifyDebitor() (*model.Inspections, error) {

	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)
	rows, err := r.store.db.Query(
		"SELECT COUNT (transfer.id), Coalesce (vendor.name_debitor, '') FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) WHERE transfer.location ='отгружено на ВК' AND transfer.status IS NULL GROUP BY vendor.name_debitor ORDER BY COUNT(transfer.id) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&listInspection.Count,
			&listInspection.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &listInspectionList, nil
}

// NotVerifyDebitor Подсчет еще не проверенных дебиторов
func (r *InspectionRepository) NotVerifyDebitorP5() (*model.Inspections, error) {

	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)
	rows, err := r.store.db.Query(
		"SELECT COUNT (transferp5.id), Coalesce (vendor.name_debitor, '') FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) WHERE transferp5.location ='отгружено на ВК' AND transferp5.status IS NULL GROUP BY vendor.name_debitor ORDER BY COUNT(transferp5.id) DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&listInspection.Count,
			&listInspection.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &listInspectionList, nil
}

func (r *InspectionRepository) CountVerifyComponents() (int, error) {

	//	count := &model.Inspection{}
	var countOK int
	//	var countNG int
	//	var countTotal int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transfer.id) FROM transfer WHERE transfer.status='OK' AND transfer.location ='отгружено на ВК'",
	)
	ok := row.Scan(&countOK)
	if ok != nil {
		log.Fatal(ok)
	}
	/*	row = r.store.db.QueryRow(
			"SELECT COUNT (transfer.id) FROM transfer WHERE transfer.status='NG' AND transfer.location ='отгружено на ВК'",
		)
		ng := row.Scan(&countNG)
		if ng != nil {
			log.Fatal(ng)
		}
		countTotal = countOK + countNG*/

	return countOK, nil

}

func (r *InspectionRepository) CountVerifyComponentsP5() (int, error) {

	//	count := &model.Inspection{}
	var countOK int
	//	var countNG int
	//	var countTotal int
	row := r.store.db.QueryRow(
		"SELECT COUNT (transferp5.id) FROM transferp5 WHERE transferp5.status='OK' AND transferp5.location ='отгружено на ВК'",
	)
	ok := row.Scan(&countOK)
	if ok != nil {
		log.Fatal(ok)
	}
	/*	row = r.store.db.QueryRow(
			"SELECT COUNT (transferp5.id) FROM transferp5 WHERE transferp5.status='NG' AND transferp5.location ='отгружено на ВК'",
		)
		ng := row.Scan(&countNG)
		if ng != nil {
			log.Fatal(ng)
		}
		countTotal = countOK + countNG*/

	return countOK, nil

}

func (r *InspectionRepository) EditInspection(id int) (*model.Inspection, error) {

	u := &model.Inspection{}
	//	fmt.Println("EditInspection -", groups)
	//	if groups == groupQuality {
	if err := r.store.db.QueryRow(
		"SELECT id, Coalesce (status, ''), Coalesce (note, '') FROM transfer WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		//	&u.IdRoll,
		&u.Status,
		&u.Note,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	//	}
	return u, nil

}

func (r *InspectionRepository) EditInspectionP5(id int) (*model.Inspection, error) {

	u := &model.Inspection{}
	//	fmt.Println("EditInspection -", groups)
	//	if groups == groupQuality {
	if err := r.store.db.QueryRow(
		"SELECT id, Coalesce (status, ''), Coalesce (note, '') FROM transferp5 WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		//	&u.IdRoll,
		&u.Status,
		&u.Note,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	//	}
	return u, nil

}

func (r *InspectionRepository) UpdateInspection(s *model.Inspection) error {

	//	if groups == groupQuality {
	_, err := r.store.db.Exec(
		"UPDATE transfer SET status = $1, note = $2, update = $3, dateupdate = $4, timeupdate = $5 FROM users WHERE transfer.id = $6 AND users.groups = $7",
		s.Status,
		s.Note,
		s.Update,
		s.Dateupdate,
		s.Timeupdate,
		s.ID,
		s.Groups,
	)

	if err != nil {
		panic(err)
	}
	//	}

	return nil
}

func (r *InspectionRepository) UpdateInspectionP5(s *model.Inspection) error {

	//	if groups == groupQuality {
	_, err := r.store.db.Exec(
		"UPDATE transferp5 SET status = $1, note = $2, update = $3, dateupdate = $4, timeupdate = $5 FROM users WHERE transferp5.id = $6 AND users.groups = $7",
		s.Status,
		s.Note,
		s.Update,
		s.Dateupdate,
		s.Timeupdate,
		s.ID,
		s.Groups,
	)

	if err != nil {
		panic(err)
	}
	//	}

	return nil
}

func (r *InspectionRepository) DeleteItemInspection(s *model.Inspection) error {
	_, err := r.store.db.Exec(
		"DELETE FROM transfer WHERE id = $1",
		s.ID,
	)
	if err != nil {
		panic(err)
	}

	return nil
}

func (r *InspectionRepository) DeleteItemInspectionP5(s *model.Inspection) error {
	_, err := r.store.db.Exec(
		"DELETE FROM transferp5 WHERE id = $1",
		s.ID,
	)
	if err != nil {
		panic(err)
	}

	return nil
}

/*
func (r *InspectionRepository) EditInspectionForWarehouse(id int) (*model.Inspection, error) {

	u := &model.Inspection{}
	//	fmt.Println("EditInspection -", groups)
	//	if groups == groupQuality {
	if err := r.store.db.QueryRow(
		"SELECT id, Coalesce (notewarehouse, '') FROM transfer WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		//	&u.IdRoll,
		//&u.Status,
		&u.Notewarehouse,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	//	}
	return u, nil

}
*/
func (r *InspectionRepository) ListAcceptWHInspection() (*model.Inspections, error) {
	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)

	rows, err := r.store.db.Query(
		"SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, Coalesce (transfer.status, ''), Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2 FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) WHERE (transfer.location ='отгружено на ВК' AND transfer.status = 'OK') OR (transfer.location ='отгружено на ВК' AND transfer.status = 'NG')",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&listInspection.ID,
			&listInspection.IdMaterial,
			&listInspection.SAP,
			&listInspection.Lot,
			&listInspection.IdRoll,
			&listInspection.ProductionDate,
			&listInspection.NameDebitor,
			&listInspection.Location,
			&listInspection.Status,
			&listInspection.Note,
			&listInspection.Update,
			&listInspection.Dateupdate2,
			&listInspection.Timeupdate2,
			&listInspection.Date2,
			&listInspection.Time2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		listInspectionList = append(listInspectionList, listInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &listInspectionList, nil
}

func (r *InspectionRepository) EditAcceptWarehouseInspection(id int) (*model.Inspection, error) {
	u := &model.Inspection{}
	fmt.Println("EditAcceptWarehouseInspection -")

	//	if groups == groupWarehouse {

	if err := r.store.db.QueryRow(
		"SELECT id, Coalesce (location, '') FROM transfer WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Location,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	//	}
	return u, nil
}

func (r *InspectionRepository) EditAcceptWarehouseInspectionP5(id int) (*model.Inspection, error) {
	u := &model.Inspection{}
	fmt.Println("EditAcceptWarehouseInspectionP5 -")

	//	if groups == groupWarehouse {

	if err := r.store.db.QueryRow(
		"SELECT id, Coalesce (location, '') FROM transferp5 WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Location,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	//	}
	return u, nil
}

func (r *InspectionRepository) AcceptWarehouseInspection(s *model.Inspection) error {
	//	if groups == groupWarehouse {
	_, err := r.store.db.Exec(
		"UPDATE transfer SET location = $1, lastnameaccept = $2, dateaccept = $3, timeaccept = $4 FROM users WHERE transfer.id = $5 AND users.groups = $6",
		s.Location,
		s.Lastnameaccept,
		s.Dateaccept,
		s.Timeaccept,
		s.ID,
		s.Groups,
	)
	if err != nil {
		panic(err)
	}
	//	}
	return nil
}

func (r *InspectionRepository) AcceptWarehouseInspectionP5(s *model.Inspection) error {
	//	if groups == groupWarehouse {
	_, err := r.store.db.Exec(
		"UPDATE transferp5 SET location = $1, lastnameaccept = $2, dateaccept = $3, timeaccept = $4 FROM users WHERE transferp5.id = $5 AND users.groups = $6",
		s.Location,
		s.Lastnameaccept,
		s.Dateaccept,
		s.Timeaccept,
		s.ID,
		s.Groups,
	)
	if err != nil {
		panic(err)
	}
	//	}
	return nil
}

func (r *InspectionRepository) AcceptGroupsWarehouseInspection(s *model.Inspection) error {
	//	if groups == groupWarehouse {
	_, err := r.store.db.Exec(
		"UPDATE transfer SET location = $1, lastnameaccept = $2, dateaccept = $3, timeaccept = $4 FROM users WHERE transfer.idmaterial = $5 AND users.groups = $6",
		s.Location,
		s.Lastnameaccept,
		s.Dateaccept,
		s.Timeaccept,
		s.IdMaterial,
		s.Groups,
	)
	if err != nil {
		panic(err)
	}
	//	}
	return nil
}

func (r *InspectionRepository) AcceptGroupsWarehouseInspectionP5(s *model.Inspection) error {
	//	if groups == groupWarehouse {
	_, err := r.store.db.Exec(
		"UPDATE transferp5 SET location = $1, lastnameaccept = $2, dateaccept = $3, timeaccept = $4 FROM users WHERE transferp5.idmaterial = $5 AND users.groups = $6",
		s.Location,
		s.Lastnameaccept,
		s.Dateaccept,
		s.Timeaccept,
		s.IdMaterial,
		s.Groups,
	)
	if err != nil {
		panic(err)
	}
	//	}
	return nil
}

func (r *InspectionRepository) ListShowDataByEO(eo string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)
	//showDataByDateList := model.Inspections{}

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
		Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
		Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
		Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) 
		WHERE idroll = $1;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		eo)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataByEOP5(eo string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll, 
	transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''), 
		Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''), 
		Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2, 
		Coalesce (transferp5.lastnameaccept, ''), Coalesce(TO_CHAR(transferp5.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transferp5.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) 
		WHERE idroll = $1;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		eo)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataBySap(sap int) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
		Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
		Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
		Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) 
		WHERE sap = $1;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		sap)

	if err != nil {
		fmt.Println("ошибка в selectDate")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			/*	if err == sql.ErrNoRows {
					return nil, store.ErrRecordNotFound
				}
				return nil, err
			*/
			if errors.Is(err, sql.ErrNoRows) {
				return nil, store.ErrRecordNotFound
			} else {
				return nil, err
			}
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

// ListShowDataBySapPagination функция для создания постраничного вывода списка
func (r *InspectionRepository) ListShowDataBySapPagination(sap, begin, limit int) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
		Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
		Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
		Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) 
		WHERE sap = $1 offset $2 limit $3;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		sap,
		begin,
		limit)

	if err != nil {
		fmt.Println("ошибка в selectDate")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			/*	if err == sql.ErrNoRows {
					return nil, store.ErrRecordNotFound
				}
				return nil, err
			*/
			if errors.Is(err, sql.ErrNoRows) {
				return nil, store.ErrRecordNotFound
			} else {
				return nil, err
			}
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataBySapP5(sap int) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll,
	transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''),
		Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''),
		Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2,
		Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2,
		TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2,
		Coalesce (transferp5.lastnameaccept, ''), Coalesce(TO_CHAR(transferp5.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,
		Coalesce(TO_CHAR(transferp5.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor)
		WHERE sap = $1;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		sap,
	)

	if err != nil {
		fmt.Println("ошибка в selectDate")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			//	if err == sql.ErrNoRows {
			//		return nil, store.ErrRecordNotFound
			//	}
			//	return nil, err
			//
			if errors.Is(err, sql.ErrNoRows) {
				return nil, store.ErrRecordNotFound
			} else {
				return nil, err
			}
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

/*
func (r *InspectionRepository) ListShowDataBySapP5(sap, begin, limit int) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll,
	transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''),
		Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''),
		Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2,
		Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2,
		TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2,
		Coalesce (transferp5.lastnameaccept, ''), Coalesce(TO_CHAR(transferp5.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,
		Coalesce(TO_CHAR(transferp5.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor)
		WHERE sap = $1 offset $2 limit $3;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		sap,
		begin,
		limit)

	if err != nil {
		fmt.Println("ошибка в selectDate")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			//	if err == sql.ErrNoRows {
			//		return nil, store.ErrRecordNotFound
			//	}
			//	return nil, err
			//
			if errors.Is(err, sql.ErrNoRows) {
				return nil, store.ErrRecordNotFound
			} else {
				return nil, err
			}
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}
*/
func (r *InspectionRepository) ListShowDataByDate(updateDate1 string, updateDate2 string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
		Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
		Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
		Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) 
		WHERE dateupdate BETWEEN $1 AND $2;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		updateDate1, updateDate2)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataByDateP5(updateDate1 string, updateDate2 string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll, 
	transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''), 
		Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''), 
		Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2, 
		Coalesce (transferp5.lastnameaccept, ''), Coalesce(TO_CHAR(transferp5.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transferp5.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) 
		WHERE dateupdate BETWEEN $1 AND $2;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		updateDate1, updateDate2)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataByDateAndSAP(updateDate1 string, updateDate2 string, sap int) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
		Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
		Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
		Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) 
		WHERE dateupdate BETWEEN $1 AND $2 AND sap = $3;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		updateDate1, updateDate2, sap)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataByDateAndSAPP5(updateDate1 string, updateDate2 string, sap int) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll, 
	transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''), 
		Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''), 
		Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2, 
		Coalesce (transferp5.lastnameaccept, ''), Coalesce(TO_CHAR(transferp5.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transferp5.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) 
		WHERE dateupdate BETWEEN $1 AND $2 AND sap = $3;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		updateDate1, updateDate2, sap)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataByDateAndEO(updateDate1 string, updateDate2 string, eo string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
		Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
		Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
		Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) 
		WHERE dateupdate BETWEEN $1 AND $2 AND idroll = $3;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		updateDate1, updateDate2, eo)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) ListShowDataByDateAndEOP5(updateDate1 string, updateDate2 string, eo string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transferp5.id, transferp5.idmaterial, transferp5.sap, transferp5.lot, transferp5.idroll, 
	transferp5.productiondate, Coalesce (vendor.name_debitor, ''), transferp5.location, transferp5.lastname, Coalesce (transferp5.status, ''), 
		Coalesce (transferp5.note, ''), Coalesce (transferp5.update, ''), 
		Coalesce(TO_CHAR(transferp5.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
		Coalesce(TO_CHAR(transferp5.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
		TO_CHAR(transferp5.date, 'YYYY-MM-DD') date2, TO_CHAR(transferp5.time, 'HH24:MI:SS') time2, 
		Coalesce (transferp5.lastnameaccept, ''), Coalesce(TO_CHAR(transferp5.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
		Coalesce(TO_CHAR(transferp5.timeaccept, 'HH24:MI:SS'), '') timeaccept2
		FROM transferp5 left outer join vendor on (transferp5.numbervendor = vendor.code_debitor) 
		WHERE dateupdate BETWEEN $1 AND $2 AND idroll = $3;`

	rows, err := r.store.db.Query(
		//	"SELECT * FROM transfer WHERE dateupdate BETWEEN $1 AND $2",
		selectDate,
		updateDate1, updateDate2, eo)

	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showDataByDate.ID,
			&showDataByDate.IdMaterial,
			&showDataByDate.SAP,
			&showDataByDate.Lot,
			&showDataByDate.IdRoll,
			&showDataByDate.ProductionDate,
			&showDataByDate.NameDebitor,
			&showDataByDate.Location,
			&showDataByDate.Lastname,
			&showDataByDate.Status,
			&showDataByDate.Note,
			&showDataByDate.Update,
			&showDataByDate.Dateupdate2,
			&showDataByDate.Timeupdate2,
			&showDataByDate.Date2,
			&showDataByDate.Time2,
			&showDataByDate.Lastnameaccept,
			&showDataByDate.Dateaccept2,
			&showDataByDate.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showDataByDateList = append(showDataByDateList, showDataByDate)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showDataByDateList, nil
}

func (r *InspectionRepository) CountInspection() (int, error) {

	//select := `SELECT COUNT (transfer.id) FROM transfer;`
	var count int
	row := r.store.db.QueryRow(
		//			select
		"SELECT COUNT (transfer.id) FROM transfer",
	)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil
}

func (r *InspectionRepository) PaginationInspection(begin, limit int) (s *model.Inspections, err error) {
	showPaginationInspection := model.Inspection{}
	showPaginationInspectionList := make(model.Inspections, 0)

	selectOffset := `transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
	transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, transfer.lastname, Coalesce (transfer.status, ''), 
	Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), 
	Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, 
	Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, 
	TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2, 
	Coalesce (transfer.lastnameaccept, ''), Coalesce(TO_CHAR(transfer.dateaccept, 'YYYY-MM-DD'), '') dateaccept2,  
	Coalesce(TO_CHAR(transfer.timeaccept, 'HH24:MI:SS'), '') timeaccept2
	FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) offset = $1 limit = $2;`

	rows, err := r.store.db.Query(
		selectOffset,
		begin,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&showPaginationInspection.ID,
			&showPaginationInspection.IdMaterial,
			&showPaginationInspection.SAP,
			&showPaginationInspection.Lot,
			&showPaginationInspection.IdRoll,
			&showPaginationInspection.ProductionDate,
			&showPaginationInspection.NameDebitor,
			&showPaginationInspection.Location,
			&showPaginationInspection.Lastname,
			&showPaginationInspection.Status,
			&showPaginationInspection.Note,
			&showPaginationInspection.Update,
			&showPaginationInspection.Dateupdate2,
			&showPaginationInspection.Timeupdate2,
			&showPaginationInspection.Date2,
			&showPaginationInspection.Time2,
			&showPaginationInspection.Lastnameaccept,
			&showPaginationInspection.Dateaccept2,
			&showPaginationInspection.Timeaccept2,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		showPaginationInspectionList = append(showPaginationInspectionList, showPaginationInspection)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &showPaginationInspectionList, nil
}
