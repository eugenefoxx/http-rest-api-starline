package sqlstore

import (
	"database/sql"
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
func (r *ShipmentbysapRepository) ShowDate() (*model.Shipmentbysaps, error) {

	shipment := model.Shipmentbysap{}
	shipmentList := make(model.Shipmentbysaps, 0)

	rows, err := r.store.db.Query(
		"SELECT material, qty, comment, shipment_date, shipment_time, lastname FROM shipmentbysap",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			//	&shipment.ID,
			&shipment.Material,
			&shipment.Qty,
			&shipment.Comment,
			&shipment.ShipmentDate,
			&shipment.ShipmentTime,
			&shipment.LastName,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		shipmentList = append(shipmentList, shipment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//	r.store.db.Close()

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
