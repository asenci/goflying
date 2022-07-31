package bme280

import (
	"fmt"
	"sync"
	"time"

	"github.com/westphae/goflying"
)

// MeasurementData represents the measurement values read from the chip
// measurement registers. It uses the factory calibration data to compensate the
// raw measurement data and return formatted values. Temperature must be read
// first so the t_fine compensation value can be calculated for Humidity and
// Pressure compensation.
type MeasurementData struct {
	mux sync.RWMutex

	calibrationData *CalibrationData

	rawHumidityDataMSB byte
	rawHumidityDataLSB byte

	rawPressureDataMSB  byte
	rawPressureDataLSB  byte
	rawPressureDataXLSB byte

	rawTemperatureDataMSB  byte
	rawTemperatureDataLSB  byte
	rawTemperatureDataXLSB byte

	timestamp time.Time
}

// Humidity returns the relative humidity in % as a float
func (d *MeasurementData) Humidity() goflying.RelativeHumidity {
	return formatHumidityRH(d.calibrationData.CompensateHumidity(d.rawHumidity()))
}

// Pressure returns the air pressure in hPa (same as millibars) as a float
func (d *MeasurementData) Pressure() goflying.HPa {
	return formatPressureHPa(d.calibrationData.CompensatePressure(d.rawPressure()))
}

func (d *MeasurementData) rawHumidity() int32 {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return (int32(d.rawHumidityDataMSB) << 8) | int32(d.rawHumidityDataLSB)
}

func (d *MeasurementData) rawPressure() int32 {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return (int32(d.rawPressureDataMSB) << 12) | (int32(d.rawPressureDataLSB) << 4) | (int32(d.rawPressureDataXLSB) >> 4)
}

func (d *MeasurementData) rawTemperature() int32 {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return (int32(d.rawTemperatureDataMSB) << 12) | (int32(d.rawTemperatureDataLSB) << 4) | (int32(d.rawTemperatureDataXLSB) >> 4)
}

// SetCalibrationData stores the factory calibration data for measurement
// compensation
func (d *MeasurementData) SetCalibrationData(cal *CalibrationData) {
	d.calibrationData = cal
}

// Temperature returns the ambient temperature in ºC as a float
func (d *MeasurementData) Temperature() goflying.Celsius {
	return formatTemperatureC(d.calibrationData.CompensateTemperature(d.rawTemperature()))
}

// Update sets the raw measurement data as read from the chip measurement registers
func (d *MeasurementData) Update(data []byte, timestamp time.Time) {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.rawPressureDataMSB = data[0]
	d.rawPressureDataLSB = data[1]
	d.rawPressureDataXLSB = data[2]

	d.rawTemperatureDataMSB = data[3]
	d.rawTemperatureDataLSB = data[4]
	d.rawTemperatureDataXLSB = data[5]

	d.rawHumidityDataMSB = data[6]
	d.rawHumidityDataLSB = data[7]

	d.timestamp = timestamp
}

func (d *MeasurementData) String() string {
	return fmt.Sprintf(
		"temperature: %.2fºC (raw data: %d), humidity: %.2f%% (raw data: %d), pressure: %.2fhPa (raw data: %d)",
		d.Temperature(), d.rawTemperature(),
		d.Humidity(), d.rawHumidity(),
		d.Pressure(), d.rawPressure(),
	)
}

func (d *MeasurementData) Timestamp() time.Time {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return d.timestamp
}

func NewMeasurementData(cal *CalibrationData) *MeasurementData {
	return &MeasurementData{calibrationData: cal}
}

// Humidity is given in 1024 * % relative humidity
func formatHumidityRH(value uint32) goflying.RelativeHumidity {
	return goflying.RelativeHumidity(float64(value) / 1024)
}

// Pressure is given in 100 * Pascal
func formatPressureHPa(value uint32) goflying.HPa {
	return goflying.HPa(float64(value) / 100 / 100)
}

// Temperature is given in 100 * ºCelsius
func formatTemperatureC(value int32) goflying.Celsius {
	return goflying.Celsius(float64(value) / 100)
}
