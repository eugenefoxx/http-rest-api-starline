package model

// NewPanacimStock The function will return a pointer to the array of type PanacimStock
func NewPanacimStock() *PanacimStocks {
	return &PanacimStocks{}
}

// PanacimStock ..
type PanacimStock struct {
	ZpuNumber                string `json:"zpu_number"`
	Slot                     int    `json:"slot"`
	Subslot                  int    `json:"subslot"`
	CheckInTime              string `json:"check_in_time"`
	CheckInOperator          string `json:"check_in_operator"`
	LocationPana             string `json:"location_pana"`
	StatePana                string `json:"state_pana"`
	StorageUnit              string `json:"storage_unit"`
	StorageSubUnit           string `json:"storage_sub_unit"`
	ZoneState                string `json:"zone_state"`
	ItemState                string `json:"item_state"`
	PartNumber               int    `json:"part_number"`
	LotNo                    string `json:"lot_no"`
	InitialQuantity          int    `json:"initial_quantity"`
	MaterialBarcode          int64  `json:"material_barcode"`
	EstimatedCurrentQuantity int    `json:"estimated_current_quantity"`
	ComponentCountUpdate     string `json:"component_count_update"`
	ReservedQuantity         int    `json:"reserved_quantity"`
	WorkOrderName            string `json:"work_order_name"`
	CartID                   string `json:"cart_id"`
	SpliceReel               int    `json:"splice_reel"`
	SpliceLotNumber          string `json:"splice_lot_number"`
	SpliceEstimatedQuantity  int    `json:"splice_estimated_quantity"`
	SpliceTotalQuantity      int    `json:"splice_total_quantity"`
	FeederBirthDate          string `json:"feeder_birth_date"`
	ComponentsFed            int64  `json:"components_fed"`
	TotalComponentsFed       int64  `json:"total_components_fed"`
	FeederCountUpdate        string `json:"feeder_count_update"`
	LastMaintained           string `json:"last_maintained"`
	FeederBarcode            string `json:"feeder_barcode"`
	Column31                 string `json:"column_31"`
}

// PanacimStocks ...
type PanacimStocks []PanacimStock
