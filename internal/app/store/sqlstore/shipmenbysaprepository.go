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

type Shipmentbysap2 struct {
	Material     int
	Qty          int64
	ShipmentDate time.Time
	ShipmentTime time.Time
	ID           int
	LastName     string
	Comment      string
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
		"SELECT id, material, qty, comment, shipment_date, shipment_time, lastname FROM shipmentbysap",
	)
	if err != nil {
		return nil, err
	}
	//	defer rows.Close()

	for rows.Next() {
		//	p := Shipmentbysap2{}
		//	p := &model.Shipmentbysap{}
		err := rows.Scan(
			&shipment.ID,
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

	r.store.db.Close()

	/*
		for rows.Next() {
			p := showdate1{}
			err := rows.Scan(&p.Material, &p.Qty, &p.Comment, &p.ShipmentDate, &p.ShipmentTime, &p.LastName)
			if err != nil {
				//		fmt.Println(err)
				panic(err)
				//	continue
			}
			showdate1 = append(showdate1, p)
		}

		//	for _, p := range showdate {
		//		fmt.Println(p.Material, p.Qty, p.Comment, p.ShipmentDate, p.ShipmentTime, &p.LastName)
		//	}
	*/
	return &shipmentList, nil

	/*
		if err := r.store.db.QueryRow(
			"SELECT id, material, qty, comment, shipment_date, shipment_time, lastname FROM shipmentbysap WHERE material=$1", 1014040,
		).Scan(
			&showdate.ID,
			&showdate.Material,
			&showdate.Qty,
			&showdate.Comment,
			&showdate.ShipmentDate,
			&showdate.ShipmentTime,
			&showdate.LastName,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}

		return showdate, nil
	*/

}

//var shipmentDate time.Time
//var shipmentTime time.Time
/*
   type rawTime []byte

   func (t rawTime) Time() (time.Time, error) {
   	return time.Parse("15:04:05", string(t))
   }

   type rawDate []byte

   func (t rawDate) Time() (time.Time, error) {
   	return time.Parse("2020-02-10", string(t))
   }
*/
