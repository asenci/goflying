package bme280

import (
	"fmt"
	"strconv"
)

//goland:noinspection GoCommentStart
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
	RegisterCtrlHum            = 0xF2
	RegisterStatus             = 0xF3
	RegisterCtrlMeas           = 0xF4
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

type Config byte

func (b Config) FilterCoefficient() FilterCoefficient {
	return FilterCoefficient(getValue(byte(b), FilterCoefficientOversamplingSize, FilterCoefficientOversamplingShift))
}

func (b Config) SetFilterCoefficient(value FilterCoefficient) Config {
	return Config(setValue(byte(b), FilterCoefficientOversamplingSize, FilterCoefficientOversamplingShift, byte(value)))
}

func (b Config) InactiveDuration() InactiveDuration {
	return InactiveDuration(getValue(byte(b), InactiveDurationSize, InactiveDurationShift))
}

func (b Config) SetInactiveDuration(value InactiveDuration) Config {
	return Config(setValue(byte(b), InactiveDurationSize, InactiveDurationShift, byte(value)))
}

func (b Config) String() string {
	return fmt.Sprintf(
		"inactive duration: %s, filter coefficient: %s",
		b.InactiveDuration(),
		b.FilterCoefficient(),
	)
}

type CtrlHum byte

func (b CtrlHum) HumidityOversampling() HumidityOversampling {
	return HumidityOversampling(getValue(byte(b), HumidityOversamplingSize, HumidityOversamplingShift))
}

func (b CtrlHum) SetHumidityOversampling(value HumidityOversampling) CtrlHum {
	return CtrlHum(setValue(byte(b), HumidityOversamplingSize, HumidityOversamplingShift, byte(value)))
}

func (b CtrlHum) String() string {
	return fmt.Sprintf(
		"humidity oversampling: %s",
		b.HumidityOversampling(),
	)
}

type CtrlMeas byte

func (b CtrlMeas) Mode() Mode {
	return Mode(getValue(byte(b), ModeSize, ModeShift))
}

func (b CtrlMeas) PressureOversampling() PressureOversampling {
	return PressureOversampling(getValue(byte(b), PressureOversamplingSize, PressureOversamplingShift))
}

func (b CtrlMeas) SetMode(value Mode) CtrlMeas {
	return CtrlMeas(setValue(byte(b), ModeSize, ModeShift, byte(value)))
}

func (b CtrlMeas) SetPressureOversampling(value PressureOversampling) CtrlMeas {
	return CtrlMeas(setValue(byte(b), PressureOversamplingSize, PressureOversamplingShift, byte(value)))
}

func (b CtrlMeas) SetTemperatureOversampling(value TemperatureOversampling) CtrlMeas {
	return CtrlMeas(setValue(byte(b), TemperatureOversamplingSize, TemperatureOversamplingShift, byte(value)))
}

func (b CtrlMeas) TemperatureOversampling() TemperatureOversampling {
	return TemperatureOversampling(getValue(byte(b), TemperatureOversamplingSize, TemperatureOversamplingShift))
}

func (b CtrlMeas) String() string {
	return fmt.Sprintf(
		"temperature oversampling: %s, pressure oversampling: %s, mode: %s",
		b.TemperatureOversampling(),
		b.PressureOversampling(),
		b.Mode(),
	)
}

type FilterCoefficient byte

func (b FilterCoefficient) String() string {
	if b > 0 {
		return strconv.Itoa(1 << b)
	}

	return "off"
}

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

type HumidityOversampling byte

func (b HumidityOversampling) String() string {
	return oversamplingString(oversamplingValue(byte(b)))
}

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

func (b I2CAddress) String() string {
	return fmt.Sprintf("0x%02X", byte(b))
}

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

func (d InactiveDuration) String() string {
	return fmt.Sprintf("%.2fms", d.Milliseconds())
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

type PressureOversampling byte

func (b PressureOversampling) String() string {
	return oversamplingString(oversamplingValue(byte(b)))
}

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

type Mode byte

const (
	ModeShift = 0
	ModeSize  = 2

	ModeSleep  Mode = 0x00
	ModeForced Mode = 0x01
	ModeNormal Mode = 0x03
)

func (b Mode) String() string {
	switch b {
	case ModeSleep:
		return "sleep"
	case ModeForced:
		return "forced"
	case ModeNormal:
		return "normal"
	default:
		return fmt.Sprintf("unknown run mode: 0x%02X", byte(b))
	}
}

type TemperatureOversampling byte

func (b TemperatureOversampling) String() string {
	return oversamplingString(oversamplingValue(byte(b)))
}

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

func oversamplingString(n int) string {
	return fmt.Sprintf("x%d", n)
}

func oversamplingValue(b byte) int {
	if b > 0 {
		return 1 << (b - 1)
	}

	return 0
}
