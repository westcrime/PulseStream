package models

import "time"

type TelemetryEvent struct {
	DeviceID  string    `json:"device_id"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
