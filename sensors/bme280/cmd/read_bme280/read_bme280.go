package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kidoman/embd"

	"github.com/westphae/goflying"
	"github.com/westphae/goflying/sensors/bme280"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	i2cbus := &goflying.I2CBus{I2CBus: embd.NewI2CBus(1)}

	var bmes []*bme280.Sensor

	var stopFuncs []func()
	defer func() {
		for _, f := range stopFuncs {
			f()
		}
	}()

	for _, address := range []bme280.I2CAddress{bme280.I2CAddressLow, bme280.I2CAddressHigh} {
		bme, err := bme280.NewSensor(i2cbus, address,
			bme280.WithFilterCoefficient(bme280.FilterCoefficient16),
			bme280.WithInactiveDuration(bme280.InactiveDuration62_5ms),
			bme280.WithPressureOversampling(bme280.PressureOversampling1x),
			bme280.WithTemperatureOversampling(bme280.TemperatureOversampling1x),
			bme280.WithHumidityOversampling(bme280.HumidityOversampling1x),
		)

		if err != nil {
			fmt.Printf("no sensor at address %s: %s\n", address, err)
			continue
		}

		bmes = append(bmes, bme)

		stopFunc, err := bme.Start(ctx)
		if err != nil {
			fmt.Printf("error starting sensor polling: %s", err)
			continue
		}

		stopFuncs = append(stopFuncs, stopFunc)
	}

	if len(bmes) == 0 {
		return
	}

	delay, _ := bmes[0].MeasurementDuration(true)
	clock := time.NewTicker(delay)
	defer clock.Stop()

	fmt.Println("timestamp, chip address, temperature, pressure, humidity")
	for {
		select {
		case <-clock.C:
			for _, bme := range bmes {
				fmt.Printf("%s, %s, %s, %s, %s\n", bme.Data.Timestamp().Format(time.StampMilli), bme.I2CAddress(), bme.Data.Temperature(), bme.Data.Pressure(), bme.Data.Humidity())
			}
		case <-ctx.Done():
			break
		}
	}
}
