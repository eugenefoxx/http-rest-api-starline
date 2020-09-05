package sqlstore

import (
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
)

// IDReturnRepository ...
type IDReturnRepository struct {
	store *Store
}

// InterDate ...
func (r *IDReturnRepository) InterDate(s *model.IDReturn) error {
	_, err := r.store.db.Exec(
		"INSERT INTO id_return (material, id_roll, lot, qty_fact, qty_sap, qty_panacim, id, lastname) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		s.Material,
		s.IDRoll,
		s.Lot,
		s.QtyFact,
		s.QtySAP,
		s.QtyPanacim,
		s.ID,
		s.LastName,
	)

	if err != nil {

		panic(err)
	}

	return nil
}
