package sqlstore

import (
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
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
		"INSERT INTO shipmentbysap (material, qty, comment, id) VALUES ($1, $2, $3, $4)",
		s.Material,
		s.Qty,
		s.Comment,
		s.ID,
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
