package model

// MainHUMOSAPStock The function will return a pointer to the array of type PanacimStock
func MainHUMOSAPStock() *HUMOSAPStocks {
	return &HUMOSAPStocks{}
}

// HUMOSAPStock ...
type HUMOSAPStock struct {
	ID        int     `json:"id"`
	Package   int     `json:"package"`
	Plant     string  `json:"plant"`
	Warehouse int     `json:"warehouse"`
	Material  int     `json:"material"`
	Quantity  float64 `json:"quantity"`
	Pcs       string  `json:"pcs"`
	Lot       string  `json:"lot"`
	Status1   string  `json:"status1"`
	Status2   string  `json:"status2"`
	SPP       string  `json:"spp"`
}

// HUMOSAPStocks ...
type HUMOSAPStocks []HUMOSAPStock
