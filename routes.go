package main

import (
	"FrontReferralAPI/entity"
	"FrontReferralAPI/referral_code"
	"FrontReferralAPI/repository"
	"encoding/json"
	"fmt"
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
	device_id := "00004"
	referrer_id := "C047BA"
	unique_id := referral_code.RandomString() // 6 digit random string referral code
	record.DeviceID = device_id               // Serial Number of the device
	record.UniqueID = unique_id               // Referral Code for particular user referral
	record.ReferrerID = referrer_id           // Referred ID
	// Save all data to firestore
	existing_device := repo.IsExists(device_id)
	if existing_device {
		response.WriteHeader(409) // Send response
		json.NewEncoder(response).Encode("Device already exists")
	} else {
		existing_record, err := repo.FindByReferrer(referrer_id) // Find the record by referral code
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(err)
			return
		}
		if existing_record.UniqueID == referrer_id {
			repo.Update(referrer_id, device_id)
		}
		repo.Save(&record)
		response.WriteHeader(http.StatusOK) // Send response
		json.NewEncoder(response).Encode(record)
	}
}

// get device by device_id
func GetDevice(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	params := mux.Vars(request)
	device_id := params["device_id"]
	fmt.Println("Device ID: ", device_id)
	record, err := repo.FindDevice(device_id)
	fmt.Println("Record: ", record)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(err)
		return
	}
	if record == nil {
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode("Device not found")
		return
	}
	response.WriteHeader(http.StatusOK) // Send response
	json.NewEncoder(response).Encode(record)
}
