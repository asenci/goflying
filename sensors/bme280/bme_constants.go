package bme280

import (
	"encoding/binary"
	"sync"
)

const (
	ChipID    = 0x60
	ResetCode = 0xB6

	CalibrationDataSize    = 26
	CalibrationHumDataSize = 7

	MeasurementDataSize = 8
	MeasurementHumMax   = 102400
	MeasurementHumMin   = 0
	MeasurementPressMax = 11000000
	MeasurementPressMin = 3000000
	MeasurementTempMax  = 8000
	MeasurementTempMin  = -4000

	// Register addresses
	RegisterCalibrationData    = 0x88
	RegisterChipID             = 0xD0
	RegisterSoftReset          = 0xE0
	RegisterHumCalibrationData = 0xE1
	RegisterHumControl         = 0xF2
	RegisterStatus             = 0xF3
	RegisterControl            = 0xF4
	RegisterConfig             = 0xF5
	RegisterPressDataMSB       = 0xF7
	RegisterPressDataLSB       = 0xF8
	RegisterPressDataXLSB      = 0xF9
	RegisterTempDataMSB        = 0xFA
	RegisterTempDataLSB        = 0xFB
	RegisterTempDataXLSB       = 0xFC
	RegisterHumDataMSB         = 0xFD
	RegisterHumDataLSB         = 0xFE
)

type CalibrationData struct {
	rawCalibrationData         []byte
	rawHumidityCalibrationData []byte

	tFine int32
}

// Compensate humidity measurements using calibration data
// Black magic imported from the datasheet
func (cal CalibrationData) CompensateHumidity(rawHum int32) uint32 {
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

// Compensate pressure measurements using calibration data
// Black magic imported from the datasheet
func (cal CalibrationData) CompensatePressure(rawPress int32) uint32 {
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

// Compensate temperature measurements using calibration data
// Black magic imported from the datasheet
func (cal CalibrationData) CompensateTemperature(rawTemp int32) int32 {
	var1 := (rawTemp >> 3) - (int32(cal.digT1()) << 1)
	var1 = (var1 * int32(cal.digT2())) >> 11

	var2 := (rawTemp >> 4) - int32(cal.digT1())
	var2 = (((var2 * var2) >> 12) * int32(cal.digH3())) >> 14

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

func (cal CalibrationData) digH1() byte {
	return cal.rawCalibrationData[25]
}

func (cal CalibrationData) digH2() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawHumidityCalibrationData[0:2]))
}

func (cal CalibrationData) digH3() byte {
	return cal.rawHumidityCalibrationData[2]
}

func (cal CalibrationData) digH4() int16 {
	return int16((cal.rawHumidityCalibrationData[3] << 4) | (cal.rawHumidityCalibrationData[4] & 0x0F))
}

func (cal CalibrationData) digH5() int16 {
	return int16((cal.rawHumidityCalibrationData[5] << 4) | (cal.rawHumidityCalibrationData[4] >> 4))
}

func (cal CalibrationData) digH6() byte {
	return cal.rawHumidityCalibrationData[6]
}

func (cal CalibrationData) digP1() uint16 {
	return binary.LittleEndian.Uint16(cal.rawCalibrationData[6:8])
}

func (cal CalibrationData) digP2() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[8:10]))
}

func (cal CalibrationData) digP3() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[10:12]))
}

func (cal CalibrationData) digP4() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[12:14]))
}

func (cal CalibrationData) digP5() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[14:16]))
}

func (cal CalibrationData) digP6() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[16:18]))
}

func (cal CalibrationData) digP7() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[18:20]))
}

func (cal CalibrationData) digP8() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[20:22]))
}

func (cal CalibrationData) digP9() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[22:24]))
}

func (cal CalibrationData) digT1() uint16 {
	return binary.LittleEndian.Uint16(cal.rawCalibrationData[0:2])
}

func (cal CalibrationData) digT2() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[2:4]))
}

func (cal CalibrationData) digT3() int16 {
	return int16(binary.LittleEndian.Uint16(cal.rawCalibrationData[4:6]))
}

func NewCalibrationData(calibrationData, humidityCalibrationData []byte) *CalibrationData {
	return &CalibrationData{
		rawCalibrationData:         calibrationData,
		rawHumidityCalibrationData: humidityCalibrationData,
	}
}

type ConfigByte byte

func (b ConfigByte) FilterCoefficient() FilterCoefficient {
	return FilterCoefficient(getValue(byte(b), FilterCoefficientOversamplingSize, FilterCoefficientOversamplingShift))
}

func (b ConfigByte) SetFilterCoefficient(value FilterCoefficient) ConfigByte {
	return ConfigByte(setValue(byte(b), FilterCoefficientOversamplingSize, FilterCoefficientOversamplingShift, byte(value)))
}

func (b ConfigByte) InactiveDuration() InactiveDuration {
	return InactiveDuration(getValue(byte(b), InactiveDurationSize, InactiveDurationShift))
}

func (b ConfigByte) SetInactiveDuration(value InactiveDuration) ConfigByte {
	return ConfigByte(setValue(byte(b), InactiveDurationSize, InactiveDurationShift, byte(value)))
}

type ControlByte byte

func (b ControlByte) PressureOversampling() PressureOversampling {
	return PressureOversampling(getValue(byte(b), PressureOversamplingSize, PressureOversamplingShift))
}

func (b ControlByte) SetPressureOversampling(value PressureOversampling) ControlByte {
	return ControlByte(setValue(byte(b), PressureOversamplingSize, PressureOversamplingShift, byte(value)))
}

func (b ControlByte) RunMode() RunMode {
	return RunMode(getValue(byte(b), RunModeSize, RunModeShift))
}

func (b ControlByte) SetRunMode(value RunMode) ControlByte {
	return ControlByte(setValue(byte(b), RunModeSize, RunModeShift, byte(value)))
}

func (b ControlByte) TemperatureOversampling() TemperatureOversampling {
	return TemperatureOversampling(getValue(byte(b), TemperatureOversamplingSize, TemperatureOversamplingShift))
}

func (b ControlByte) SetTemperatureOversampling(value TemperatureOversampling) ControlByte {
	return ControlByte(setValue(byte(b), TemperatureOversamplingSize, TemperatureOversamplingShift, byte(value)))
}

type FilterCoefficient byte

func (b FilterCoefficient) Value() int {
	if b > 0 {
		return 1 << b
	}

	return 0
}

const (
	FilterCoefficientOversamplingShift = 2
	FilterCoefficientOversamplingSize  = 3

	FilterCoefficientOff FilterCoefficient = 0x00
	FilterCoefficient2   FilterCoefficient = 0x01
	FilterCoefficient4   FilterCoefficient = 0x02
	FilterCoefficient8   FilterCoefficient = 0x03
	FilterCoefficient16  FilterCoefficient = 0x04
)

type HumidityControlByte byte

func (b HumidityControlByte) HumidityOversampling() HumidityOversampling {
	return HumidityOversampling(getValue(byte(b), HumidityOversamplingSize, HumidityOversamplingShift))
}

func (b HumidityControlByte) SetHumidityOversampling(value HumidityOversampling) HumidityControlByte {
	return HumidityControlByte(setValue(byte(b), HumidityOversamplingSize, HumidityOversamplingShift, byte(value)))
}

type HumidityOversampling byte

func (b HumidityOversampling) Value() int {
	return oversamplingValue(byte(b))
}

const (
	HumidityOversamplingShift = 0
	HumidityOversamplingSize  = 3

	HumidityOversamplingSkipped HumidityOversampling = 0x00
	HumidityOversampling1x      HumidityOversampling = 0x01
	HumidityOversampling2x      HumidityOversampling = 0x02
	HumidityOversampling4x      HumidityOversampling = 0x03
	HumidityOversampling8x      HumidityOversampling = 0x04
	HumidityOversampling16x     HumidityOversampling = 0x05
)

type I2CAddress byte

const (
	I2CAddressLow  I2CAddress = 0x76
	I2CAddressHigh I2CAddress = 0x77
)

type InactiveDuration byte

func (d InactiveDuration) Milliseconds() float64 {
	switch d {
	case InactiveDuration0_5ms:
		return 0.5
	case InactiveDuration62_5ms:
		return 62.5
	case InactiveDuration125ms:
		return 125
	case InactiveDuration250ms:
		return 250
	case InactiveDuration500ms:
		return 500
	case InactiveDuration1000ms:
		return 1000
	case InactiveDuration10ms:
		return 10
	case InactiveDuration20ms:
		return 20
	default:
		return 0
	}
}

const (
	InactiveDurationShift = 5
	InactiveDurationSize  = 3

	InactiveDuration0_5ms  InactiveDuration = 0x00
	InactiveDuration62_5ms InactiveDuration = 0x01
	InactiveDuration125ms  InactiveDuration = 0x02
	InactiveDuration250ms  InactiveDuration = 0x03
	InactiveDuration500ms  InactiveDuration = 0x04
	InactiveDuration1000ms InactiveDuration = 0x05
	InactiveDuration10ms   InactiveDuration = 0x06
	InactiveDuration20ms   InactiveDuration = 0x07
)

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

type PressureOversampling byte

func (b PressureOversampling) Value() int {
	return oversamplingValue(byte(b))
}

const (
	PressureOversamplingShift = 2
	PressureOversamplingSize  = 3

	PressureOversamplingSkipped PressureOversampling = 0x00
	PressureOversampling1x      PressureOversampling = 0x01
	PressureOversampling2x      PressureOversampling = 0x02
	PressureOversampling4x      PressureOversampling = 0x03
	PressureOversampling8x      PressureOversampling = 0x04
	PressureOversampling16x     PressureOversampling = 0x05
)

type RunMode byte

const (
	RunModeShift = 0
	RunModeSize  = 2

	RunModeSleep  RunMode = 0x00
	RunModeForced RunMode = 0x01
	RunModeNormal RunMode = 0x03
)

type SettingFunc func(configByte ConfigByte, controlByte ControlByte, humControlByte HumidityControlByte) (ConfigByte, ControlByte, HumidityControlByte, error)

type TemperatureOversampling byte

func (b TemperatureOversampling) Value() int {
	return oversamplingValue(byte(b))
}

const (
	TemperatureOversamplingShift = 5
	TemperatureOversamplingSize  = 3

	TemperatureOversamplingSkipped TemperatureOversampling = 0x00
	TemperatureOversampling1x      TemperatureOversampling = 0x01
	TemperatureOversampling2x      TemperatureOversampling = 0x02
	TemperatureOversampling4x      TemperatureOversampling = 0x03
	TemperatureOversampling8x      TemperatureOversampling = 0x04
	TemperatureOversampling16x     TemperatureOversampling = 0x05
)

func getValue(b byte, size, shift int) byte {
	valueMask := byte((1 << size) - 1)
	getMask := valueMask << shift

	return (b & getMask) >> shift
}

func setValue(b byte, size, shift int, value byte) byte {
	valueMask := byte((1 << size) - 1)
	getMask := valueMask << shift
	clearMask := ^getMask

	return (b & clearMask) | ((value & valueMask) << shift)
}

func oversamplingValue(b byte) int {
	if b > 0 {
		return 1 << (b - 1)
	}

	return 0
}
