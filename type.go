package goflying

import "fmt"

type Celsius float64

func (t Celsius) String() string {
	return fmt.Sprintf("%.2fºC", t)
}

func (t Celsius) ToKelvin() Kelvin {
	return Kelvin(t + 273.15)
}

type Feet float64

func (d Feet) String() string {
	return fmt.Sprintf("%.2fft", d)
}

type Kelvin float64

func (t Kelvin) String() string {
	return fmt.Sprintf("%.2fK", t)
}

type Kilometers float64

func (d Kilometers) String() string {
	return fmt.Sprintf("%.2fkm", d)
}

func (d Kilometers) ToFeet() Feet {
	return Feet(d * 3280.84)
}

type HPa float64

func (p HPa) String() string {
	return fmt.Sprintf("%.2fhPa", p)
}

type RelativeHumidity float64

func (f RelativeHumidity) String() string {
	return fmt.Sprintf("%.2f%%", f)
}
