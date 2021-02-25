package sqlstore

import (
	"database/sql"
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

func (r *InspectionRepository) ListInspection() (*model.Inspections, error) {
	listInspection := model.Inspection{}
	listInspectionList := make(model.Inspections, 0)

	rows, err := r.store.db.Query(
		"SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, Coalesce (transfer.status, ''), Coalesce (transfer.note, ''), Coalesce (transfer.update, ''), Coalesce(TO_CHAR(transfer.dateupdate, 'YYYY-MM-DD'), '') dateupdate2, Coalesce(TO_CHAR(transfer.timeupdate, 'HH24:MI:SS'), '') timeupdate2, TO_CHAR(transfer.date, 'YYYY-MM-DD') date2, TO_CHAR(transfer.time, 'HH24:MI:SS') time2 FROM transfer left outer join vendor on (transfer.numbervendor = vendor.code_debitor) WHERE transfer.location ='отгружено на ВК'",
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

func (r *InspectionRepository) EditInspection(id int, groups string) (*model.Inspection, error) {

	u := &model.Inspection{}
	fmt.Println("EditInspection -", groups)
	if groups == "качество" {
		if err := r.store.db.QueryRow(
			"SELECT id, Coalesce (status, ''), Coalesce (note, '') FROM transfer WHERE id = $1",
			id,
		).Scan(
			&u.ID,
			&u.Status,
			&u.Note,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
	}
	return u, nil

}

func (r *InspectionRepository) UpdateInspection(s *model.Inspection, groups string) error {

	if groups == "качество" {
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
	}

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

func (r *InspectionRepository) EditAcceptWarehouseInspection(id int, groups string) (*model.Inspection, error) {
	u := &model.Inspection{}
	fmt.Println("EditAcceptWarehouseInspection -", groups)

	if groups == "склад" {

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

	}
	return u, nil
}

func (r *InspectionRepository) AcceptWarehouseInspection(s *model.Inspection, groups string) error {
	if groups == "склад" {
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
	}
	return nil
}

func (r *InspectionRepository) ListShowDataByDate(updateDate1 string, updateDate2 string) (s *model.Inspections, err error) {
	showDataByDate := model.Inspection{}
	showDataByDateList := make(model.Inspections, 0)

	selectDate := `SELECT transfer.id, transfer.idmaterial, transfer.sap, transfer.lot, transfer.idroll, 
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, Coalesce (transfer.status, ''), 
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
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, Coalesce (transfer.status, ''), 
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
		transfer.productiondate, Coalesce (vendor.name_debitor, ''), transfer.location, Coalesce (transfer.status, ''), 
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
