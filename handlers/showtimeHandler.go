package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"
	"strconv"
)

type ShowtimeHandler struct {
	showtimeService *services.Service[models.ShowTime]
}

func NewShowtimeHandler(showtimeService *services.Service[models.ShowTime]) *ShowtimeHandler {
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

		if err := st.showtimeService.Add(showtime); err != nil {
			fmt.Println(err)
			http.Error(w, "failed to create showtime", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/showtime/%d", showtime.ID), http.StatusSeeOther)

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (st *ShowtimeHandler) GetShowtime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		showtimeID := r.PathValue("id")

		id, err := strconv.Atoi(showtimeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		showtime, err := st.showtimeService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(showtime)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (st *ShowtimeHandler) GetAllShowtimes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		showtimes, err := st.showtimeService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(showtimes)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (st *ShowtimeHandler) UpdateShowtime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Showtime Handler...")
	if r.Method == http.MethodPut {
		var showtime *models.ShowTime
		showtimeID := r.PathValue("id")

		id, err := strconv.Atoi(showtimeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewDecoder(r.Body).Decode(&showtime); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := st.showtimeService.Update(id, showtime); err != nil {
			fmt.Println(err)
			http.Error(w, "failed to update showtime", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/showtime/%d", id), http.StatusSeeOther)

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (st *ShowtimeHandler) DeleteShowtime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Showtime Handler...")
	if r.Method == http.MethodDelete {
		showtimeID := r.PathValue("id")

		id, err := strconv.Atoi(showtimeID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := st.showtimeService.Delete(id); err != nil {
			fmt.Println(err)
			http.Error(w, "failed to delete showtime", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(map[string]string{"message": "showtime deleted successfully"})
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}

}
