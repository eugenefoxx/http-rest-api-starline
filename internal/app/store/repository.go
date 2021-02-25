package store

import "github.com/eugenefoxx/http-rest-api-starline/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string, string) (*model.User, error)
	UpdatePass(*model.User) error
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

// VendorRepository ...
type VendorRepository interface {
	InsertVendor(*model.Vendor) error
	EditVendor(int) (*model.Vendor, error)
	UpdateVendor(*model.Vendor) error
	ListVendor() (*model.Vendors, error)
	DeleteVendor(*model.Vendor) error
}

type InspectionRepository interface {
	InInspection(*model.Inspection) error
	ListInspection() (*model.Inspections, error)
	EditInspection(int, string) (*model.Inspection, error)
	UpdateInspection(*model.Inspection, string) error
	ListAcceptWHInspection() (*model.Inspections, error)
	EditAcceptWarehouseInspection(int, string) (*model.Inspection, error)
	AcceptWarehouseInspection(*model.Inspection, string) error
	CountTotalInspection() (int, error)
	HoldInspection() (int, error)
	NotVerifyComponents() (int, error)
	DeleteItemInspection(*model.Inspection) error
	//	CountDebitor() (int, string, error)
	CountDebitor() (*model.Inspections, error)
	HoldCountDebitor() (*model.Inspections, error)
	NotVerifyDebitor() (*model.Inspections, error)
	ListShowDataByDate(string, string) (*model.Inspections, error)
	ListShowDataByDateAndSAP(string, string, int) (*model.Inspections, error)
	ListShowDataByDateAndEO(string, string, string) (*model.Inspections, error)
}
