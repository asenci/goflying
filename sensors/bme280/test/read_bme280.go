package main

import (
	"fmt"
	"time"

	"github.com/kidoman/embd"

	"github.com/westphae/goflying/sensors/bme280"
)

func main() {
	i2cbus := embd.NewI2CBus(1)

	var bmes []*bme280.BME280
	defer func() {
		for _, bme := range bmes {
			bme.Close()
		}
	}()

	for i, address := range []bme280.I2CAddress{bme280.I2CAddressLow, bme280.I2CAddressHigh} {
		bme, err := bme280.NewBME280(i2cbus, address, func(configByte bme280.ConfigByte, controlByte bme280.ControlByte, humControlByte bme280.HumidityControlByte) (bme280.ConfigByte, bme280.ControlByte, bme280.HumidityControlByte, error) {
			configByte.SetFilterCoefficient(bme280.FilterCoefficient16)
			configByte.SetInactiveDuration(bme280.InactiveDuration62_5ms)

			controlByte.SetRunMode(bme280.RunModeNormal)
			controlByte.SetPressureOversampling(bme280.PressureOversampling1x)
			controlByte.SetTemperatureOversampling(bme280.TemperatureOversampling1x)

			humControlByte.SetHumidityOversampling(bme280.HumidityOversampling1x)

			return configByte, controlByte, humControlByte, nil
		})

		if err != nil {
			fmt.Printf("no BME280 at address %d: %s\n", i, err)
			continue
		}

		bmes = append(bmes, bme)
	}

	if len(bmes) == 0 {
		return
	}

	fmt.Println("t,chip,temp,press,hum,alt")
	delay, _ := bmes[0].MeasurementDuration(true)
	clock := time.NewTicker(delay)
	for {
		t := <-clock.C
		for _, bme := range bmes {
			fmt.Printf("%s,0x%02X,%.2f,%.2f,%.2f,%.1f\n", t.Format(time.StampMilli), bme.I2CAddress(), bme.Data.Temperature(), bme.Data.Pressure(), bme.Data.Humidity(), bme.Altitude())
		}
	}
}
