package bme280

import "sync"

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
}

func (d *MeasurementData) Humidity() float64 {
	humidity := d.calibrationData.CompensateHumidity(int32(d.rawHumidity()))

	// Humidity is given in 1024 * % relative humidity
	return float64(humidity) / 1024
}

func (d *MeasurementData) Pressure() float64 {
	pressure := d.calibrationData.CompensatePressure(int32(d.rawPressure()))

	// Pressure is given in 100 * Pascal
	return float64(pressure) / 100 / 100
}

func (d *MeasurementData) rawHumidity() uint32 {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return (uint32(d.rawHumidityDataMSB) << 8) | uint32(d.rawHumidityDataLSB)
}

func (d *MeasurementData) rawPressure() uint32 {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return (uint32(d.rawPressureDataMSB) << 12) | (uint32(d.rawPressureDataLSB) << 4) | (uint32(d.rawPressureDataXLSB) >> 4)
}

func (d *MeasurementData) rawTemperature() uint32 {
	d.mux.RLock()
	defer d.mux.RUnlock()

	return (uint32(d.rawTemperatureDataMSB) << 12) | (uint32(d.rawTemperatureDataLSB) << 4) | (uint32(d.rawTemperatureDataXLSB) >> 4)
}

func (d *MeasurementData) SetCalibrationData(cal *CalibrationData) {
	d.calibrationData = cal
}

func (d *MeasurementData) Temperature() float64 {
	temperature := d.calibrationData.CompensateTemperature(int32(d.rawTemperature()))

	// Pressure is given in 100 * ºCelsius
	return float64(temperature) / 100
}

func (d *MeasurementData) Update(data []byte) {
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
}
