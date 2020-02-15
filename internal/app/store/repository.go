package store

import "github.com/eugenefoxx/http-rest-api-starline/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

// ShipmentbysapRepository ...
type ShipmentbysapRepository interface {
	InterDate(*model.Shipmentbysap) error
	ShowDate() (*model.Shipmentbysaps, error)
}
