package Controllers

import (
	"RateLimitAPI/Helpers"
	"RateLimitAPI/Models"
	"RateLimitAPI/Respositories"
	"errors"
	"log"
	"net/http"
)

type signUpBody struct {
	Username string
	Password string
}

type SignUpController struct {
	UserRepository Respositories.UserRepositoryInterface
}

func NewSignUpController(
	UserRepository Respositories.UserRepositoryInterface,
) *SignUpController {
	return &SignUpController{
		UserRepository,
	}
}

func (c *SignUpController) Handler(w http.ResponseWriter, r *http.Request) {
	var p signUpBody

	err := Helpers.DecodeJSONBody(w, r, &p)
	if err != nil {
		var mr *Helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	exist, dbErr := c.UserRepository.UsernameExist(p.Username)
	if dbErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if exist {
		http.Error(w, "This username already exist", http.StatusForbidden)
	} else {
		createErr := c.UserRepository.Save(Models.UserModel{
			Username: p.Username,
			Password: p.Password,
		})

		if createErr != nil {
			http.Error(w, "Can not signup", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Signed up successfully"))
		}
	}
}
