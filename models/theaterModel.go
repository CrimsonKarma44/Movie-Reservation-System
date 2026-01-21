package models

import (
	"gorm.io/gorm"
)

type Theater struct {
	gorm.Model

	Name     string
	Capacity int
	Location string

	ShowTimes []ShowTime `gorm:"foreignKey:TheaterID"`
}

func NewTheater(name string, capacity int, location string) Theater {
	return Theater{
		Name:     name,
		Capacity: capacity,
		Location: location,
	}
}

func (theater Theater) Type() string {
	return "Theater"
}
