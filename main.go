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

package main

import (
	"RateLimitAPI/Controllers"
	"RateLimitAPI/DB"
	"RateLimitAPI/Middlewares"
	"log"
	"net/http"
)

func main() {

	migrator, dbErr := DB.New()
	migrator.AutoMigrateModels()
	if dbErr != nil {
		log.Print("database has error")
	}

	signupController := Controllers.NewSignUpController(migrator)
	signinController := Controllers.NewSignInController(migrator)

	limiter := Middlewares.NewRateLimiter(20, 1)

	commonMiddleware := []Middlewares.Middleware{
		limiter.Limiter,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", Middlewares.MultipleMiddleware(okHandler, commonMiddleware...))
	mux.HandleFunc("/signup", Middlewares.MultipleMiddleware(signupController.Handler, commonMiddleware...))
	mux.HandleFunc("/signin", Middlewares.MultipleMiddleware(signinController.Handler, commonMiddleware...))

	// Wrap the servemux with the limit middleware.
	log.Print("Listening on :4000...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Print("Failed to Serve")
	}
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
