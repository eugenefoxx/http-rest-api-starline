package store

import "github.com/eugenefoxx/http-rest-api-starline/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string, string) (*model.User, error)
	UpdatePass(*model.User) error
	CreateUserByManager(*model.User) error
	ListUsersQuality() (*model.Users, error)
	ListUsersQualityP5() (*model.Users, error)
	EditUserByManager(int) (*model.User, error)
	UpdateUserByManager(*model.User) error
	DeleteUserByManager(*model.User) error
	ListUsersWarehouse(string) (*model.Users, error)
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
	InInspectionP5(*model.Inspection) error
	ListInspection() (*model.Inspections, error)
	ListInspectionP5() (*model.Inspections, error)
	EditInspection(int) (*model.Inspection, error)
	EditInspectionP5(int) (*model.Inspection, error)
	UpdateInspection(*model.Inspection) error
	UpdateInspectionP5(*model.Inspection) error
	ListAcceptWHInspection() (*model.Inspections, error)
	EditAcceptWarehouseInspection(int) (*model.Inspection, error)
	EditAcceptWarehouseInspectionP5(int) (*model.Inspection, error)
	AcceptWarehouseInspection(*model.Inspection) error
	AcceptWarehouseInspectionP5(*model.Inspection) error
	AcceptGroupsWarehouseInspection(*model.Inspection) error
	AcceptGroupsWarehouseInspectionP5(*model.Inspection) error
	CountTotalInspection() (int, error)
	CountTotalInspectionP5() (int, error)
	HoldInspection() (int, error)
	HoldInspectionP5() (int, error)
	NotVerifyComponents() (int, error)
	NotVerifyComponentsP5() (int, error)
	CountVerifyComponents() (int, error)
	CountVerifyComponentsP5() (int, error)
	DeleteItemInspection(*model.Inspection) error
	DeleteItemInspectionP5(*model.Inspection) error
	//	CountDebitor() (int, string, error)
	CountDebitor() (*model.Inspections, error)
	CountDebitorP5() (*model.Inspections, error)
	HoldCountDebitor() (*model.Inspections, error)
	HoldCountDebitorP5() (*model.Inspections, error)
	NotVerifyDebitor() (*model.Inspections, error)
	NotVerifyDebitorP5() (*model.Inspections, error)
	ListShowDataByDate(string, string) (*model.Inspections, error)
	ListShowDataByDateP5(string, string) (*model.Inspections, error)
	ListShowDataByDateAndSAP(string, string, int) (*model.Inspections, error)
	ListShowDataByDateAndSAPP5(string, string, int) (*model.Inspections, error)
	ListShowDataByDateAndEO(string, string, string) (*model.Inspections, error)
	ListShowDataByDateAndEOP5(string, string, string) (*model.Inspections, error)
	ListShowDataBySap(int) (*model.Inspections, error)
	//ListShowDataBySapPagination(int, int, int) (*model.Inspections, error)
	ListShowDataBySapP5(int) (*model.Inspections, error)
	ListShowDataByEO(string) (*model.Inspections, error)
	ListShowDataByEOP5(string) (*model.Inspections, error)
	CountInspection() (int, error)
	PaginationInspection(int, int) (*model.Inspections, error)
}

type RoleQualityRepository interface {
	ListRoleQuality() (*model.RoleQualitys, error)
}
