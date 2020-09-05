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
	//	ShowDateBySearch(*model.Shipmentbysap) (*model.Shipmentbysaps, error)
	ShowDateBySearch(string, string, string, int) (*model.Shipmentbysaps, error)
	ShowDataByDate(string, string) (*model.Shipmentbysaps, error)
}

// IDReturnRepository ...
type IDReturnRepository interface {
	InterDate(*model.IDReturn) error
}

// PanacimStockRepository ...
type PanacimStockRepository interface {
	ImportDate()
}

//HUMOSAPStockRepository ...
type HUMOSAPStockRepository interface {
	ImportDate()
}

// MB52SAPStockRepository ...
type MB52SAPStockRepository interface {
	ImportDate()
}

// ShowdateidreturnRepository ...
type ShowdateidreturnRepository interface {
	ShowDataByDate(string, string) (*model.Showdateidreturns, error)
}
