package Models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Username string
	Password string
}
