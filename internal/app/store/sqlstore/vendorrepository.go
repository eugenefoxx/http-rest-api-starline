package sqlstore

import (
	"database/sql"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
)

type VendorRepository struct {
	store *Store
}

// InsertVendor ...
func (r *VendorRepository) InsertVendor(s *model.Vendor) error {
	_, err := r.store.db.Exec(
		"INSERT INTO vendor (code_debitor, name_debitor) VALUES ($1, $2)",
		s.CodeDebitor,
		s.NameDebitor,
	)

	if err != nil {
		panic(err)
	}

	return nil
}

func (r *VendorRepository) ListVendor() (*model.Vendors, error) {

	editVendor := model.Vendor{}
	editVendorList := make(model.Vendors, 0)

	rows, err := r.store.db.Query(
		"SELECT id, code_debitor, name_debitor FROM vendor",
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&editVendor.ID,
			&editVendor.CodeDebitor,
			&editVendor.NameDebitor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		editVendorList = append(editVendorList, editVendor)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &editVendorList, nil

}

func (r *VendorRepository) EditVendor(id int) (*model.Vendor, error) {

	u := &model.Vendor{}
	if err := r.store.db.QueryRow(
		"SELECT id, code_debitor, name_debitor FROM vendor WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.CodeDebitor,
		&u.NameDebitor,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *VendorRepository) UpdateVendor(s *model.Vendor) error {

	//	updateVendor := model.Vendor{}
	//	updateVendorList := make(model.Vendors, 0)

	_, err := r.store.db.Exec(
		"UPDATE vendor SET code_debitor = $1, name_debitor = $2 WHERE id = $3",
		s.CodeDebitor,
		s.NameDebitor,
		s.ID,
	)

	if err != nil {
		panic(err)
	}

	return nil
}

func (r *VendorRepository) DeleteVendor(s *model.Vendor) error {
	_, err := r.store.db.Exec(
		"DELETE FROM vendor WHERE id = $1",
		s.ID,
	)
	if err != nil {
		panic(err)
	}

	return nil
}
