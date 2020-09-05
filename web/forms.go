package web

import "encoding/gob"

func init() {
	gob.Register(IDReturn{})
}

type FormErrors map[string]string

type IDReturn struct {
	ScanID     string
	QtyFact    int
	QtySAP     int
	QtyPanacim int

	Errors FormErrors
}

func (f *IDReturn) Validate() bool {
	f.Errors = FormErrors{}
	if f.ScanID == "" {
		f.Errors["ScanID"] = "Пустое поле для сканирования."
	}

	if f.QtyFact < 0 {
		f.Errors["QtyFact"] = "Пустое поле для количества по факту"
	}

	if f.QtySAP < 0 {
		f.Errors["QtySAP"] = "Пустое поле для количества по SAP"
	}

	if f.QtyPanacim < 0 {
		f.Errors["QtyPanacim"] = "Пустое поле для количества по PanaCIM"
	}

	return len(f.Errors) == 0
}
