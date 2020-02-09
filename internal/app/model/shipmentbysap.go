package model

import "time"

// Shipmentbysap ...
type Shipmentbysap struct {
	Material     int
	Qty          int64
	ShipmentDate time.Time
	ID           int
	Comment      string
}
