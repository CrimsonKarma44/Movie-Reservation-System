package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"movie-reservation-system/utils"
	"net/http"
	"strconv"
)

type ReservationHandler struct {
	reservationService *services.Service[models.Reservation]
}

func NewReservationHandler(reservationService *services.Service[models.Reservation]) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}

func (st *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Reservation Handler...")
	if r.Method == http.MethodPost {
		claims, ok := r.Context().Value(utils.UserContextKey).(*models.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var reservation *models.Reservation
		if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		reservation.UserID = int(claims.ID)

		if err := st.reservationService.Add(reservation); err != nil {
			fmt.Println(err)
			http.Error(w, "failed to create showtime", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/reservation/%d", reservation.ID), http.StatusSeeOther)

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (rh *ReservationHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		claims, ok := r.Context().Value(utils.UserContextKey).(*models.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		reservationID := r.PathValue("id")

		id, err := strconv.Atoi(reservationID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// if !claims.IsAdmin {
		// 	if id != int(claims.ID) {
		// 		http.Error(w, "operation not allowed for this user", http.StatusUnauthorized)
		// 	}
		// }

		reservation, err := rh.reservationService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(reservation)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (rh *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		claims, ok := r.Context().Value(utils.UserContextKey).(*models.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		reservations, err := rh.reservationService.GetAllByID(int(claims.ID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(reservations)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (rh *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Reservation Handler...")
	if r.Method == http.MethodPut {
		var reservation *models.Reservation
		reservationID := r.PathValue("id")

		id, err := strconv.Atoi(reservationID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := rh.reservationService.Update(id, reservation); err != nil {
			fmt.Println(err)
			http.Error(w, "failed to update reservation", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/reservation/%d", id), http.StatusSeeOther)

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (rh *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Reservation...")
	if r.Method == http.MethodDelete {
		reservationID := r.PathValue("id")

		id, err := strconv.Atoi(reservationID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := rh.reservationService.Delete(id); err != nil {
			fmt.Println(err)
			http.Error(w, "failed to delete reservation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(map[string]string{"message": "reservation deleted successfully"})
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}

}
