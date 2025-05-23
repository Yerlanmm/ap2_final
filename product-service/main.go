package main

import (
	"log"
	"net/http"
	"product-service/config"
	"product-service/handler"

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
	
	r.HandleFunc("/main", handler.MainPage).Methods("GET")
	r.HandleFunc("/payment", handler.PaymentPage).Methods("GET")




	// Start the server
	log.Println("✅ User service running on http://localhost:8001...")
	log.Fatal(http.ListenAndServe(":8001", r))
}
