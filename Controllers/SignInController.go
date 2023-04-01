package Controllers

import (
	"RateLimitAPI/Helpers"
	"RateLimitAPI/Respositories"
	"errors"
	"log"
	"net/http"
)

type signInBody struct {
	Username string
	Password string
}

type SignInController struct {
	UserRepository Respositories.UserRepositoryInterface
}

func NewSignInController(
	UserRepository Respositories.UserRepositoryInterface,
) *SignInController {
	return &SignInController{
		UserRepository,
	}
}

func (c *SignInController) Handler(w http.ResponseWriter, r *http.Request) {
	var p signInBody

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

	userModel, dbErr := c.UserRepository.FindByUsername(p.Username)
	if dbErr != nil {
		http.Error(w, "Username Or Password Incorrect", http.StatusInternalServerError)
		return
	}
	if userModel != nil && userModel.Password == p.Password {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Signed in successfully"))
	} else {
		http.Error(w, "Username Or Password Incorrect", http.StatusInternalServerError)
	}
}
