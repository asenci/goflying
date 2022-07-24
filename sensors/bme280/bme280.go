/*
Reference 1: https://github.com/BoschSensortec/BME280_driver
Reference 2: https://github.com/adafruit/Adafruit_CircuitPython_BME280
*/

package bme280

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	_ "github.com/kidoman/embd/host/rpi"

	"github.com/westphae/goflying/sensors"
)

const (
	QNH            = 1013.25          // Sea level reference pressure in hPa
	BufSize        = 256              // Buffer size for reading data from BME
	ExtraReadDelay = time.Millisecond // Delay between chip reading polls
)

type BME280 struct {
	i2CAddress byte
	i2CBus     embd.I2CBus

	Data *MeasurementData

	C      <-chan *sensors.BMEData
	CBuf   <-chan *sensors.BMEData
	cClose chan bool
}

/*
NewBME280 returns a BME280 object with the chosen settings:
address is one of bme280.I2CAddressLow (0x76) or bme280.I2CAddressHigh (0x77).
powerMode is one of bme280.ModeSleep, bme280.ModeForced, or bme280.ModeNormal.
standby is one of the bme280.StandbyTimeX (0.5ms up to 1000ms).
filter is one of bme280.FilterCoeffX.
tempRes is one of bme280.OversampX (up to 16x).
presRes is one of bme280.OversampX (up to 16x).
humRes is one of bme280.OversampX (up to 16x).
See BME280 datasheet for details.
*/
func NewBME280(i2cbus embd.I2CBus, address I2CAddress, settings ...SettingFunc) (*BME280, error) {
	bme := new(BME280)
	bme.i2CBus = i2cbus
	bme.i2CAddress = byte(address)
	bme.Data = &MeasurementData{}

	// Make sure we can connect to the chip and read a valid ChipID
	chipID, err := bme.ChipID()
	if err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	if chipID != ChipID {
		return nil, fmt.Errorf("bme280: wrong ChipID: expected %02X, got %02X", ChipID, chipID)
	}

	if err := bme.Reset(); err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	if err := bme.SetConfiguration(settings...); err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	if err := bme.ReadCalibrationData(); err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	go bme.Run()

	return bme, nil
}

func (bme *BME280) Close() error {
	bme.cClose <- true

	if err := bme.SetRunMode(RunModeSleep); err != nil {
		return fmt.Errorf("bme280: %w", err)
	}

	return nil
}

func (bme *BME280) Run() error {
	measurementDuration, err := bme.MeasurementDuration(false)
	if err != nil {
		return err
	}

	clock := time.NewTicker(measurementDuration)
	defer clock.Stop()

	rawData := make([]byte, MeasurementDataSize)

	for {
		select {
		case <-clock.C:
			if err := bme.i2CBus.ReadFromReg(bme.i2CAddress, RegisterPressDataMSB, rawData); err != nil {
				log.Printf("bme280: error reading sensor data: %s", err)
				continue
			}

			bme.Data.Update(rawData)

		case <-bme.cClose: // Stop the goroutine, ease up on the CPU
			break
		}
	}
}

func (bme *BME280) Altitude() float64 {
	return 145366.45 * (1.0 - math.Pow(bme.Data.Pressure()/QNH, 0.190284))
}

func (bme *BME280) ChipID() (byte, error) {
	chipID, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterChipID)
	if err != nil {
		return 0, fmt.Errorf("failed to read chip ID:  %w", err)
	}

	return chipID, nil
}

func (bme *BME280) ConfigByte() (ConfigByte, error) {
	value, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterConfig)

	if err != nil {
		return 0, fmt.Errorf("failed to read from Config register")
	}

	return ConfigByte(value), nil
}

func (bme *BME280) ControlByte() (ControlByte, error) {
	value, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterControl)

	if err != nil {
		return 0, fmt.Errorf("failed to read from Control register")
	}

	return ControlByte(value), nil
}

func (bme *BME280) FilterCoefficient() (FilterCoefficient, error) {
	configByte, err := bme.ConfigByte()
	if err != nil {
		return 0, err
	}

	return configByte.FilterCoefficient(), nil
}

func (bme *BME280) HumControlByte() (HumidityControlByte, error) {
	value, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterHumControl)

	if err != nil {
		return 0, fmt.Errorf("failed to read from Humidity Control register")
	}

	return HumidityControlByte(value), nil
}

func (bme *BME280) HumidityOversampling() (HumidityOversampling, error) {
	humControlByte, err := bme.HumControlByte()
	if err != nil {
		return 0, err
	}

	return humControlByte.HumidityOversampling(), nil
}

func (bme *BME280) I2CAddress() I2CAddress {
	return I2CAddress(bme.i2CAddress)
}

func (bme *BME280) InactiveDuration() (InactiveDuration, error) {
	configByte, err := bme.ConfigByte()
	if err != nil {
		return 0, err
	}

	return configByte.InactiveDuration(), nil
}

func (bme *BME280) MeasurementDuration(max bool) (time.Duration, error) {
	multiplier := 1.0
	totalDuration := 1.0

	if max {
		multiplier = 1.15
		totalDuration = 1.25
	}

	inactiveDuration, err := bme.InactiveDuration()
	if err != nil {
		return 0, err
	}

	totalDuration += inactiveDuration.Milliseconds()

	temperatureOversampling, err := bme.TemperatureOversampling()
	if err != nil {
		return 0, err
	}

	if temperatureOversampling > 0 {
		totalDuration += 2 * multiplier * float64(temperatureOversampling.Value())
	}

	pressureOversampling, err := bme.PressureOversampling()
	if err != nil {
		return 0, err
	}

	if pressureOversampling > 0 {
		totalDuration += 0.5 + (2 * multiplier * float64(pressureOversampling.Value()))
	}

	humidityOversampling, err := bme.HumidityOversampling()
	if err != nil {
		return 0, err
	}

	if humidityOversampling > 0 {
		totalDuration += 0.5 + (2 * multiplier * float64(humidityOversampling.Value()))
	}

	return time.Duration(totalDuration*1000) * time.Microsecond, nil
}

func (bme *BME280) PressureOversampling() (PressureOversampling, error) {
	controlByte, err := bme.ControlByte()
	if err != nil {
		return 0, err
	}

	return controlByte.PressureOversampling(), nil
}

func (bme *BME280) ReadCalibrationData() error {
	rawData := make([]byte, CalibrationDataSize)
	if err := bme.i2CBus.ReadFromReg(bme.i2CAddress, RegisterCalibrationData, rawData); err != nil {
		return fmt.Errorf("failed to read calibration data: %w", err)
	}

	rawHumData := make([]byte, CalibrationHumDataSize)
	if err := bme.i2CBus.ReadFromReg(bme.i2CAddress, RegisterHumCalibrationData, rawHumData); err != nil {
		return fmt.Errorf("failed to read humidity calibration data: %w", err)
	}

	bme.Data.SetCalibrationData(NewCalibrationData(rawData, rawHumData))

	return nil
}

func (bme *BME280) Reset() error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterSoftReset, ResetCode); err != nil {
		return fmt.Errorf("failed to reset sensor: %w", err)
	}

	// Wait start-up time
	time.Sleep(2 * time.Millisecond)

	return nil
}

func (bme *BME280) RunMode() (RunMode, error) {
	controlByte, err := bme.ControlByte()
	if err != nil {
		return 0, err
	}

	return controlByte.RunMode(), nil
}

func (bme *BME280) SetConfiguration(settings ...SettingFunc) error {
	var err error

	configByte, err := bme.ConfigByte()
	if err != nil {
		return err
	}

	newConfigByte := configByte

	controlByte, err := bme.ControlByte()
	if err != nil {
		return err
	}

	newControlByte := controlByte

	humControlByte, err := bme.HumControlByte()
	if err != nil {
		return err
	}

	newHumControlByte := humControlByte

	for _, f := range settings {
		newConfigByte, newControlByte, newHumControlByte, err = f(newConfigByte, newControlByte, newHumControlByte)
		if err != nil {
			return fmt.Errorf("error setting sensor configuration: %w", err)
		}
	}

	if newConfigByte != configByte {
		if err := bme.SetConfigByte(newConfigByte); err != nil {
			return err
		}
	}

	if newControlByte != controlByte {
		if err := bme.SetControlByte(newControlByte); err != nil {
			return err
		}
	}

	if newHumControlByte != humControlByte {
		if err := bme.SetHumidityControlByte(newHumControlByte); err != nil {
			return err
		}
	}

	return nil
}

func (bme *BME280) SetConfigByte(value ConfigByte) error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterConfig, byte(value)); err != nil {
		return fmt.Errorf("failed to write to Config register")
	}

	return nil
}

func (bme *BME280) SetControlByte(value ControlByte) error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterControl, byte(value)); err != nil {
		return fmt.Errorf("failed to write to Control register")
	}

	return nil
}

func (bme *BME280) SetFilterCoefficient(value FilterCoefficient) error {
	oldByte, err := bme.ConfigByte()
	if err != nil {
		return err
	}

	newByte := oldByte.SetFilterCoefficient(value)

	return bme.SetConfigByte(newByte)
}

func (bme *BME280) SetHumidityOversampling(value HumidityOversampling) error {
	oldByte, err := bme.HumControlByte()
	if err != nil {
		return err
	}

	newByte := oldByte.SetHumidityOversampling(value)

	return bme.SetHumidityControlByte(newByte)
}

func (bme *BME280) SetHumidityControlByte(value HumidityControlByte) error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterHumControl, byte(value)); err != nil {
		return fmt.Errorf("failed to write to Humidity Control register")
	}

	return nil
}

func (bme *BME280) SetInactiveDuration(value InactiveDuration) error {
	oldByte, err := bme.ConfigByte()
	if err != nil {
		return err
	}

	newByte := oldByte.SetInactiveDuration(value)

	return bme.SetConfigByte(newByte)
}

func (bme *BME280) SetPressureOversampling(value PressureOversampling) error {
	oldByte, err := bme.ControlByte()
	if err != nil {
		return err
	}

	newByte := oldByte.SetPressureOversampling(value)

	return bme.SetControlByte(newByte)
}

func (bme *BME280) SetRunMode(value RunMode) error {
	oldByte, err := bme.ControlByte()
	if err != nil {
		return err
	}

	newByte := oldByte.SetRunMode(value)

	return bme.SetControlByte(newByte)
}

func (bme *BME280) SetTemperatureOversampling(value TemperatureOversampling) error {
	oldByte, err := bme.ControlByte()
	if err != nil {
		return err
	}

	newByte := oldByte.SetTemperatureOversampling(value)

	return bme.SetControlByte(newByte)
}

func (bme *BME280) TemperatureOversampling() (TemperatureOversampling, error) {
	controlByte, err := bme.ControlByte()
	if err != nil {
		return 0, err
	}

	return controlByte.TemperatureOversampling(), nil
}
