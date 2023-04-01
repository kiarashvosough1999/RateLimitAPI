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
