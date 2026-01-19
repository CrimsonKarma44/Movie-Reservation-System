package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"
)

type ShowtimeHandler struct {
	showtimeService *services.ShowtimeService
}

func NewShowtimeHandler(showtimeService *services.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{
		showtimeService: showtimeService,
	}
}

func (st *ShowtimeHandler) CreateShowtime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Showtime Handler...")
	if r.Method == http.MethodPost {
		var showtime *models.ShowTime
		if err := json.NewDecoder(r.Body).Decode(&showtime); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := st.showtimeService.CreateShowtime(showtime); err != nil {
			http.Error(w, "failed to create showtime", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(showtime)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}

}
