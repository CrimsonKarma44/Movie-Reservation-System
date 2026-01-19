package services

import (
	// "fmt"
	"movie-reservation-system/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

// func NewUserService(db *gorm.DB) *UserService {
// 	fmt.Println("UserService initialized")
// 	return &UserService{db: db}
// }

func (us UserService) CreateUser(user *models.User) error {
	return us.DB.Create(user).Error
}

func (us UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := us.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us UserService) UpdateUser(user *models.User) error {
	return us.DB.Save(user).Error
}
