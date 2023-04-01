//  Copyright 2023 KiarashVosough and other contributors
//
//  Permission is hereby granted, free of charge, to any person obtaining
//  a copy of this software and associated documentation files (the
//  Software"), to deal in the Software without restriction, including
//  without limitation the rights to use, copy, modify, merge, publish,
//  distribute, sublicense, and/or sell copies of the Software, and to
//  permit persons to whom the Software is furnished to do so, subject to
//  the following conditions:
//
//  The above copyright notice and this permission notice shall be
//  included in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
//  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
//  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
//  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
