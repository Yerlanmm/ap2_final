package main

import (
	"log"
	"net/http"
	"user-service/config"
	"user-service/handler"

	"github.com/gorilla/mux"
)

func main() {
	// Connect to MongoDB
	config.ConnectDB()

	// Create new router
	r := mux.NewRouter()

	// Serve static files (CSS, JS, images)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve HTML pages (GET)
	r.HandleFunc("/register", handler.RegisterForm).Methods("GET")
	r.HandleFunc("/login", handler.LoginForm).Methods("GET")
	r.HandleFunc("/logout", handler.Logout).Methods("GET")



	// Handle form submissions (POST)
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")

	// Start the server
	log.Println("✅ User service running on http://localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

