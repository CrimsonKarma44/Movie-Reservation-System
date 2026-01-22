package handlers

import (
	"encoding/json"
	"fmt"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	movieService *services.Service[models.Movie]
}

func NewMovieHandler(movieService *services.Service[models.Movie]) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (mh *MovieHandler) AddMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AddMovie Handler...")
	if r.Method == http.MethodPost {
		var movie *models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		err := mh.movieService.Add(movie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/movie/%d", movie.ID), http.StatusSeeOther)
	}
}

func (mh *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateMovie Handler...")
	if r.Method == http.MethodPut {
		var movie *models.Movie
		pathValue := r.PathValue("id")

		id, err := strconv.Atoi(pathValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		
		fmt.Println("reached this point")
		err = mh.movieService.Update(id, movie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		http.Redirect(w, r, fmt.Sprintf("/movie/%d", id), http.StatusSeeOther)
	}
}

func (mh *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteMovie Handler...")
	if r.Method == http.MethodDelete {
		pathValue := r.PathValue("id")

		id, err := strconv.Atoi(pathValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := mh.movieService.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (mh *MovieHandler) GetMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMovie Handler...")
	if r.Method == http.MethodGet {
		var movie *models.Movie
		pathValue := r.PathValue("id")

		id, err := strconv.Atoi(pathValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		movie, err = mh.movieService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(movie); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (mh *MovieHandler) GetAllMoviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var movie []models.Movie

		movie, err := mh.movieService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(movie); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
