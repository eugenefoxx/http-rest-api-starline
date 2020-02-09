package model

import "time"

// Shipmentbysap ...
type Shipmentbysap struct {
	Material     int
	Qty          int64
	ShipmentDate time.Time
	ShipmentTime time.Time
	ID           int
	LastName     string
	Comment      string
}
