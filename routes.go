package main

import (
	"FrontReferralAPI/entity"
	"FrontReferralAPI/referral_code"
	"FrontReferralAPI/repository"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	repo repository.DeviceRepository = repository.NewRepository()
)

// create referral code "/referrals/{device_id}/create"
func ReferralData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	var record entity.Device

	params := mux.Vars(request)               // Get params
	device_id := params["device_id"]          // Get device_id
	referrer_id := "9189C9"                   // Get referrer_id
	unique_id := referral_code.RandomString() // 6 digit random string referral code
	record.DeviceID = device_id               // Serial Number of the device
	record.UniqueID = unique_id               // Referral Code for particular user referral
	record.ReferrerID = referrer_id           // Referred ID

	existing_device := repo.IsExists(device_id) // Check if device already exists
	if existing_device {
		response.WriteHeader(http.StatusNotFound) // Send response
		json.NewEncoder(response).Encode("Device already exists")
	} else {
		existing_referrer, err := repo.FindByReferrer(referrer_id) // Check if referrer already exists
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(err)
			return
		}
		if existing_referrer.UniqueID == referrer_id {
			repo.Update(referrer_id, device_id)
			response.WriteHeader(http.StatusOK) // Send response
			json.NewEncoder(response).Encode(existing_referrer)
		}
		repo.Save(&record)
		response.WriteHeader(http.StatusOK)      // Send response
		json.NewEncoder(response).Encode(record) // Send response
	}
}

// get device by device_id "/referrals/{device_id}"
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
