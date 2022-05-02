package helper

import (
	"FrontReferralAPI/entity"
	"fmt"
)

// Record exists or not
func IsExists(device_id string, records []entity.Device) (result bool) {
	fmt.Println(records, device_id)
	result = false
	for _, record := range records {
		if record.DeviceID == device_id {
			result = true
			break
		}
	}
	return result
}
