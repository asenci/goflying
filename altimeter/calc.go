package altimeter

import (
	"math"

	"github.com/westphae/goflying"
)

// Altimeter calculates the altitude based on the pressure and altimeter
// setting/QNH using the National Weather Service approximation (with units
// converted to hPa)
func Altimeter(pressure, qnh goflying.HPa) goflying.Feet {
	return goflying.Feet((math.Pow(float64(qnh), 0.190263) - math.Pow(float64(pressure), 0.190263)) / 0.00002569)
}

// DensityAltitude calculates the density altitude using the National Weather
// Service approximation (with units converted to hPa and kelvin)
func DensityAltitude(pressure goflying.HPa, temperature goflying.Celsius) goflying.Feet {
	return goflying.Feet(145442.16 * (1 - math.Pow((0.28424266*float64(pressure))/float64(temperature.ToKelvin()), 0.235)))
}

// DensityAltitudeWet calculates the density altitude using the relative
// humidity (https://wahiduddin.net/calc/density_altitude.htm equation 12)
func DensityAltitudeWet(pressure goflying.HPa, temperature goflying.Celsius, humidity goflying.RelativeHumidity) goflying.Feet {
	return goflying.Kilometers(44.3308 - (42.2665 * math.Pow(airDensity(pressure, temperature, humidity), 0.234969))).ToFeet()
}

// PressureAltitude calculates the pressure altitude using the National Weather
// Service approximation
func PressureAltitude(pressure goflying.HPa) goflying.Feet {
	return goflying.Feet(145442.16 * (1.0 - math.Pow(float64(pressure)/1013.25, 0.190263)))
}

// airDensity calculates the air density in kg/m3
// (https://wahiduddin.net/calc/density_altitude.htm equation 4b)
func airDensity(pressure goflying.HPa, temperature goflying.Celsius, humidity goflying.RelativeHumidity) float64 {
	return (float64(pressure) / (2.8705 * float64(temperature.ToKelvin()))) * (1 - ((0.378 * float64(vaporPressure(temperature, humidity))) / float64(pressure)))
}

// saturationVaporPressure calculates the saturation vapor pressure using
// Tetens' Formula (https://wahiduddin.net/calc/density_altitude.htm equation 6)
func saturationVaporPressure(temperature goflying.Celsius) goflying.HPa {
	return goflying.HPa(6.1078 * math.Pow(10, (7.5*float64(temperature))/(237.3+float64(temperature))))
}

// vaporPressure calculates the vapor pressure based on the relative humidity
// (https://wahiduddin.net/calc/density_altitude.htm equation 7b)
func vaporPressure(temperature goflying.Celsius, humidity goflying.RelativeHumidity) goflying.HPa {
	return goflying.HPa((float64(humidity) / 100) * float64(saturationVaporPressure(temperature)))
}
