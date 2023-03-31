package Respositories

import (
	"RateLimitAPI/Models"
)

type UserRepositoryInterface interface {
	UsernameExist(username string) (bool, error)
	FindByUsername(username string) (*Models.UserModel, error)
	Save(user Models.UserModel) error
}
