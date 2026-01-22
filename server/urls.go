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

	authHandlers *handlers.AuthHandler

	movieHandlers       *handlers.MovieHandler
	reservationHandlers *handlers.ReservationHandler
	theaterHandlers     *handlers.TheaterHandler
	showtimeHandlers    *handlers.ShowtimeHandler

	authMiddleware *middleware.AuthMiddleware
}

func NewServer(db *gorm.DB, env *models.Env, refreshStore map[uint]string) *url {

	return &url{
		db: db,

		authHandlers: handlers.NewAuthHandler(services.NewAuthService(db), env, refreshStore),

		movieHandlers:       handlers.NewMovieHandler(services.NewService(db, models.Movie{})),
		theaterHandlers:     handlers.NewTheaterHandler(services.NewService(db, models.Theater{})),
		reservationHandlers: handlers.NewReservationHandler(services.NewService(db, models.Reservation{})),
		showtimeHandlers:    handlers.NewShowtimeHandler(services.NewService(db, models.ShowTime{})),

		authMiddleware: middleware.NewAuthMiddleware(env, refreshStore),
	}

}

func (s *url) Run() {
	fmt.Println("API Server Initializing...")

	// routes
	mux := http.NewServeMux()

	s.auth(mux)
	s.movie(mux)
	s.theater(mux)
	s.showtime(mux)
	s.reservation(mux)

	http.ListenAndServe(":8080", mux)
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
	mux.Handle("/movie/{id}/update", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.UpdateMovieHandler))))
	mux.Handle("/movie/{id}/delete", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.DeleteMovieHandler))))
}

func (s *url) theater(mux *http.ServeMux) {
	mux.Handle("/theaters", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.theaterHandlers.GetAllTheaters)))
	mux.Handle("/theater/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.theaterHandlers.GetTheater)))

	// admin
	mux.Handle("/theater/add", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.CreateTheater))))
	mux.Handle("/theater/{id}/update", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.UpdateTheater))))
	mux.Handle("/theater/{id}/delete", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.DeleteTheater))))
}

func (s *url) showtime(mux *http.ServeMux) {
	mux.Handle("/showtimes", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.showtimeHandlers.GetAllShowtimes)))
	mux.Handle("/showtime/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.showtimeHandlers.GetShowtime)))

	// admin
	mux.Handle("/showtime/add", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.showtimeHandlers.CreateShowtime))))
	mux.Handle("/showtime/{id}/update", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.showtimeHandlers.UpdateShowtime))))
	mux.Handle("/showtime/{id}/delete", s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.showtimeHandlers.DeleteShowtime))))
}

func (s *url) reservation(mux *http.ServeMux) {
	mux.Handle("/reservation/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.GetReservation)))
	mux.Handle("/reservations", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.GetAllReservations)))
	mux.Handle("/reservation/add", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.CreateReservation)))
	mux.Handle("/reservation/{id}/update", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.UpdateReservation)))
	mux.Handle("/reservation/{id}/delete", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.DeleteReservation)))
}
