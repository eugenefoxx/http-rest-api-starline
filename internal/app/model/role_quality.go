package model

func NewRole_Quality() *RoleQualitys {
	return &RoleQualitys{}
}

// Role_Quality model base
type Role_Quality struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

type RoleQualitys []Role_Quality
