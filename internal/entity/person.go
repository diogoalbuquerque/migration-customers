package entity

import "time"

type Person struct {
	NI         string     `json:"ni" example:"37930615073"`
	Name       string     `json:"name" example:"Roberto Silva"`
	PersonType PersonType `json:"personType" example:"1"`
	BirthDate  *time.Time `json:"birthDate,omitempty" example:"1987-04-06"`
}

type PersonType int

const (
	PF PersonType = iota + 1
	PJ
)
