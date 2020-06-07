package model

import "time"

// Main The function will return a pointer to the array of type Shipmentbysap
func Main() *Shipmentbysaps {
	return &Shipmentbysaps{}
}

// Shipmentbysap ...
type Shipmentbysap struct {
	Material      int       `db:"material"`
	Qty           int       `db:"qty"`
	ShipmentDate  time.Time `db:"shipment_date"`
	ShipmentDate2 string    `db:"shipment_date"`
	ShipmentTime  time.Time `db:"shipment_time"`
	ShipmentTime2 string    `db:"shipment_time"`
	ID            int       `db:"id"`
	LastName      string    `db:"lastname"`
	Comment       string    `db:"comment"`
}

// Shipmentbysaps ...
type Shipmentbysaps []Shipmentbysap

type rawTime []byte

func (t rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type ShipmentDate []byte

func (t ShipmentDate) Time() (time.Time, error) {
	return time.Parse("2020-02-10", string(t))
}
