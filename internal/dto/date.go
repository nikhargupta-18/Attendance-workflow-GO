package dto

import (
	"encoding/json"
	"time"
)

// Date is a custom type that accepts YYYY-MM-DD format in JSON
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	// Try YYYY-MM-DD format first
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		d.Time = t
		return nil
	}

	// Fallback to RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		d.Time = t
		return nil
	}

	// Try RFC3339 without timezone
	if t, err := time.Parse("2006-01-02T15:04:05Z", dateStr); err == nil {
		d.Time = t
		return nil
	}

	return &time.ParseError{
		Layout:  "2006-01-02",
		Value:   dateStr,
		Message: "date must be in YYYY-MM-DD format",
	}
}

// MarshalJSON implements json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

