package bme280

import (
	"testing"
)

func TestMeasurementData_Calibrated(t *testing.T) {
	cal := NewCalibrationData(
		[]byte{0xEE, 0x6E, 0xCA, 0x68, 0x32, 0x00, 0x48, 0x92, 0xD6, 0xD6, 0xD0, 0x0B, 0x7F, 0x19, 0x1F, 0x00, 0xF9, 0xFF, 0xAC, 0x26, 0x0A, 0xD8, 0xBD, 0x10, 0x00, 0x4B},
		[]byte{0x80, 0x01, 0x00, 0x10, 0x2D, 0x03, 0x1E},
	)

	tests := []struct {
		name        string
		data        []byte
		cal         *CalibrationData
		Humidity    float64
		Pressure    float64
		Temperature float64
	}{
		{"initial data", []byte{0x80, 0x00, 0x00, 0x80, 0x00, 0x00, 0x80, 0x00}, cal, 90.7587890625, 696.5183999999999, 22.36},
		{"case 01", []byte{0x52, 0xB9, 0x50, 0x80, 0x92, 0xA0, 0x68, 0x4D}, cal, 55.40234375, 1006.4693, 23.11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MeasurementData{}
			d.SetCalibrationData(tt.cal)
			d.Update(tt.data)

			// Get temperature first to update tFine
			if got := d.Temperature(); (got) != (tt.Temperature) {
				t.Errorf("Temperature() = %g, want %g", got, tt.Temperature)
			}

			if got := d.Humidity(); (got) != (tt.Humidity) {
				t.Errorf("Humidity() = %g, want %g", got, tt.Humidity)
			}

			if got := d.Pressure(); (got) != (tt.Pressure) {
				t.Errorf("Pressure() = %g, want %g", got, tt.Pressure)
			}

		})
	}
}

func TestMeasurementData_Raw(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		rawHumidity    uint32
		rawPressure    uint32
		rawTemperature uint32
	}{
		{"initial data", []byte{0x80, 0x00, 0x00, 0x80, 0x00, 0x00, 0x80, 0x00}, 32768, 524288, 524288},
		{"case 01", []byte{0x52, 0xB9, 0x50, 0x80, 0x92, 0xA0, 0x68, 0x4D}, 26701, 338837, 526634},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MeasurementData{}
			d.Update(tt.data)

			if got := d.rawHumidity(); got != tt.rawHumidity {
				t.Errorf("rawHumidity() = %d, want %d", got, tt.rawHumidity)
			}

			if got := d.rawPressure(); got != tt.rawPressure {
				t.Errorf("rawPressure() = %d, want %d", got, tt.rawPressure)
			}

			if got := d.rawTemperature(); got != tt.rawTemperature {
				t.Errorf("rawTemperature() = %d, want %d", got, tt.rawTemperature)
			}
		})
	}
}
