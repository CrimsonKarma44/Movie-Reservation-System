package moviereservationsystem

import (
	"errors"
	"fmt"
	"log"
	"movie-reservation-system/models"
	"movie-reservation-system/services"
	"movie-reservation-system/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	// defaultAdmin models.User
}

func (app *App) Init() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}
	var ENV = models.NewEnv()

	// Initialize thread-safe refresh token store
	refreshStore := models.NewSafeTokenStore()

	// starts the Database server
	db, err := services.NewDatabaseService(app.dns(ENV))
	if err != nil {
		return fmt.Errorf("Error initializing database")
	}

	// performst the initial migrations
	db.Migrate()

	// creates the default admin user
	app.CreateDefaultAdmin(db.DB, ENV)

	// starts the server
	server := NewServer(db.DB, ENV, refreshStore)
	server.Run()

	return nil
}

// func (app *App) Init(refreshStore map[uint]string) error {
// 	err := godotenv.Load()
// 	if err != nil {
// 		return fmt.Errorf("Error loading .env file")
// 	}
// 	var ENV = models.NewEnv()

// 	// starts the Database server
// 	db, err := services.NewDatabaseService(app.dns(ENV))
// 	if err != nil {
// 		return fmt.Errorf("Error initializing database")
// 	}

// 	// performst the initial migrations
// 	db.Migrate()

// 	// creates the default admin user
// 	app.CreateDefaultAdmin(db.DB, ENV)

// 	// starts the server
// 	server := NewServer(db.DB, ENV, refreshStore)
// 	server.Run()

// 	return nil
// }

func (app *App) dns(env *models.Env) string {
	if env.DBHost == "" || env.DBPassword == "" || env.DBName == "" {
		log.Fatalf("Missing required environment variables")
	}

	if env.DBPort == "" {
		env.DBPort = "5432"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.DBHost,
		env.DBUser,
		env.DBPassword,
		env.DBName,
		env.DBPort,
	)
}

func (app *App) CreateDefaultAdmin(db *gorm.DB, env *models.Env) (models.User, error) {
	var admin models.User

	// Check if admin already exists
	err := db.Where("is_admin = ?", true).First(&admin).Error
	if err == nil {
		// Admin exists → do nothing
		return admin, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other DB error
		return models.User{}, err
	}

	// Read env vars
	email := env.AdminEmail
	password := env.AdminPassword
	name := env.AdminUsername

	if email == "" || password == "" {
		return models.User{}, errors.New("admin credentials not set in env")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	admin = models.User{
		UserName: name,
		Email:    email,
		Password: hashedPassword,
		IsAdmin:  true,
	}

	err = db.Create(&admin).Error
	if err != nil {
		return models.User{}, err
	}

	return admin, nil
}
