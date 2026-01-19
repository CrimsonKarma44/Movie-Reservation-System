package services

import (
	"errors"
	"fmt"
	"movie-reservation-system/models"
	"time"

	"gorm.io/gorm"
)

type MovieService struct {
	db *gorm.DB
}

func NewMovieService(db *gorm.DB) *MovieService {
	fmt.Println("MovieService initialized")
	return &MovieService{db: db}
}

func (ms *MovieService) AddMovie(movie *models.Movie) error {
	return ms.db.Create(&movie).Error
}

func (ms *MovieService) GetMovieByID(id int) (*models.Movie, error) {
    var movie models.Movie
    if err := ms.db.Preload("ShowTime").Preload("ShowTime.Theater").First(&movie, id).Error; err != nil {
        return nil, err
    }
    return &movie, nil
}

func (ms *MovieService) GetMovies() ([]models.Movie, error) {
	var movies []models.Movie
	if err := ms.db.Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}
func (ms *MovieService) GetMoviesByShowTime(date time.Time) ([]models.Movie, error) {
	var movies []models.Movie
	if err := ms.db.Where("show_time = ?", date).Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (ms *MovieService) UpdateMovie(id int, updateMovie *models.Movie) error {
	movie, err := ms.GetMovieByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return fmt.Errorf("ID %d retrival error: %w", id, err)
	}

	if err := ms.db.Model(movie).Updates(updateMovie).Error; err != nil {
		return fmt.Errorf("failed to update movie with id %d: %w", id, err)
	}

	return nil
}

func (ms *MovieService) DeleteMovie(id int) error {
	movie, err := ms.GetMovieByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("record with id %d not found", id)
		}
		return fmt.Errorf("ID %d retrival error: %w", id, err)
	}

	if err := ms.db.Delete(movie).Error; err != nil {
		return fmt.Errorf("failed to delete movie with id %d: %w", id, err)
	}

	return nil
}
