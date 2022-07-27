package bme280

import (
	"testing"
)

func TestCalibrationData_CompensateHumidity(t *testing.T) {
	cal := NewCalibrationData(
		[]byte{0xEE, 0x6E, 0xCA, 0x68, 0x32, 0x00, 0x48, 0x92, 0xD6, 0xD6, 0xD0, 0x0B, 0x7F, 0x19, 0x1F, 0x00, 0xF9, 0xFF, 0xAC, 0x26, 0x0A, 0xD8, 0xBD, 0x10, 0x00, 0x4B},
		[]byte{0x80, 0x01, 0x00, 0x10, 0x2D, 0x03, 0x1E},
	)

	tests := []struct {
		name    string
		cal     *CalibrationData
		rawTemp int32
		rawHum  int32
		want    uint32
	}{
		{"initial data", cal, 524288, 32768, 92947},
		{"case 01", cal, 526634, 32768, 93031},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// calculate tFine
			tt.cal.CompensateTemperature(tt.rawTemp)

			if got := tt.cal.CompensateHumidity(tt.rawHum); got != tt.want {
				t.Errorf("CompensateHumidity() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalibrationData_CompensatePressure(t *testing.T) {
	cal := NewCalibrationData(
		[]byte{0xEE, 0x6E, 0xCA, 0x68, 0x32, 0x00, 0x48, 0x92, 0xD6, 0xD6, 0xD0, 0x0B, 0x7F, 0x19, 0x1F, 0x00, 0xF9, 0xFF, 0xAC, 0x26, 0x0A, 0xD8, 0xBD, 0x10, 0x00, 0x4B},
		[]byte{0x80, 0x01, 0x00, 0x10, 0x2D, 0x03, 0x1E},
	)

	tests := []struct {
		name     string
		cal      *CalibrationData
		rawTemp  int32
		rawPress int32
		want     uint32
	}{
		{"initial data", cal, 524288, 524288, 6965154},
		{"case 01", cal, 526634, 338837, 10064647},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// calculate tFine
			tt.cal.CompensateTemperature(tt.rawTemp)

			if got := tt.cal.CompensatePressure(tt.rawPress); got != tt.want {
				t.Errorf("CompensatePressure() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalibrationData_CompensateTemperature(t *testing.T) {
	cal := NewCalibrationData(
		[]byte{0xEE, 0x6E, 0xCA, 0x68, 0x32, 0x00, 0x48, 0x92, 0xD6, 0xD6, 0xD0, 0x0B, 0x7F, 0x19, 0x1F, 0x00, 0xF9, 0xFF, 0xAC, 0x26, 0x0A, 0xD8, 0xBD, 0x10, 0x00, 0x4B},
		[]byte{0x80, 0x01, 0x00, 0x10, 0x2D, 0x03, 0x1E},
	)

	tests := []struct {
		name    string
		cal     *CalibrationData
		rawTemp int32
		want    int32
	}{
		{"initial data", cal, 524288, 2236},
		{"case 01", cal, 526634, 2311},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cal.CompensateTemperature(tt.rawTemp); got != tt.want {
				t.Errorf("CompensateTemperature() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNewCalibrationData(t *testing.T) {
	type calParams struct {
		digH1 byte
		digH2 int16
		digH3 byte
		digH4 int16
		digH5 int16
		digH6 byte
		digP1 uint16
		digP2 int16
		digP3 int16
		digP4 int16
		digP5 int16
		digP6 int16
		digP7 int16
		digP8 int16
		digP9 int16
		digT1 uint16
		digT2 int16
		digT3 int16
	}

	tests := []struct {
		name                    string
		calibrationData         []byte
		humidityCalibrationData []byte
		want                    calParams
	}{
		{
			"case 01",
			[]byte{0xEE, 0x6E, 0xCA, 0x68, 0x32, 0x00, 0x48, 0x92, 0xD6, 0xD6, 0xD0, 0x0B, 0x7F, 0x19, 0x1F, 0x00, 0xF9, 0xFF, 0xAC, 0x26, 0x0A, 0xD8, 0xBD, 0x10, 0x00, 0x4B},
			[]byte{0x80, 0x01, 0x00, 0x10, 0x2D, 0x03, 0x1E},
			calParams{digH1: 0x4B, digH2: 384, digH3: 0x00, digH4: 269, digH5: 50, digH6: 0x1E, digP1: 37448, digP2: -10538, digP3: 3024, digP4: 6527, digP5: 31, digP6: -7, digP7: 9900, digP8: -10230, digP9: 4285, digT1: 28398, digT2: 26826, digT3: 50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := NewCalibrationData(tt.calibrationData, tt.humidityCalibrationData)

			if got := cal.digH1(); got != tt.want.digH1 {
				t.Errorf("digH1() = %02X, want %02X", got, tt.want.digH1)
			}

			if got := cal.digH2(); got != tt.want.digH2 {
				t.Errorf("digH2() = %d, want %d", got, tt.want.digH2)
			}

			if got := cal.digH3(); got != tt.want.digH3 {
				t.Errorf("digH3() = %02X, want %02X", got, tt.want.digH3)
			}

			if got := cal.digH4(); got != tt.want.digH4 {
				t.Errorf("digH4() = %d, want %d", got, tt.want.digH4)
			}

			if got := cal.digH5(); got != tt.want.digH5 {
				t.Errorf("digH5() = %d, want %d", got, tt.want.digH5)
			}

			if got := cal.digH6(); got != tt.want.digH6 {
				t.Errorf("digH6() = %02X, want %02X", got, tt.want.digH6)
			}

			if got := cal.digP1(); got != tt.want.digP1 {
				t.Errorf("digP1() = %d, want %d", got, tt.want.digP1)
			}

			if got := cal.digP2(); got != tt.want.digP2 {
				t.Errorf("digP2() = %d, want %d", got, tt.want.digP2)
			}

			if got := cal.digP3(); got != tt.want.digP3 {
				t.Errorf("digP3() = %d, want %d", got, tt.want.digP3)
			}

			if got := cal.digP4(); got != tt.want.digP4 {
				t.Errorf("digP4() = %d, want %d", got, tt.want.digP4)
			}

			if got := cal.digP5(); got != tt.want.digP5 {
				t.Errorf("digP5() = %d, want %d", got, tt.want.digP5)
			}

			if got := cal.digP6(); got != tt.want.digP6 {
				t.Errorf("digP6() = %d, want %d", got, tt.want.digP6)
			}

			if got := cal.digP7(); got != tt.want.digP7 {
				t.Errorf("digP7() = %d, want %d", got, tt.want.digP7)
			}

			if got := cal.digP8(); got != tt.want.digP8 {
				t.Errorf("digP8() = %d, want %d", got, tt.want.digP8)
			}

			if got := cal.digP9(); got != tt.want.digP9 {
				t.Errorf("digP9() = %d, want %d", got, tt.want.digP9)
			}

			if got := cal.digT1(); got != tt.want.digT1 {
				t.Errorf("digT1() = %d, want %d", got, tt.want.digT1)
			}

			if got := cal.digT2(); got != tt.want.digT2 {
				t.Errorf("digT2() = %d, want %d", got, tt.want.digT2)
			}

			if got := cal.digT3(); got != tt.want.digT3 {
				t.Errorf("digT3() = %d, want %d", got, tt.want.digT3)
			}
		})
	}
}
