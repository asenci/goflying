/*
Reference 1: https://github.com/BoschSensortec/BME280_driver
Reference 2: https://github.com/adafruit/Adafruit_CircuitPython_BME280
*/

package bme280

import (
	"context"
	"fmt"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"

	"github.com/westphae/goflying"
)

type Sensor struct {
	i2CAddress byte
	i2CBus     embd.I2CBus

	Data *MeasurementData
}

// NewSensor returns a Sensor object connected to the specified I2C bus and with
// the specified I2C address. One or more SettingFunc functions can be specified
// for updating the configuration and control bytes during initialisation.
func NewSensor(bus embd.I2CBus, address I2CAddress, settings ...SettingFunc) (*Sensor, error) {
	bme := &Sensor{
		i2CBus:     bus,
		i2CAddress: byte(address),
	}

	// Make sure we can connect to the chip and read a valid ChipID
	chipID, err := bme.ChipID()
	if err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	if chipID != ChipID {
		return nil, fmt.Errorf("bme280: wrong ChipID: expected 0x%02X, got 0x%02X", ChipID, chipID)
	}

	goflying.Debugf("bme280: chip ID: 0x%02X\n", chipID)

	cal, err := bme.CalibrationData()
	if err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	goflying.Debugf("bme280: calibration data: %s\n", cal)

	bme.Data = NewMeasurementData(cal)

	goflying.Debugln("bme280: resetting chip")

	if err := bme.Reset(); err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	goflying.Debugln("bme280: adjusting config and control bytes")

	if err := bme.Configure(settings...); err != nil {
		return nil, fmt.Errorf("bme280: %w", err)
	}

	return bme, nil
}

// Configure runs all the specified SettingFunc and writes the modified config
// and control bytes after they have all finished running
func (bme *Sensor) Configure(settings ...SettingFunc) error {
	var err error

	ctrlHum, err := bme.CtrlHum()
	if err != nil {
		return err
	}
	goflying.Debugf("bme280: initial ctrl_hum: %s\n", ctrlHum)

	newCtrlHum := ctrlHum

	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return err
	}
	goflying.Debugf("bme280: initial ctrl_meas: %s\n", ctrlMeas)

	newCtrlMeas := ctrlMeas

	config, err := bme.Config()
	if err != nil {
		return err
	}
	goflying.Debugf("bme280: initial config: %s\n", config)

	newConfig := config

	for _, f := range settings {
		newConfig, newCtrlMeas, newCtrlHum, err = f(newConfig, newCtrlMeas, newCtrlHum)
		if err != nil {
			return fmt.Errorf("error setting sensor configuration: %w", err)
		}
	}

	if newCtrlHum != ctrlHum {
		goflying.Debugf("bme280: writing new ctrl_hum: %s\n", newCtrlHum)

		if err := bme.SetCtrlHum(newCtrlHum); err != nil {
			return err
		}
	}

	if newCtrlMeas != ctrlMeas {
		goflying.Debugf("bme280: writing new ctrl_meas: %s\n", newCtrlMeas)

		if err := bme.SetCtrlMeas(newCtrlMeas); err != nil {
			return err
		}
	}

	if newConfig != config {
		goflying.Debugf("bme280: writing new config: %s\n", newConfig)

		if err := bme.SetConfig(newConfig); err != nil {
			return err
		}
	}

	return nil
}

// Start starts the polling goroutine and returns a stop function
func (bme *Sensor) Start(ctx context.Context) (func(), error) {
	measurementDuration, err := bme.MeasurementDuration(false)
	if err != nil {
		return nil, err
	}

	poolCtx, cancel := context.WithCancel(ctx)

	go bme.poll(poolCtx, measurementDuration)

	return cancel, nil
}

// poll sets the chip to normal mode and read the new measurement values
// according to the specified frequency until the ctx context is stopped
func (bme *Sensor) poll(ctx context.Context, frequency time.Duration) {
	defer func(bme *Sensor, value Mode) {
		goflying.Debugln("bme280: setting chip mode to Sleep")
		if err := bme.SetMode(value); err != nil {
			goflying.Logger.Printf("bme280: error setting chip mode to Sleep: %w\n", err)
		}
	}(bme, ModeSleep)

	goflying.Debugln("bme280: setting chip mode to Normal")
	if err := bme.SetMode(ModeNormal); err != nil {
		goflying.Logger.Printf("bme280: error setting chip mode to Normal: %w\n", err)
	}

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	rawData := make([]byte, MeasurementDataSize)

	for {
		select {
		case <-ticker.C:
			goflying.Debugln("bme280: reading measurement data")

			if err := bme.i2CBus.ReadFromReg(bme.i2CAddress, RegisterPressDataMSB, rawData); err != nil {
				goflying.Logger.Printf("bme280: error reading sensor data: %w\n", err)
				continue
			}

			bme.Data.Update(rawData)

			goflying.Debugf("bme280: new measurement data: %s\n", bme.Data)

		case <-ctx.Done():
			break
		}
	}
}

// CalibrationData reads chip factory calibration data from registers
// RegisterCalibrationData and RegisterHumCalibrationData
func (bme *Sensor) CalibrationData() (*CalibrationData, error) {
	calibrationData := make([]byte, CalibrationDataSize)
	if err := bme.i2CBus.ReadFromReg(bme.i2CAddress, RegisterCalibrationData, calibrationData); err != nil {
		return nil, fmt.Errorf("failed to read calibration data: %w", err)
	}

	humidityCalibrationData := make([]byte, CalibrationHumDataSize)
	if err := bme.i2CBus.ReadFromReg(bme.i2CAddress, RegisterHumCalibrationData, humidityCalibrationData); err != nil {
		return nil, fmt.Errorf("failed to read humidity calibration data: %w", err)
	}

	return NewCalibrationData(calibrationData, humidityCalibrationData), nil
}

// ChipID reads chip ID from register RegisterChipID. Expected value is 0x60 for
// BME280
func (bme *Sensor) ChipID() (byte, error) {
	chipID, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterChipID)
	if err != nil {
		return 0, fmt.Errorf("failed to read chip ID: %w", err)
	}

	return chipID, nil
}

// Config reads the config register byte from RegisterConfig. The config
// register contains the chip inactive duration and IIR filter coefficient
// values
func (bme *Sensor) Config() (Config, error) {
	value, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterConfig)

	if err != nil {
		return 0, fmt.Errorf("failed to read from Config register")
	}

	return Config(value), nil
}

// CtrlHum reads the ctrl_hum register byte from RegisterCtrlHum. The ctrl_hum
// register contains the humidity oversampling value
func (bme *Sensor) CtrlHum() (CtrlHum, error) {
	value, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterCtrlHum)

	if err != nil {
		return 0, fmt.Errorf("failed to read from Humidity Control register")
	}

	return CtrlHum(value), nil
}

// CtrlMeas reads the ctrl_meas register byte from RegisterCtrlMeas. The
// ctrl_meas register contains the pressure and temperature oversampling, and
// the chip mode values
func (bme *Sensor) CtrlMeas() (CtrlMeas, error) {
	value, err := bme.i2CBus.ReadByteFromReg(bme.i2CAddress, RegisterCtrlMeas)

	if err != nil {
		return 0, fmt.Errorf("failed to read from Control register")
	}

	return CtrlMeas(value), nil
}

// FilterCoefficient reads the IIR filter coefficient bits from the config
// register
func (bme *Sensor) FilterCoefficient() (FilterCoefficient, error) {
	config, err := bme.Config()
	if err != nil {
		return 0, err
	}

	return config.FilterCoefficient(), nil
}

// HumidityOversampling reads the humidity oversampling bits from the ctrl_hum
// register
func (bme *Sensor) HumidityOversampling() (HumidityOversampling, error) {
	ctrlHum, err := bme.CtrlHum()
	if err != nil {
		return 0, err
	}

	return ctrlHum.HumidityOversampling(), nil
}

// I2CAddress returns the formatted chip I2C address
func (bme *Sensor) I2CAddress() I2CAddress {
	return I2CAddress(bme.i2CAddress)
}

// InactiveDuration reads the inactive duration bits from the config register
func (bme *Sensor) InactiveDuration() (InactiveDuration, error) {
	config, err := bme.Config()
	if err != nil {
		return 0, err
	}

	return config.InactiveDuration(), nil
}

// MeasurementDuration returns the typical interval between measurements if max
// is false, otherwise it will return the maximum interval between measurements.
// The total measurement interval consists in a configurable inactive interval
// (set on the config register) and a variable active interval which depends on
// the oversampling settings. See appendix B of the BME280 datasheet for more
// details.
func (bme *Sensor) MeasurementDuration(max bool) (time.Duration, error) {
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

// Mode reads the chip mode from the ctrl_meas register
func (bme *Sensor) Mode() (Mode, error) {
	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return 0, err
	}

	return ctrlMeas.Mode(), nil
}

// PressureOversampling reads the pressure oversampling bits from the ctrl_meas
// register
func (bme *Sensor) PressureOversampling() (PressureOversampling, error) {
	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return 0, err
	}

	return ctrlMeas.PressureOversampling(), nil
}

// Reset triggers a soft reset on the chip. This will also reset the config and
// control registers.
func (bme *Sensor) Reset() error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterSoftReset, ResetCode); err != nil {
		return fmt.Errorf("failed to reset sensor: %w", err)
	}

	// Wait start-up time
	time.Sleep(2 * time.Millisecond)

	return nil
}

// TemperatureOversampling reads the temperature oversampling bits from the
// ctrl_meas register
func (bme *Sensor) TemperatureOversampling() (TemperatureOversampling, error) {
	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return 0, err
	}

	return ctrlMeas.TemperatureOversampling(), nil
}

// SetConfig sets the config register byte on the RegisterConfig register
func (bme *Sensor) SetConfig(value Config) error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterConfig, byte(value)); err != nil {
		return fmt.Errorf("failed to write to Config register")
	}

	return nil
}

// SetCtrlHum sets the ctrl_hum register byte on the RegisterCtrlHum register.
// Changes to this register only become effective after a write operation to
// RegisterCtrlMeas.
func (bme *Sensor) SetCtrlHum(value CtrlHum) error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterCtrlHum, byte(value)); err != nil {
		return fmt.Errorf("failed to write to Humidity Control register")
	}

	return nil
}

// SetCtrlMeas sets the ctrl_meas register byte on the RegisterCtrlMeas register
func (bme *Sensor) SetCtrlMeas(value CtrlMeas) error {
	if err := bme.i2CBus.WriteByteToReg(bme.i2CAddress, RegisterCtrlMeas, byte(value)); err != nil {
		return fmt.Errorf("failed to write to Control register")
	}

	return nil
}

// SetFilterCoefficient sets the IIR filter coefficient bits on the config
// register
func (bme *Sensor) SetFilterCoefficient(value FilterCoefficient) error {
	config, err := bme.Config()
	if err != nil {
		return err
	}

	return bme.SetConfig(config.SetFilterCoefficient(value))
}

// SetHumidityOversampling sets the humidity oversampling bits on the ctrl_hum
// register
func (bme *Sensor) SetHumidityOversampling(value HumidityOversampling) error {
	ctrlHum, err := bme.CtrlHum()
	if err != nil {
		return err
	}

	return bme.SetCtrlHum(ctrlHum.SetHumidityOversampling(value))
}

// SetInactiveDuration sets the inactive duration bits on the config register
func (bme *Sensor) SetInactiveDuration(value InactiveDuration) error {
	config, err := bme.Config()
	if err != nil {
		return err
	}

	return bme.SetConfig(config.SetInactiveDuration(value))
}

// SetMode sets the chip measurement mode on the ctrl_meas register. Set to
// ModeNormal for polling mode, or ModeForced to perform a single measurement.
func (bme *Sensor) SetMode(value Mode) error {
	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return err
	}

	return bme.SetCtrlMeas(ctrlMeas.SetMode(value))
}

// SetPressureOversampling sets the pressure oversampling bits on the ctrl_meas
// register
func (bme *Sensor) SetPressureOversampling(value PressureOversampling) error {
	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return err
	}

	return bme.SetCtrlMeas(ctrlMeas.SetPressureOversampling(value))
}

// SetTemperatureOversampling sets the temperature oversampling bits on the
// ctrl_meas register
func (bme *Sensor) SetTemperatureOversampling(value TemperatureOversampling) error {
	ctrlMeas, err := bme.CtrlMeas()
	if err != nil {
		return err
	}

	return bme.SetCtrlMeas(ctrlMeas.SetTemperatureOversampling(value))
}

// SettingFunc represents a function that modifies one of the config or control
// registers and returns the modified or unmodified values for the next
// SettingFunc. A SettingFunc slice is passed to the NewSensor function to
// minimise the number of writes to the register during initialisation. See
// read_bme280.go for an example.
type SettingFunc func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error)

func WithFilterCoefficient(value FilterCoefficient) SettingFunc {
	return func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error) {
		return config.SetFilterCoefficient(value), ctrlMeas, ctrlHum, nil
	}
}

func WithHumidityOversampling(value HumidityOversampling) SettingFunc {
	return func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error) {
		return config, ctrlMeas, ctrlHum.SetHumidityOversampling(value), nil
	}
}

func WithInactiveDuration(value InactiveDuration) SettingFunc {
	return func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error) {
		return config.SetInactiveDuration(value), ctrlMeas, ctrlHum, nil
	}
}

func WithMode(value Mode) SettingFunc {
	return func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error) {
		return config, ctrlMeas.SetMode(value), ctrlHum, nil
	}
}

func WithPressureOversampling(value PressureOversampling) SettingFunc {
	return func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error) {
		return config, ctrlMeas.SetPressureOversampling(value), ctrlHum, nil
	}
}

func WithTemperatureOversampling(value TemperatureOversampling) SettingFunc {
	return func(config Config, ctrlMeas CtrlMeas, ctrlHum CtrlHum) (Config, CtrlMeas, CtrlHum, error) {
		return config, ctrlMeas.SetTemperatureOversampling(value), ctrlHum, nil
	}
}
