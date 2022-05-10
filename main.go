package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	port := os.Getenv("PORT")

	router.HandleFunc("/", HomePage)

	router.HandleFunc("/api/v1/referral", ReferralData).Methods("POST")            // Create a new referral code
	router.HandleFunc("/api/v1/referral", GetDevice).Methods("GET")                // Get device by device_id
	router.HandleFunc("/api/v1/referrals", GetAllDevices).Methods("GET")           // Get all devices
	router.HandleFunc("/api/v1/referral/counts", GetReferredCounts).Methods("GET") // Get referral counts
	// log.Println("Server listening on port", port)
	// log.Fatalln(http.ListenAndServe(port, router))
	// log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router)) // for production
	if len(port) == 0 {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
