package bme280

import (
	"encoding/binary"
	"fmt"
	"sync"
)

// CalibrationData represents the factory calibration bits on the BME280 chip.
// It is used for performing value compensation on the raw measurement values.
// CompensateTemperature must be run before CompensateHumidity and
// CompensatePressure so the t_fine value is updated.
type CalibrationData struct {
	rawCalibrationData         []byte
	rawHumidityCalibrationData []byte

	tFine    int32
	tFineMux sync.RWMutex
}

// CompensateHumidity compensates the humidity measurements using the factory
// calibration data. Black-box algorithm imported from the datasheet. Humidity
// is given in 1024 * % relative humidity. i.e.: 79974.4 represents 78.1%
// relative humidity.
func (cal *CalibrationData) CompensateHumidity(rawHum int32) uint32 {
	cal.tFineMux.RLock()
	defer cal.tFineMux.RUnlock()

	var1 := cal.tFine - 76800

	var2 := rawHum << 14

	var3 := int32(cal.digH4()) << 20

	var4 := int32(cal.digH5()) * var1

	var5 := (var2 - var3 - var4 + 16384) >> 15

	var2 = (var1 * int32(cal.digH6())) >> 10

	var3 = (var1 * int32(cal.digH3())) >> 11

	var4 = ((var2 * (var3 + 32768)) >> 10) + 2097152

	var2 = ((var4 * int32(cal.digH2())) + 8192) >> 14

	var3 = var5 * var2

	var4 = ((var3 >> 15) * (var3 >> 15)) >> 7

	var5 = var3 - ((var4 * int32(cal.digH1())) >> 4)

	if var5 < 0 {
		return MeasurementHumMin
	}

	humidity := var5 >> 12

	if humidity > MeasurementHumMax {
		return MeasurementHumMax
	}

	return uint32(humidity)
}

// CompensatePressure compensates the pressure measurements using the factory
// calibration data. Black-box algorithm imported from the datasheet. Pressure
// is given in 100 * Pascal. i.e.: 10132500.00 represents 1013.25hPa.
func (cal *CalibrationData) CompensatePressure(rawPress int32) uint32 {
	cal.tFineMux.RLock()
	defer cal.tFineMux.RUnlock()

	var1 := int64(cal.tFine) - 128000

	var2 := var1 * var1 * int64(cal.digP6())
	var2 += (var1 * int64(cal.digP5())) << 17
	var2 += int64(cal.digP4()) << 35

	var1 = ((var1 * var1 * int64(cal.digP3())) >> 8) + ((var1 * int64(cal.digP2())) << 12)

	var3 := int64(1) << 47

	var1 = ((var3 + var1) * int64(cal.digP1())) >> 33

	// Avoid division by zero
	if var1 == 0 {
		return MeasurementPressMin
	}

	var4 := 1048576 - int64(rawPress)
	var4 = (((var4 << 31) - var2) * 3125) / var1

	var1 = (int64(cal.digP9()) * (var4 >> 13) * (var4 >> 13)) >> 25

	var2 = (int64(cal.digP8()) * var4) >> 19

	var4 = ((var4 + var1 + var2) >> 8) + (int64(cal.digP7()) << 4)

	pressure := ((var4 >> 1) * 100) >> 7

	if pressure < MeasurementPressMin {
		return MeasurementPressMin
	}

	if pressure > MeasurementPressMax {
		return MeasurementPressMax
	}

	return uint32(pressure)
}

// CompensateTemperature compensates the temperature measurements using the
// factory calibration data. Black-box algorithm imported from the datasheet.
// Temperature is given in 100 * ºCelsius. i.e.: 1500.00 represents 15.0ºC.
func (cal *CalibrationData) CompensateTemperature(rawTemp int32) int32 {
	cal.tFineMux.Lock()
	defer cal.tFineMux.Unlock()

	var1 := (rawTemp >> 3) - (int32(cal.digT1()) << 1)
	var1 = (var1 * int32(cal.digT2())) >> 11

	var2 := (rawTemp >> 4) - int32(cal.digT1())
	var2 = (((var2 * var2) >> 12) * int32(cal.digT3())) >> 14

	cal.tFine = var1 + var2

	temperature := ((cal.tFine * 5) + 128) >> 8

	if temperature < MeasurementTempMin {
		return MeasurementTempMin
	}

	if temperature > MeasurementTempMax {
		return MeasurementTempMax
	}

	return temperature
}

func (cal *CalibrationData) digH1() byte {
	return cal.rawCalibrationData[25]
}

func (cal *CalibrationData) digH2() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawHumidityCalibrationData[0:2]))
}

func (cal *CalibrationData) digH3() byte {
	return cal.rawHumidityCalibrationData[2]
}

func (cal *CalibrationData) digH4() int16 {
	return (int16(cal.rawHumidityCalibrationData[3]) << 4) | (int16(cal.rawHumidityCalibrationData[4]) & 0x0F)
}

func (cal *CalibrationData) digH5() int16 {
	return (int16(cal.rawHumidityCalibrationData[5]) << 4) | (int16(cal.rawHumidityCalibrationData[4]) >> 4)
}

func (cal *CalibrationData) digH6() byte {
	return cal.rawHumidityCalibrationData[6]
}

func (cal *CalibrationData) digP1() uint16 {
	return binary.LittleEndian.Uint16(cal.rawCalibrationData[6:8])
}

func (cal *CalibrationData) digP2() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[8:10]))
}

func (cal *CalibrationData) digP3() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[10:12]))
}

func (cal *CalibrationData) digP4() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[12:14]))
}

func (cal *CalibrationData) digP5() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[14:16]))
}

func (cal *CalibrationData) digP6() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[16:18]))
}

func (cal *CalibrationData) digP7() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[18:20]))
}

func (cal *CalibrationData) digP8() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[20:22]))
}

func (cal *CalibrationData) digP9() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[22:24]))
}

func (cal *CalibrationData) digT1() uint16 {
	return binary.LittleEndian.Uint16(cal.rawCalibrationData[0:2])
}

func (cal *CalibrationData) digT2() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[2:4]))
}

func (cal *CalibrationData) digT3() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[4:6]))
}

func (cal *CalibrationData) String() string {
	return fmt.Sprintf(
		"dig_H1: 0x%02X, dig_H2: %d, dig_H3: 0x%02X, dig_H4: %d, dig_H5: %d, dig_H6: 0x%02X, dig_P1: %d, dig_P2: %d, dig_P3: %d, dig_P4: %d, dig_P5: %d, dig_P6: %d, dig_P7: %d, dig_P8: %d, dig_P9: %d, dig_T1: %d, dig_T2: %d, dig_T3: %d",
		cal.digH1(), cal.digH2(), cal.digH3(), cal.digH4(), cal.digH5(), cal.digH6(),
		cal.digP1(), cal.digP2(), cal.digP3(), cal.digP4(), cal.digP5(), cal.digP6(), cal.digP7(), cal.digP8(), cal.digP9(),
		cal.digT1(), cal.digT2(), cal.digT3(),
	)
}

func NewCalibrationData(calibrationData, humidityCalibrationData []byte) *CalibrationData {
	return &CalibrationData{
		rawCalibrationData:         calibrationData,
		rawHumidityCalibrationData: humidityCalibrationData,
	}
}
