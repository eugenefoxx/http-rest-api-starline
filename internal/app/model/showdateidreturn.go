package model

// NewShowdateidreturn The function will return a pointer to the array of type Showdateidreturn
func NewShowdateidreturn() *Showdateidreturns {
	return &Showdateidreturns{}
}

// Showdateidreturn ...
type Showdateidreturn struct {
	IDRoll                   int    `json:"idroll"`
	Material                 int    `json:"material"`
	MaterialDescription      string `json:"materialdescription"`
	Lot                      string `json:"lot"`
	QtyFact                  int    `json:"qtyfact"`
	QtySAP                   int    `json:"qtysap"`
	QtyPanacim               int    `json:"qtypanacim"`
	SPP                      string `json:"spp,omitempty"`
	Summa                    int64  `json:"summa,omitempty"`
	Warehouse                string `json:"warehouse,omitempty"`
	QtyIDSAP                 int    `json:"qty_id_sap"`
	EstimatedCurrentQuantity int    `json:"estimatedcurrentquantity,omitempty"`
	StorageUnit              string `json:"storageunit,omitempty"`
	StorageSubUnit           string `json:"storagesubunit,omitempty"`
	ShipmentDate2            string `json:"shipmentdate2"`
	ShipmentTime2            string `json:"shipmenttime2"`
	Lastname                 string `json:"lastname"`
}

// Showdateidreturns ...
type Showdateidreturns []Showdateidreturn
