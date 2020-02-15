package model

import "time"

// Main The function will return a pointer to the array of type Shipmentbysap
func Main() *Shipmentbysaps {
	return &Shipmentbysaps{}
}

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

// Shipmentbysaps ...
type Shipmentbysaps []Shipmentbysap

type rawTime []byte

func (t rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type rawDate []byte

func (t rawDate) Time() (time.Time, error) {
	return time.Parse("2020-02-10", string(t))
}
