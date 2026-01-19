package models

import (
	// "movie-reservation-system/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email    string `json:"email" gorm:"unique"`
	UserName string `json:"username"`
	Password string `json:"password"`

	IsAdmin bool `json:"omitempty"`

	Reservations []Reservation
}

func (admin *User) UserPromotion(user *User) error {
	if admin.IsAdmin {
		if user.IsAdmin {
			return ErrAlreadyAdmin{}
		}

		user.IsAdmin = true
		return nil
	}
	return ErrNotAdmin{}
}

func (admin *User) UserDemotion(user *User) error {
	if admin.IsAdmin {
		if !user.IsAdmin {
			return ErrNotAdmin{}
		}

		user.IsAdmin = false
		return nil
	}
	return ErrNotAdmin{}
}
