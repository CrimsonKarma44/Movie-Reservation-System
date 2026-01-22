package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"
	"strconv"
)

type TheaterHandler struct {
	theaterService *services.Service[models.Theater]
}

func NewTheaterHandler(theaterService *services.Service[models.Theater]) *TheaterHandler {
	return &TheaterHandler{
		theaterService: theaterService,
	}
}

func (th *TheaterHandler) CreateTheater(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Theater Handler...")
	if r.Method == http.MethodPost {
		var theater *models.Theater
		if err := json.NewDecoder(r.Body).Decode(&theater); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := th.theaterService.Add(theater); err != nil {
			http.Error(w, "failed to create theater", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/theater/%d", theater.ID), http.StatusSeeOther)
	}
}

func (th *TheaterHandler) GetTheater(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Theater Handler...")
	if r.Method == http.MethodGet {
		var theater *models.Theater
		theaterID := r.PathValue("id")

		id, err := strconv.Atoi(theaterID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		theater, err = th.theaterService.GetByID(id)
		if err != nil {
			http.Error(w, "failed to get theater", http.StatusInternalServerError)
			return
		}

		// w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(theater)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
func (th *TheaterHandler) GetAllTheaters(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get All Theaters Handler...")
	if r.Method == http.MethodGet {
		theaters, err := th.theaterService.GetAll()
		if err != nil {
			http.Error(w, "failed to get theater", http.StatusInternalServerError)
			return
		}

		// w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(theaters)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (th *TheaterHandler) UpdateTheater(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Theater Handler...")
	if r.Method == http.MethodPut {
		var theater *models.Theater
		theaterID := r.PathValue("id")

		id, err := strconv.Atoi(theaterID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewDecoder(r.Body).Decode(&theater); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := th.theaterService.Update(id, theater); err != nil {
			http.Error(w, "failed to update theater", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/theater/%d", id), http.StatusSeeOther)
	}
}

func (th *TheaterHandler) DeleteTheater(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Theater Handler...")
	if r.Method == http.MethodDelete {
		theaterID := r.PathValue("id")

		id, err := strconv.Atoi(theaterID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := th.theaterService.Delete(id); err != nil {
			http.Error(w, "failed to delete theater", http.StatusInternalServerError)
			return
		}

		// w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(map[string]string{"message": "Theater deleted successfully"})
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
