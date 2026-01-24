package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model

	UserID     int
	ShowTimeID int

	Seats int
	Cost  int

	expiryStatus bool
}

func NewReservation(userID int, showTimeID int, seats int, cost int) Reservation {
	return Reservation{
		UserID:     userID,
		ShowTimeID: showTimeID,
		Seats:      seats,
		Cost:       cost,
	}
}

func (rm Reservation) Type() string {
	return "Reservation"
}

func (rm *Reservation) SetExpired() {
	rm.expiryStatus = true
}

func (rm *Reservation) BeforeSave(tx *gorm.DB) error {
	if err := tx.First(&User{}, rm.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("dependent Record not found: %v", err)
		}
	}

	var showtime ShowTime
	if err := tx.First(&showtime, rm.ShowTimeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("dependent Record not found: %v", err)
		}
	}

	rm.Cost = rm.Seats * showtime.Price

	return nil
}

// kept aside for now
type ReservationsModel struct {
	reservations []Reservation
}

func NewReservationModel() *ReservationsModel {
	return &ReservationsModel{
		reservations: []Reservation{},
	}
}

func (rm *ReservationsModel) Create(reservation Reservation) {
	rm.reservations = append(rm.reservations, reservation)
}

func (rm *ReservationsModel) Get() []Reservation {
	return rm.reservations
}
