package sqlstore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
)

// ShipmentbysapRepository ...
type ShipmentbysapRepository struct {
	store *Store
}

// InterDate ...
func (r *ShipmentbysapRepository) InterDate(s *model.Shipmentbysap) error {
	//	s := model.Shipmenbysap{}
	/*
		return r.store.db.QueryRow(
			"INSERT INTO shipmentbysap (material, qty, comment) VALUES ($1, $2, $3)",
			s.Material,
			s.Qty,
			s.Comment,
		//	s.ShipmentDate,
		//	s.ID,
		).Scan(
			&s.Material,
			&s.Qty,
			&s.Comment,
		//	&s.ShipmentDate,
		//	&s.ID,
		)

	*/
	const layoutISO = "2006-01-02"
	const layoutTime = "15:04:05"
	_, err := r.store.db.Exec(
		// if err := r.store.db.QueryRow(
		"INSERT INTO shipmentbysap (material, qty, comment, id, lastname) VALUES ($1, $2, $3, $4, $5)",
		s.Material,
		s.Qty,
		s.Comment,
		s.ID,
		s.LastName,

	//	s.ID,
	) /*.Scan(
		&s.Material,
		&s.Qty,
		&s.Comment,
	//	&s.ID,
	)*/

	if err != nil {
		//	if err == sql.ErrNoRows {
		//	return store.ErrRecordNotFound
		panic(err)
	}
	//	return err
	//}
	//	return add

	//	lastInsertID, err := result.LastInsertId()
	//	rowsAffected, err := result.RowsAffected()

	//	fmt.Printf("Product with id=%d created successfully (%d row affected)\n", lastInsertID, rowsAffected)

	return nil
}

// ShowDate ...
func (r *ShipmentbysapRepository) ShowDate() (*model.Shipmentbysaps, error) { // *model.Shipmentbysaps

	shipment := model.Shipmentbysap{}
	shipmentList := make(model.Shipmentbysaps, 0)

	const layoutISO = "2006-01-02"
	const layoutTime = "15:04:05"

	type ShipmentFormat struct {
		Material     int    `db:"material"`
		Qty          int    `db:"qty"`
		Comment      string `db:"comment"`
		ShipmentDate string `db:"shipment_date"`
		ShipmentTime string `db:"shipment_time"`
		LastName     string `db:"lastname"`
	}

	type Datapoint struct {
		Date time.Time
		Time time.Time
	}
	//	var dp Datapoint
	rows, err := r.store.db.Query(
		"SELECT material, qty, comment, TO_CHAR(shipment_date, 'YYYY-MM-DD') shipment_date2, TO_CHAR(shipment_time, 'HH24:MI:SS') shipment_time2, lastname FROM shipmentbysap ORDER BY shipment_date2 DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//pp := ShipmentFormat{}

	for rows.Next() {

		err := rows.Scan(
			//	&shipment.ID,
			&shipment.Material,
			&shipment.Qty,
			&shipment.Comment,
			//	&shipment.ShipmentDate, //&dp.Date, ShipmentDate
			&shipment.ShipmentDate2,
			&shipment.ShipmentTime2, //&dp.Time,
			&shipment.LastName,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		/*
			for _, v := range pp {
				u := &model.Shipmentbysap{
					shipment.Material:                        v.Material,
					shipment.Qty:                             v.Qty,
					shipment.Comment:                         v.Comment,
					shipment.ShipmentDate.Format(layoutISO):  v.ShipmentDate,
					shipment.ShipmentTime.Format(layoutTime): v.ShipmentTime,
					shipment.LastName:                        v.LastName,

					fmt.Println(u),
				}
			}*/
		//	shipment = shipment.Material
		//	shipment = shipment.Qty
		//	shipment = shipment.Comment
		//	shipment = shipment.ShipmentDate.Format(layoutISO)
		//	shipment = shipment.ShipmentTime.Format(layoutTime)
		//	shipment = shipment.LastName
		//	pp.Material = shipment.Material
		//	pp.Qty = shipment.Qty
		//	pp.Comment = shipment.Comment
		//	pp.ShipmentDate = shipment.ShipmentDate.Format(layoutISO)
		//	pp.ShipmentTime = shipment.ShipmentTime.Format(layoutTime)
		//	pp.LastName = shipment.LastName

		//	material := shipment.Material
		//	p := shipment.ShipmentDate.Format(layoutISO)
		//	fmt.Println(pp.Material, pp.Qty, pp.Comment, pp.ShipmentDate, pp.ShipmentTime, pp.LastName)
		//	fmt.Println(shipment.ShipmentDate.Format(layoutISO))
		//	shipment.ShipmentTime.Format(layoutTime)
		//	fmt.Println(db.Date.Format(layoutISO))
		shipmentList = append(shipmentList, shipment)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//	r.store.db.Close()

	//	fmt.Println("YYY", dp.Date.Format(layoutISO))

	return &shipmentList, nil

}

var shipmentDate time.Time
var shipmentTime time.Time

type rawTime []byte

func (t rawTime) TimeTime() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type rawDate []byte

func (t rawDate) TimeDate() (time.Time, error) {
	return time.Parse("2020-02-10", string(t))
}

// ShowDateBySearch ...
func (r *ShipmentbysapRepository) ShowDateBySearch(lastname string, shipment_date2 string, shipment_date3 string, material int) (*model.Shipmentbysaps, error) {

	shipment := model.Shipmentbysap{}
	shipmentList := make(model.Shipmentbysaps, 0)

	rows, err := r.store.db.Query(
		"SELECT material, qty, comment, TO_CHAR(shipment_date, 'YYYY-MM-DD') shipment_date2, TO_CHAR(shipment_time, 'HH24:MI:SS') shipment_time2, lastname FROM shipmentbysap WHERE lastname = $1 AND shipment_date BETWEEN $2 AND $3 AND material = $4",
		lastname, shipment_date2, shipment_date3, material)
	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			//	&shipment.ID,
			&shipment.Material,
			&shipment.Qty,
			&shipment.Comment,
			&shipment.ShipmentDate2, // rawDate ShipmentDate
			&shipment.ShipmentTime2,
			&shipment.LastName,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("ошибка в rows.Scan")
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		shipmentList = append(shipmentList, shipment)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("ошибка в rows.Err")
		return nil, err
	}

	//	r.store.db.Close()

	return &shipmentList, nil

}

// ShowDataByDate ...
func (r *ShipmentbysapRepository) ShowDataByDate(shipmentDate2 string, shipmentDate3 string) (*model.Shipmentbysaps, error) {

	shipment := model.Shipmentbysap{}
	shipmentList := make(model.Shipmentbysaps, 0)

	rows, err := r.store.db.Query(
		"SELECT material, qty, comment, TO_CHAR(shipment_date, 'YYYY-MM-DD') shipment_date2, TO_CHAR(shipment_time, 'HH24:MI:SS') shipment_time2, lastname FROM shipmentbysap WHERE shipment_date BETWEEN $1 AND $2",
		shipmentDate2, shipmentDate3)
	if err != nil {
		fmt.Println("ошибка в select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			//	&shipment.ID,
			&shipment.Material,
			&shipment.Qty,
			&shipment.Comment,
			&shipment.ShipmentDate2, // rawDate ShipmentDate
			&shipment.ShipmentTime2,
			&shipment.LastName,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("ошибка в rows.Scan")
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		shipmentList = append(shipmentList, shipment)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("ошибка в rows.Err")
		return nil, err
	}

	//	r.store.db.Close()

	return &shipmentList, nil

}
