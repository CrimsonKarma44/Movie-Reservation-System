package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model

	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Year        int    `gorm:"not null"`
	Duration    int    `gorm:"not null"`
	// ShowTime    time.Time `gorm:"not null"`

	ShowTime *[]ShowTime `gorm:"foreignKey:MovieID"`
}

func (movie Movie) New(title string, description string, year int, duration int) Movie {
	return Movie{
		Title:       title,
		Description: description,
		Year:        year,
		Duration:    duration,
	}
}

func (movie Movie) Type() string {
	return "Movie"
}

func (movie *Movie) Preload(db *gorm.DB) *gorm.DB {
	return db.Preload("ShowTime").Preload("ShowTime.Theater")
}