package Respositories

import (
	"RateLimitAPI/Models"
)

// UserRepositoryInterface create a layer of abstraction to access `Store`
// in which we store user data or fetch it.
type UserRepositoryInterface interface {
	UsernameExist(username string) (bool, error)
	FindByUsername(username string) (*Models.UserModel, error)
	Save(user Models.UserModel) error
}
