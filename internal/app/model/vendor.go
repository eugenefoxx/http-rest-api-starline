package model

// NewVendor The function will return a pointer to the array of type Vendor
func NewVendor() *Vendors {
	return &Vendors{}
}

// Vendor ...
type Vendor struct {
	ID          int    `json:"id"`
	CodeDebitor string `json:"code_debitor"`
	NameDebitor string `json:"name_debitor"`
}

// Vendors ...
type Vendors []Vendor
