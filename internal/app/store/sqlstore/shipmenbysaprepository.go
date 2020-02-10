package sqlstore

import (
	"database/sql"

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
func (r *ShipmentbysapRepository) ShowDate(showdate *model.Shipmentbysap) (*model.Shipmentbysap, error) {
	/*
		rows, err := r.store.db.Query(
			"SELECT * FROM shipmentbysap",
			s.Material,
			s.Qty,
			s.Comment,
			s.ShipmentDate,
			s.ShipmentTime,
			s.LastName,
		)

		if err != nil {
			//	if err == sql.ErrNoRows {
			//	return store.ErrRecordNotFound
			panic(err)
		}

		defer rows.Close()
	*/
	//	showdate := &model.Shipmentbysap{}
	/*
		for rows.Next() {
			p := showdate{}
			err := rows.Scan(&p.Material, &p.Qty, &p.Comment, &p.ShipmentDate, &p.ShipmentTime, &p.LastName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			showdate = append(showdate, p)
		}

		for _, p := range showdate {
			fmt.Println(p.Material, p.Qty, p.Comment, p.ShipmentDate, p.ShipmentTime, &p.LastName)
		}
	*/
	if err := r.store.db.QueryRow(
		"SELECT * FROM shipmentbysap",
	).Scan(
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
}
