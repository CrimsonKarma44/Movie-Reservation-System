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
		Title:       movie.Title,
		Description: movie.Description,
		Year:        movie.Year,
		Duration:    movie.Duration,
	}
}

func (movie Movie) Type() string {
	return "Movie"
}