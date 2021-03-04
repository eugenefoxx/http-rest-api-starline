package model

import "time"

// MainInspection The function will return a pointer to the array of type Inspection
func MainInspection() *Inspections {
	return &Inspections{}
}

// Inspection
type Inspection struct {
	ID             int       `json:"id"`
	IdMaterial     string    `json:"idmaterial"`
	SAP            int       `json:"sap"`
	Lot            string    `json:"lot"`
	IdRoll         int       `json:"idroll"`
	Qty            int       `json:"qty"` // initial qty
	ProductionDate string    `json:"productiondate"`
	NumberVendor   string    `json:"numbervendor"`
	NameDebitor    string    `json:"name_debitor"`
	Note           string    `json:"note"`
	Location       string    `json:"location"`
	Status         string    `json:"status"`
	Date           time.Time `json:"date"`
	Date2          string    `json:"date2"`
	Time           time.Time `json:"time"`
	Time2          string    `json:"time2"`
	Lastname       string    `json:"lastname"`
	Role           string    `json:"role"`
	Groups         string    `json:"groups"`
	Update         string    `json:"update"`
	Dateupdate     time.Time `json:"dateupdate"`
	Dateupdate2    string    `json:"dateupdate2"`
	Timeupdate     time.Time `json:"timeupdate"`
	Timeupdate2    string    `json:"timeupdate2"`
	Lastnameaccept string    `json:"lastnameaccept"`
	Dateaccept     time.Time `json:"dateaccept"`
	Dateaccept2    string    `json:"dateaccept2"`
	Timeaccept     time.Time `json:"timeaccept"`
	Timeaccept2    string    `json:"timeaccept2"`
	Count          int       `json:"count"`
}

// Inspections ...
type Inspections []Inspection
