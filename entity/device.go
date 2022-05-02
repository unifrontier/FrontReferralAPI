package entity

// Data structure
type Device struct {
	DeviceID    string   `json:"device_id"`    // Serial number
	UniqueID    string   `json:"unique_id"`    // ReferralCode
	ReferrerID  string   `json:"referrer_id"`  // Referred ID
	ReferredIDS []string `json:"referrer_ids"` // Referred IDs
}
