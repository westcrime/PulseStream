package event

import (
	"errors"
	"fmt"
	"time"
)

type Type string

const (
	Temperature  Type = "temperature"
	Humidity     Type = "humidity"
	Brightness   Type = "brightness"
	Pressure     Type = "pressure"
	CO2          Type = "co2"
	AirQuality   Type = "air_quality"
	Noise        Type = "noise"
	UVIndex      Type = "uv_index"
	Rain         Type = "rain"
	WindSpeed    Type = "wind_speed"
	SoilMoisture Type = "soil_moisture"
)

type Spec struct {
	Unit string
	Min  float64
	Max  float64
}

var specs = map[Type]Spec{
	Temperature:  {Unit: "°C", Min: -90, Max: 60},
	Humidity:     {Unit: "%", Min: 0, Max: 100},
	Brightness:   {Unit: "lux", Min: 0, Max: 200_000},
	Pressure:     {Unit: "hPa", Min: 300, Max: 1100},
	CO2:          {Unit: "ppm", Min: 0, Max: 10_000},
	AirQuality:   {Unit: "µg/m³", Min: 0, Max: 1000},
	Noise:        {Unit: "dB", Min: 0, Max: 200},
	UVIndex:      {Unit: "", Min: 0, Max: 20},
	Rain:         {Unit: "mm/h", Min: 0, Max: 500},
	WindSpeed:    {Unit: "m/s", Min: 0, Max: 120},
	SoilMoisture: {Unit: "%", Min: 0, Max: 100},
}

type Event struct {
	DeviceID  string    `json:"device_id"`
	Type      Type      `json:"type"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (t Type) Valid() bool {
	_, ok := specs[t]
	return ok
}

func (t Type) Spec() (Spec, bool) {
	s, ok := specs[t]
	return s, ok
}

func (e Event) Validate() error {
	if e.DeviceID == "" {
		return errors.New("device_id is empty")
	}
	spec, ok := e.Type.Spec()
	if !ok {
		return fmt.Errorf("unknown event type %q", e.Type)
	}
	if e.Value < spec.Min || e.Value > spec.Max {
		return fmt.Errorf("%s: value %v out of range [%v, %v] %s",
			e.Type, e.Value, spec.Min, spec.Max, spec.Unit)
	}
	if e.Timestamp.IsZero() {
		return errors.New("timestamp is empty")
	}
	if e.Timestamp.After(time.Now().Add(time.Minute)) {
		return errors.New("timestamp is in the future")
	}
	return nil
}
