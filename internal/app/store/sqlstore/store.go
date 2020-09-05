package sqlstore

import (
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // ...
	//	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	//	_ "github.com/golang-migrate/migrate/v4/source/github"
	//	"github.com/jackc/pgx/v4"
)

// Store ...
type Store struct {
	db *sqlx.DB
	//	db             *pgx.Conn
	userRepository *UserRepository
	/////	shipmentbysapRepository *ShipmentbysapRepository
	shipmentbysapRepository    *ShipmentbysapRepository
	idreturnRepository         *IDReturnRepository
	panacimstockRepository     *PanacimStockRepository
	humosapstockRepository     *HUMOSAPStockRepository
	mb52sapstockRepository     *MB52SAPStockRepository
	showdateidreturnRepository *ShowdateidreturnRepository
}

// New ...
// func New(db *sql.DB) *Store {
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

// store.User().Create()

// Shipmentbysap ...
func (s *Store) Shipmentbysap() store.ShipmentbysapRepository {
	if s.shipmentbysapRepository != nil {
		return s.shipmentbysapRepository
	}

	s.shipmentbysapRepository = &ShipmentbysapRepository{
		store: s,
	}

	return s.shipmentbysapRepository
}

// IDReturn ...
func (s *Store) IDReturn() store.IDReturnRepository {
	if s.idreturnRepository != nil {
		return s.idreturnRepository
	}

	s.idreturnRepository = &IDReturnRepository{
		store: s,
	}

	return s.idreturnRepository
}

// PanacimStock ...
func (s *Store) PanacimStock() store.PanacimStockRepository {
	if s.panacimstockRepository != nil {
		return s.panacimstockRepository
	}

	s.panacimstockRepository = &PanacimStockRepository{
		store: s,
	}

	return s.panacimstockRepository
}

// HUMOSAPStock ...
func (s *Store) HUMOSAPStock() store.HUMOSAPStockRepository {
	if s.humosapstockRepository != nil {
		return s.humosapstockRepository
	}

	s.humosapstockRepository = &HUMOSAPStockRepository{
		store: s,
	}

	return s.humosapstockRepository
}

// MB52SAPStock ...
func (s *Store) MB52SAPStock() store.MB52SAPStockRepository {
	if s.mb52sapstockRepository != nil {
		return s.mb52sapstockRepository
	}

	s.mb52sapstockRepository = &MB52SAPStockRepository{
		store: s,
	}

	return s.mb52sapstockRepository
}

// Showdateidreturn ...
func (s *Store) Showdateidreturn() store.ShowdateidreturnRepository {
	if s.showdateidreturnRepository != nil {
		return s.showdateidreturnRepository
	}

	s.showdateidreturnRepository = &ShowdateidreturnRepository{
		store: s,
	}

	return s.showdateidreturnRepository
}
