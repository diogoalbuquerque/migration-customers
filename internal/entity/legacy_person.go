package entity

import "time"

type LegacyPerson struct {
	NI        float64    `json:"ni" example:"37930615073"`
	Name      string     `json:"name" example:"Roberto Silva"`
	BirthDate *time.Time `json:"birthDate,omitempty" example:"1987-04-06"`
}

type ChanelLegacyPerson struct {
	LegacyPeople []LegacyPerson
	Err          error
}
