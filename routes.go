package main

import (
	"FrontReferralAPI/entity"
	"FrontReferralAPI/referral_code"
	"FrontReferralAPI/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	repo repository.DeviceRepository = repository.NewRepository()
)

// Get referral code "/referrals/{device_id}/create/{referral_code}"
func ReferralData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	var record entity.Device
	device_id := request.FormValue("device_id")       // Get device_id from form
	referrer_id := request.FormValue("referral_code") // Get referral_code from form
	if device_id == "" || referrer_id == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode("Device id or referral code is empty")
		return
	}
	unique_id := referral_code.RandomString() // 6 digit random string referral code
	record.DeviceID = device_id               // Serial Number of the device
	record.UniqueID = unique_id               // Referral Code for particular user referral
	record.ReferrerID = referrer_id           // Referred ID

	existing_device, err := repo.FindDevice(device_id)
	if err != nil {
		fmt.Println(err)
	}
	if existing_device != nil && existing_device.DeviceID == device_id {
		response.WriteHeader(http.StatusOK) // Send response
		json.NewEncoder(response).Encode("Device already exists")
		json.NewEncoder(response).Encode(record)
	} else {
		existing_referrer, err := repo.FindByReferrer(referrer_id) // Check if referrer already exists
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(err)
			return
		}
		if existing_referrer.UniqueID == referrer_id {
			repo.Update(referrer_id, device_id) // Update the device_id in the referrer table
		}
		repo.Save(&record)
		response.WriteHeader(http.StatusOK)                                      // Send response
		json.NewEncoder(response).Encode("Referral Code generated successfully") // Send response
		json.NewEncoder(response).Encode(record)                                 // Send response
	}
}

// Get device by device_id "/referrals/{device_id}"
func GetDevice(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	device_id := request.FormValue("device_id")
	if device_id == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode("Device id is required")
		return
	}
	record, err := repo.FindDevice(device_id)
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

// Get all devices "/referrals"
func GetAllDevices(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	records, err := repo.FindAll()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(err)
		return
	}
	response.WriteHeader(http.StatusOK) // Send response
	json.NewEncoder(response).Encode(records)
}

// Get Referred Counts "/referrals/{device_id}/counts"
func GetReferredCounts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "application/json")
	device_id := request.FormValue("device_id")
	if device_id == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode("Device id is required")
		return
	}
	record, err := repo.FindDevice(device_id)
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
	json.NewEncoder(response).Encode("Referred Counts")
	json.NewEncoder(response).Encode(len(record.ReferredIDS))
}
