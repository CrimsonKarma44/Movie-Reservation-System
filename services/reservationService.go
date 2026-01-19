package services

import (
	"movie-reservation-system/models"

	"gorm.io/gorm"
)

type ReservationService struct {
	db *gorm.DB
}

func NewReservationService(db *gorm.DB) *ReservationService {
	return &ReservationService{
		db: db,
	}
}

func (rs *ReservationService) CreateReservation(reservation *models.Reservation, showtime *models.ShowTime) (*models.Reservation, error) {
	err := showtime.SeatReserveAvailability(reservation.Seats)
	if err != nil {
		return nil, err
	}
	if err := rs.db.Create(&reservation).Error; err != nil {
		return nil, err
	}
	return reservation, nil
}

func (rs *ReservationService) GetReservationByID(id int) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := rs.db.First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (rs *ReservationService) CancelReservation(reservation *models.Reservation, showtime *models.ShowTime) error {
	showtime.SeatReturn(reservation.Seats)

	if err := rs.db.Delete(&reservation).Error; err != nil {
		return err
	}
	return nil
}
