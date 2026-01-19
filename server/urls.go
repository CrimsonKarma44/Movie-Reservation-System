package moviereservationsystem

import (
	"fmt"
	"movie-reservation-system/handlers"
	"movie-reservation-system/middleware"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"

	"gorm.io/gorm"
)

type url struct {
	db *gorm.DB

	authHandlers        *handlers.AuthHandler
	movieHandlers       *handlers.MovieHandler
	reservationHandlers *handlers.ReservationHandler
	theaterHandlers     *handlers.TheaterHandler
	showtimeHandlers    *handlers.ShowtimeHandler

	authMiddleware *middleware.AuthMiddleware
}

func NewServer(db *gorm.DB, env *models.Env, refreshStore map[uint]string, accessStore map[uint]string) *url {

	return &url{
		db: db,

		authHandlers:        handlers.NewAuthHandler(services.NewAuthService(db), env, refreshStore, accessStore),
		movieHandlers:       handlers.NewMovieHandler(services.NewMovieService(db)),
		reservationHandlers: handlers.NewReservationHandler(services.NewReservationService(db)),
		theaterHandlers:     handlers.NewTheaterHandler(services.NewTheaterService(db)),
		showtimeHandlers:    handlers.NewShowtimeHandler(services.NewShowtimeService(db)),

		authMiddleware: middleware.NewAuthMiddleware(env, refreshStore, accessStore),
	}

}

func (s *url) Run() {
	fmt.Println("API Server Initializing...")

	// routes
	mux := http.NewServeMux()

	s.auth(mux)
	s.movie(mux)
	s.theater(mux)

	// static Dir
	// mux.Handle("/static/", s.authMiddleware.ProtectMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir(path[1:]))).ServeHTTP))

	http.ListenAndServe(":8080", mux)
	// Implement server run logic here
}

func (s *url) auth(mux *http.ServeMux) {
	mux.Handle("/auth/register", http.HandlerFunc(s.authHandlers.RegisterHandler))
	mux.Handle("/auth/login", http.HandlerFunc(s.authHandlers.LoginHandler))
	mux.Handle("/auth/logout", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.authHandlers.LogoutHandler)))
	mux.Handle("/auth/renew-token", s.authMiddleware.RenewTokenMiddleware(
		func(w http.ResponseWriter, r *http.Request) {},
	))
}

func (s *url) movie(mux *http.ServeMux) {
	mux.Handle("/movies/", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.movieHandlers.GetAllMoviesHandler)))
	mux.Handle("/movie/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.movieHandlers.GetMovieHandler)))

	// admin
	mux.Handle("/movie/add", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.AddMovieHandler))))
	mux.Handle("/movie/update/{id}", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.UpdateMovieHandler))))
	mux.Handle("/movies/delete/{id}", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.DeleteMovieHandler))))
}

func (s *url) theater(mux *http.ServeMux) {
	mux.Handle("/theaters", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.theaterHandlers.GetAllTheaters)))
	mux.Handle("/theater/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.theaterHandlers.GetTheater)))

	// admin
	mux.Handle("/theater/add", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.CreateTheater))))
	mux.Handle("/theater/{id}/update", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.UpdateTheater))))
	mux.Handle("/theater/{id}/delete", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.DeleteTheater))))
}
