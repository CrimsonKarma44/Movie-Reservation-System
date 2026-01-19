package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ShowTime struct {
	gorm.Model
	Price          int
	StartTime      time.Time
	AvailableSeats int

	MovieID   int
	TheaterID int

	Theater      *Theater      `gorm:"foreignKey:TheaterID"`
	Reservations []Reservation `gorm:"foreignKey:ShowTimeID"`
}

func NewShowTime(price int, startTime time.Time, availableSeats int, movieID int, theaterID int) ShowTime {
	return ShowTime{
		Price:          price,
		StartTime:      startTime,
		AvailableSeats: availableSeats,
		MovieID:        movieID,
		TheaterID:      theaterID,
	}
}

func (st *ShowTime) BeforeSave(tx *gorm.DB) error {
	// Load theater if not already loaded
	if st.Theater == nil && st.TheaterID != 0 {
		var th Theater
		if err := tx.First(&th, st.TheaterID).Error; err != nil {
			return err
		}
		st.Theater = &th
	}

	// Now we can safely reference Capacity
	if st.Theater == nil {
		return fmt.Errorf("theater not found")
	}
	if st.AvailableSeats < 0 || st.AvailableSeats > st.Theater.Capacity {
		return fmt.Errorf("invalid available seats")
	}

	if st.StartTime.Before(time.Now()) {
		return fmt.Errorf("invalid start time")
	}
	return nil
}

func (st *ShowTime) SeatReserveAvailability(seats int) error {
	if seats > st.AvailableSeats {
		return fmt.Errorf("not enough seats available")
	} else {
		st.AvailableSeats -= seats
	}
	return nil
}

func (st *ShowTime) SeatReturn(seats int) {
	st.AvailableSeats += seats
}
