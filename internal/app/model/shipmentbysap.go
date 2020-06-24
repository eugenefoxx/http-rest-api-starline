package model

import "time"

// Main The function will return a pointer to the array of type Shipmentbysap
func Main() *Shipmentbysaps {
	return &Shipmentbysaps{}
}

// Shipmentbysap ...
type Shipmentbysap struct {
	Material      int       `json:"material"`
	Qty           int       `json:"qty"`
	ShipmentDate  time.Time `json:"shipment_date"`
	ShipmentDate2 string    `json:"shipment_date2"`
	ShipmentDate3 string    `json:"shipment_date3"`
	ShipmentTime  time.Time `json:"shipment_time"`
	ShipmentTime2 string    `json:"shipment_time2"`
	ID            int       `json:"id"`
	LastName      string    `json:"lastname"`
	Comment       string    `json:"comment"`
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
