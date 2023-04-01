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

	limiter := Middlewares.NewRateLimiter(20, 1)

	commonMiddleware := []Middlewares.Middleware{
		limiter.Limiter,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", Middlewares.MultipleMiddleware(okHandler, commonMiddleware...))
	mux.HandleFunc("/signup", Middlewares.MultipleMiddleware(signupController.Handler, commonMiddleware...))

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
