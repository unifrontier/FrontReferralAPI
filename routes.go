package main

import (
	"FrontReferralAPI/entity"
	"FrontReferralAPI/referral_code"
	"FrontReferralAPI/repository"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	repo repository.DeviceRepository = repository.NewRepository()
)

func ReferralData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")

	rand.Seed(time.Now().UnixNano())
	var record entity.Device
	device_id := "00001"
	referrer_id := "E31BCF"
	unique_id := referral_code.RandomString() // 6 digit random string referral code
	record.DeviceID = device_id               // Serial Number of the device
	record.UniqueID = unique_id               // Referral Code for particular user referral
	record.ReferrerID = referrer_id           // Referred ID
	// Save all data to firestore
	existing_record, err := repo.Find(referrer_id) // Find the record by referral code
	if err != nil {
		log.Println(err)
	}
	fmt.Println(existing_record)
	if existing_record.UniqueID == referrer_id {
		fmt.Println("Record exists")
		repo.Update(referrer_id, device_id)
		// existing_record.ReferredIDS = append(existing_record.ReferredIDS, device_id)
		fmt.Println(existing_record.ReferredIDS)
	}
	repo.Save(&record)
	response.WriteHeader(http.StatusOK) // Send response
	json.NewEncoder(response).Encode(record)
}

func GetData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	params := mux.Vars(request)
	device_id := params["device_id"]
	record, err := repo.Find(device_id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(err)
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(record)
}
