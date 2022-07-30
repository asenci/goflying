package altimeter

import (
	"math"
)

// Altimeter calculates the altitude based on the pressure and altimeter
// setting/QNH using the National Weather Service approximation (with units
// converted to millibar)
func Altimeter(pressure, qnh Millibar) Feet {
	return Feet((math.Pow(float64(qnh), 0.190263) - math.Pow(float64(pressure), 0.190263)) / 0.00002569)
}

// DensityAltitude calculates the density altitude using the National Weather
// Service approximation (with units converted to millibar and kelvin)
func DensityAltitude(pressure Millibar, temperature Celsius) Feet {
	return Feet(145442.16 * (1 - math.Pow((0.28424266*float64(pressure))/float64(temperature.ToKelvin()), 0.235)))
}

// DensityAltitudeWet calculates the density altitude using the relative
// humidity (https://wahiduddin.net/calc/density_altitude.htm equation 12)
func DensityAltitudeWet(pressure Millibar, temperature Celsius, humidity RelativeHumidity) Feet {
	return Kilometers(44.3308 - (42.2665 * math.Pow(airDensity(pressure, temperature, humidity), 0.234969))).ToFeet()
}

// PressureAltitude calculates the pressure altitude using the National Weather
// Service approximation
func PressureAltitude(pressure Millibar) Feet {
	return Feet(145442.16 * (1.0 - math.Pow(float64(pressure)/1013.25, 0.190263)))
}

// airDensity calculates the air density in kg/m3
// (https://wahiduddin.net/calc/density_altitude.htm equation 4b)
func airDensity(pressure Millibar, temperature Celsius, humidity RelativeHumidity) float64 {
	return (float64(pressure) / (2.8705 * float64(temperature.ToKelvin()))) * (1 - ((0.378 * float64(vaporPressure(temperature, humidity))) / float64(pressure)))
}

// saturationVaporPressure calculates the saturation vapor pressure using
// Tetens' Formula (https://wahiduddin.net/calc/density_altitude.htm equation 6)
func saturationVaporPressure(temperature Celsius) Millibar {
	return Millibar(6.1078 * math.Pow(10, (7.5*float64(temperature))/(237.3+float64(temperature))))
}

// vaporPressure calculates the vapor pressure based on the relative humidity
// (https://wahiduddin.net/calc/density_altitude.htm equation 7b)
func vaporPressure(temperature Celsius, humidity RelativeHumidity) Millibar {
	return Millibar((float64(humidity) / 100) * float64(saturationVaporPressure(temperature)))
}

type Celsius float64

func (t Celsius) ToKelvin() Kelvin {
	return Kelvin(t + 273.15)
}

type Feet float64

type Kelvin float64

type Kilometers float64

func (d Kilometers) ToFeet() Feet {
	return Feet(d * 3280.84)
}

type Millibar float64

type RelativeHumidity float64
