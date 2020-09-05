package model

// MainMB52SAPStock The function will return a pointer to the array of type MB52SAPStock
func MainMB52SAPStock() *MB52SAPStocks {
	return &MB52SAPStocks{}
}

// MB52SAPStock ...
type MB52SAPStock struct {
	Material           int     `json:"material"`
	Plant              string  `json:"plant"`
	Warehouse          string  `json:"warehouse"`
	MateralDescription string  `json:"materialdescription"`
	QtyFree            float64 `json:"qtyfree"`
	QtyInspection      float64 `json:"qtyinspection"`
	QtyBlock           float64 `json:"qtyblock"`
	Lot                string  `json:"lot"`
	SPP                string  `json:"spp"`
}

// MB52SAPStocks ...
type MB52SAPStocks []MB52SAPStock
