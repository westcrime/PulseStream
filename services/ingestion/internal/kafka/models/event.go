package models

type Event struct {
	DeviceID string  `json:"device_id"`
	Type     string  `json:"type"`
	Value    float64 `json:"value"`
}
