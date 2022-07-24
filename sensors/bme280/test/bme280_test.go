package main

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/kidoman/embd"

	"github.com/westphae/goflying/sensors/bme280"
)

func TestBME280Math(t *testing.T) {
	bme := bme280.BME280{
		DigT: map[int]int32{
			1: 27504,
			2: 26435,
			3: -1000,
		},
		DigP: map[int]int64{
			1: 36477,
			2: -10685,
			3: 3024,
			4: 2855,
			5: 140,
			6: -7,
			7: 15500,
			8: -14600,
			9: 6000,
		},
	}

	raw_temp := int32(519888)
	raw_press := int64(415148)
	temp := bme.CalcCompensatedTemp(raw_temp)
	press := bme.CalcCompensatedPress(raw_press)

	t_fine := int32(128422)
	if bme.T_fine != t_fine {
		log.Printf("t_fine mismatch: calculated %d, should be %d\n", bme.T_fine, t_fine)
		t.Fail()
	} else {
		log.Printf("t_fine matched:  calculated %d, should be %d\n", bme.T_fine, t_fine)
	}

	calc_temp := 25.08
	if math.Abs(temp-calc_temp) > 0.01 {
		log.Printf("temp mismatch:   calculated %f, should be %f\n", temp, calc_temp)
		t.Fail()
	} else {
		log.Printf("temp matched:    calculated %f, should be %f\n", temp, calc_temp)
	}

	calc_press := 1006.5327
	if math.Abs(press-calc_press) > 0.001 {
		log.Printf("press mismatch:  calculated %f, should be %f\n", press, calc_press)
		t.Fail()
	} else {
		log.Printf("press matched:   calculated %f, should be %f\n", press, calc_press)
	}

}

func TestBME280Setup(t *testing.T) {
	const (
		mode          = bme280.NormalMode
		standbyTime   = bme280.StandbyTime250ms
		filterCoeff   = bme280.FilterCoeff8
		oversampTemp  = bme280.Oversamp16x
		oversampPress = bme280.Oversamp16x
	)

	var (
		modes = []byte{
			bme280.SleepMode,
			// bme280.ForcedMode,
			bme280.NormalMode,
		}
		standbyTimes = []byte{
			bme280.StandbyTime1ms,
			bme280.StandbyTime63ms,
			bme280.StandbyTime125ms,
			bme280.StandbyTime250ms,
			bme280.StandbyTime500ms,
			bme280.StandbyTime1000ms,
			bme280.StandbyTime2000ms,
			bme280.StandbyTime4000ms,
		}
		filterCoeffs = []byte{
			bme280.FilterCoeffOff,
			bme280.FilterCoeff2,
			bme280.FilterCoeff4,
			bme280.FilterCoeff8,
			bme280.FilterCoeff16,
		}
		oversamps = []byte{
			bme280.OversampSkipped,
			bme280.Oversamp1x,
			bme280.Oversamp2x,
			bme280.Oversamp4x,
			bme280.Oversamp8x,
			bme280.Oversamp16x,
		}
		bme *bme280.BME280
		err error
	)

	var checkAll = func(newMode, newStandbyTime, newFilterCoeff, newOversampTemp, newOversampPress byte) {
		var (
			curMode, curStandbyTime, curFilterCoeff, curOversampTemp, curOversampPress byte
			err                                                                        error
		)
		curMode, err = bme.GetPowerMode()
		if err != nil {
			fmt.Printf("Error getting power mode: %s\n", err)
		}
		if curMode != newMode {
			fmt.Printf("Mode not set correctly: got %x, should be %x\n", curMode, newMode)
			t.Fail()
		}
		curStandbyTime, err = bme.GetStandbyTime()
		if err != nil {
			fmt.Printf("Error getting standby time: %s\n", err)
		}
		if curStandbyTime != newStandbyTime {
			fmt.Printf("Standby time not set correctly: got %x, should be %x\n", curStandbyTime, newStandbyTime)
			t.Fail()
		}
		curFilterCoeff, err = bme.GetFilterCoeff()
		if err != nil {
			fmt.Printf("Error getting filter coefficient: %s\n", err)
		}
		if curFilterCoeff != newFilterCoeff {
			fmt.Printf("Filter coefficient not set correctly: got %x, should be %x\n", curFilterCoeff, newFilterCoeff)
			t.Fail()
		}
		curOversampTemp, err = bme.GetOversampTemp()
		if err != nil {
			fmt.Printf("Error getting temperature oversampling: %s\n", err)
		}
		if curOversampTemp != newOversampTemp {
			fmt.Printf("Temperature oversampling not set correctly: got %x, should be %x\n", curOversampTemp, newOversampTemp)
			t.Fail()
		}
		curOversampPress, err = bme.GetOversampPress()
		if err != nil {
			fmt.Printf("Error getting pressure oversampling: %s\n", err)
		}
		if curOversampPress != newOversampPress {
			fmt.Printf("Pressure oversampling not set correctly: got %x, should be %x\n", curOversampPress, newOversampPress)
			t.Fail()
		}
	}

	i2cbus := embd.NewI2CBus(1)
	bme, err = bme280.NewBME280(&i2cbus, bme280.Address1, mode, standbyTime, filterCoeff, oversampTemp, oversampPress)
	if err != nil {
		bme, err = bme280.NewBME280(i2cbus, bme280.Address2, mode, standbyTime, filterCoeff, oversampTemp, oversampPress)
	}
	if err != nil {
		log.Println("Couldn't find a BME280")
		t.Fail()
		return
	}

	// Test all initial gets
	fmt.Println("Checking initial setup values")
	checkAll(mode, standbyTime, filterCoeff, oversampTemp, oversampPress)

	// Test changing modes
	fmt.Println("\nChecking mode changes")
	for _, x := range modes {
		fmt.Printf("Setting mode to %x\n", x)
		err = bme.SetPowerMode(x)
		if err != nil {
			fmt.Printf("Error setting power mode to %x: %s\n", x, err)
			t.Fail()
			continue
		}
		checkAll(x, standbyTime, filterCoeff, oversampTemp, oversampPress)
	}
	bme.SetPowerMode(mode)

	// Test changing standby times
	fmt.Println("\nChecking standby time changes")
	for _, x := range standbyTimes {
		fmt.Printf("Setting standby time to %x\n", x)
		err = bme.SetStandbyTime(x)
		if err != nil {
			fmt.Printf("Error setting standby time to %x: %s\n", x, err)
			t.Fail()
			continue
		}
		checkAll(mode, x, filterCoeff, oversampTemp, oversampPress)
	}
	bme.SetStandbyTime(standbyTime)

	// Test changing filter coefficients
	fmt.Println("\nChecking filter coefficient changes")
	for _, x := range filterCoeffs {
		fmt.Printf("Setting filter coefficient to %x\n", x)
		err = bme.SetFilterCoeff(x)
		if err != nil {
			fmt.Printf("Error setting filter coefficient to %x: %s\n", x, err)
			t.Fail()
			continue
		}
		checkAll(mode, standbyTime, x, oversampTemp, oversampPress)
	}
	bme.SetFilterCoeff(filterCoeff)

	// Test changing temperature oversampling
	fmt.Println("\nChecking temperature oversampling changes")
	for _, x := range oversamps {
		fmt.Printf("Setting temperature oversampling to %x\n", x)
		err = bme.SetOversampTemp(x)
		if err != nil {
			fmt.Printf("Error setting temperature oversampling to %x: %s\n", x, err)
			t.Fail()
			continue
		}
		checkAll(mode, standbyTime, filterCoeff, x, oversampPress)
	}
	bme.SetOversampTemp(oversampTemp)

	// Test changing pressure oversampling
	fmt.Println("\nChecking pressure oversampling changes")
	for _, x := range oversamps {
		fmt.Printf("Setting pressure oversampling to %x\n", x)
		err = bme.SetOversampPress(x)
		if err != nil {
			fmt.Printf("Error setting pressure oversampling to %x: %s\n", x, err)
			t.Fail()
			continue
		}
		checkAll(mode, standbyTime, filterCoeff, oversampTemp, x)
	}
	bme.SetOversampPress(oversampPress)
}
