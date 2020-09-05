package store

// Store ...
type Store interface {
	User() UserRepository
	Shipmentbysap() ShipmentbysapRepository
	IDReturn() IDReturnRepository
	PanacimStock() PanacimStockRepository
	HUMOSAPStock() HUMOSAPStockRepository
	MB52SAPStock() MB52SAPStockRepository
	Showdateidreturn() ShowdateidreturnRepository
}
