package moviereservationsystem

import (
	"fmt"
	"log"
	"movie-reservation-system/handlers"
	"movie-reservation-system/middleware"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"net/http"
	"os"

	// _ "net/http/pprof"
	// +

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

	// Rate limiters for sensitive endpoints
	authRateLimiter        *middleware.RateLimiter // Auth endpoints: 5 req/sec, burst 20
	adminRateLimiter       *middleware.RateLimiter // Admin endpoints: 10 req/sec, burst 50
	reservationRateLimiter *middleware.RateLimiter // Reservation endpoints: 15 req/sec, burst 100
}

func NewServer(db *gorm.DB, env *models.Env, refreshStore *models.SafeTokenStore) *url {
	return &url{
		db: db,

		authHandlers: handlers.NewAuthHandler(services.NewAuthService(db), env, refreshStore),

		movieHandlers:       handlers.NewMovieHandler(services.NewService(db, models.Movie{})),
		theaterHandlers:     handlers.NewTheaterHandler(services.NewService(db, models.Theater{})),
		reservationHandlers: handlers.NewReservationHandler(services.NewService(db, models.Reservation{})),
		showtimeHandlers:    handlers.NewShowtimeHandler(services.NewService(db, models.ShowTime{})),

		authMiddleware: middleware.NewAuthMiddleware(env, refreshStore),

		// Initialize rate limiters
		authRateLimiter:        middleware.NewRateLimiter(5.0, 20),   // 5 requests/sec, burst of 20
		adminRateLimiter:       middleware.NewRateLimiter(10.0, 50),  // 10 requests/sec, burst of 50
		reservationRateLimiter: middleware.NewRateLimiter(15.0, 100), // 15 requests/sec, burst of 100
	}
}


func (s *url) Run() {
	fmt.Println("API Server Initializing...")
	// log.Printf("Access pprof at: http://localhost%s/debug/pprof/", ":8080")

	// routes
	mux := http.NewServeMux()

	// SECURITY FIX: Apply CORS middleware
	corsMiddleware := middleware.NewCORSMiddleware()
	corsMiddleware.LogCORSInfo()

	// SECURITY FIX: Apply security headers middleware
	securityHeaders := middleware.NewSecurityHeadersMiddleware()
	securityHeaders.LogSecurityHeaders()

	s.auth(mux)
	s.movie(mux)
	s.theater(mux)
	s.showtime(mux)
	s.reservation(mux)

	// Use PORT environment variable (Vercel sets this), fallback to 8080 for local development
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server listening on %s\n", addr)

	// Wrap mux with security and CORS middleware
	handler := corsMiddleware.Handler()(
		securityHeaders.Handler()(mux),
	)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func (s *url) auth(mux *http.ServeMux) {
	// Auth endpoints with rate limiting (sensitive - brute force protection)
	mux.Handle("/auth/register", middleware.RateLimitMiddleware(s.authRateLimiter)(http.HandlerFunc(s.authHandlers.RegisterHandler)))
	mux.Handle("/auth/login", middleware.RateLimitMiddleware(s.authRateLimiter)(http.HandlerFunc(s.authHandlers.LoginHandler)))
	mux.Handle("/auth/logout", middleware.RateLimitMiddleware(s.authRateLimiter)(s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.authHandlers.LogoutHandler))))
	mux.Handle("/auth/renew-token", middleware.RateLimitMiddleware(s.authRateLimiter)(s.authMiddleware.RenewTokenMiddleware(
		func(w http.ResponseWriter, r *http.Request) {},
	)))
}

func (s *url) movie(mux *http.ServeMux) {
	mux.Handle("/movies/", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.movieHandlers.GetAllMoviesHandler)))
	mux.Handle("/movie/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.movieHandlers.GetMovieHandler)))

	// admin endpoints with rate limiting (sensitive - prevent abuse)
	mux.Handle("/movie/add", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.AddMovieHandler)))))
	mux.Handle("/movie/{id}/update", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.UpdateMovieHandler)))))
	mux.Handle("/movie/{id}/delete", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.movieHandlers.DeleteMovieHandler)))))
}

func (s *url) theater(mux *http.ServeMux) {
	mux.Handle("/theaters", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.theaterHandlers.GetAllTheaters)))
	mux.Handle("/theater/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.theaterHandlers.GetTheater)))

	// admin endpoints with rate limiting (sensitive - prevent abuse)
	mux.Handle("/theater/add", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.CreateTheater)))))
	mux.Handle("/theater/{id}/update", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.UpdateTheater)))))
	mux.Handle("/theater/{id}/delete", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.theaterHandlers.DeleteTheater)))))
}

func (s *url) showtime(mux *http.ServeMux) {
	mux.Handle("/showtimes", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.showtimeHandlers.GetAllShowtimes)))
	mux.Handle("/showtime/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.showtimeHandlers.GetShowtime)))

	// admin endpoints with rate limiting (sensitive - prevent abuse)
	mux.Handle("/showtime/add", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.showtimeHandlers.CreateShowtime)))))
	mux.Handle("/showtime/{id}/update", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.showtimeHandlers.UpdateShowtime)))))
	mux.Handle("/showtime/{id}/delete", middleware.RateLimitMiddleware(s.adminRateLimiter)(s.authMiddleware.ProtectMiddleware(middleware.AdminMiddleware(http.HandlerFunc(s.showtimeHandlers.DeleteShowtime)))))
}

func (s *url) reservation(mux *http.ServeMux) {
	mux.Handle("/reservation/{id}", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.GetReservation)))
	mux.Handle("/reservations", s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.GetAllReservations)))

	// Reservation modification endpoints with rate limiting (sensitive - prevent spam/abuse)
	mux.Handle("/reservation/add", middleware.RateLimitMiddleware(s.reservationRateLimiter)(s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.CreateReservation))))
	mux.Handle("/reservation/{id}/update", middleware.RateLimitMiddleware(s.reservationRateLimiter)(s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.UpdateReservation))))
	mux.Handle("/reservation/{id}/delete", middleware.RateLimitMiddleware(s.reservationRateLimiter)(s.authMiddleware.ProtectMiddleware(http.HandlerFunc(s.reservationHandlers.DeleteReservation))))
}
