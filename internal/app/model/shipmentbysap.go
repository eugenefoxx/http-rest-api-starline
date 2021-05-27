package model

import "time"

// NewShipmentbysap The function will return a pointer to the array of type Shipmentbysap
func NewShipmentbysap() *Shipmentbysaps {
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
