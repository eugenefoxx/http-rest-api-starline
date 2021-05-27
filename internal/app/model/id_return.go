package model

// NewIDReturn The function will return a pointer to the array of type IDReturn
func NewIDReturn() *IDReturns {
	return &IDReturns{}
}

// IDReturn ...
type IDReturn struct {
	Material   int    `json:"material"`
	IDRoll     int    `json:"idroll"`
	Lot        string `json:"lot"`
	QtyFact    int    `json:"qtyfact"`
	QtySAP     int    `json:"qtysap"`
	QtyPanacim int    `json:"qtypanacim"`
	ID         int    `json:"id"`
	LastName   string `json:"lastname"`
}

// IDReturns ...
type IDReturns []IDReturn
