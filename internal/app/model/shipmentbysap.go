package model

import "time"

// Shipmenbysap ...
type Shipmenbysap struct {
	Material     int
	Qty          int64
	ShipmentDate time.Time
	ID           int
	Comment      string
}
