package sqlstore

import (
	"database/sql"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	_ "github.com/lib/pq" // ...
	//	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	//	_ "github.com/golang-migrate/migrate/v4/source/github"
)

// Store ...
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	/////	shipmentbysapRepository *ShipmentbysapRepository
	shipmentbysapRepository *ShipmentbysapRepository
}

// New ...
func New(db *sql.DB) *Store {
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
