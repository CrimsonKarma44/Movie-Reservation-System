package handlers

import "movie-reservation-system/services"

type ReservationHandler struct {
	reservationService *services.ReservationService
}

func NewReservationHandler(reservationService *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}
