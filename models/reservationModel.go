package models

import (
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

// to complete
func (rm *ReservationsModel) BeforeSave(tx *gorm.DB) error {
	return nil
}
