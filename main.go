package main

import (
	"RateLimitAPI/Controllers"
	"RateLimitAPI/DB"
	"log"
	"net/http"
)

func main() {

	migrator, dbErr := DB.New()

	migrator.AutoMigrateModels()

	if dbErr != nil {
		log.Print("database has error")
	}

	signupController := Controllers.SignUpController{
		UserRepository: migrator,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	mux.HandleFunc("/signup", signupController.Handler)

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
