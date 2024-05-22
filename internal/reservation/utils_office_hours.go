package reservation

import "time"

func (o *OfficeHours) IsValid() bool {
	// Check if Open and Close are in the format xx:xx
	openTime, err := time.Parse("15:04", o.Open)
	if err != nil {
		return false
	}

	closeTime, err := time.Parse("15:04", o.Close)
	if err != nil {
		return false
	}

	// Check if Open time is before Close time
	if openTime.After(closeTime) {
		return false
	}

	return true
}