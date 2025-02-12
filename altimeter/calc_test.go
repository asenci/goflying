package altimeter

import (
	"testing"

	"github.com/westphae/goflying"
)

func TestAltimeter(t *testing.T) {
	tests := []struct {
		name     string
		pressure goflying.HPa
		qnh      goflying.HPa
		want     goflying.Feet
	}{
		{"isa", 1013.25, 1013.25, 0.0},
		{"case 01", 1013.0, 1015.0, 54.51552254459593},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Altimeter(tt.pressure, tt.qnh); got != tt.want {
				t.Errorf("Altimeter() = %g, want %g", got, tt.want)
			}
		})
	}
}

func TestDensityAltitude(t *testing.T) {
	tests := []struct {
		name        string
		pressure    goflying.HPa
		temperature goflying.Celsius
		want        goflying.Feet
	}{
		{"isa", 1013.25, 15.0, 16.74264731603165},
		{"case 01", 1008.0, 23.0, 1125.9023484025736},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DensityAltitude(tt.pressure, tt.temperature); got != tt.want {
				t.Errorf("DensityAltitude() = %g, want %g", got, tt.want)
			}
		})
	}
}

func TestDensityAltitudeWet(t *testing.T) {
	tests := []struct {
		name        string
		pressure    goflying.HPa
		temperature goflying.Celsius
		humidity    goflying.RelativeHumidity
		want        goflying.Feet
	}{
		{"isa dry", 1013.25, 15.0, 0.0, -0.3308763250435717},
		{"isa wet", 1013.25, 15.0, 100.0, 217.6003123243137},
		{"case 01", 1008.0, 23.0, 55.0, 1305.746800081177},
		{"case 02", 1008.0, 23.0, 0.0, 1108.8140489322639},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DensityAltitudeWet(tt.pressure, tt.temperature, tt.humidity); got != tt.want {
				t.Errorf("DensityAltitudeWet() = %g, want %g", got, tt.want)
			}
		})
	}
}

func TestPressureAltitude(t *testing.T) {
	tests := []struct {
		name     string
		pressure goflying.HPa
		want     goflying.Feet
	}{
		{"sea level", 1013.25, 0.0},
		{"FL180", 500.0, 18288.816087059095},
		{"FL360", 225.0, 36210.89622094748},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PressureAltitude(tt.pressure); got != tt.want {
				t.Errorf("PressureAltitude() = %g, want %g", got, tt.want)
			}
		})
	}
}

func Test_saturationVaporPressure(t *testing.T) {
	tests := []struct {
		name        string
		temperature goflying.Celsius
		want        goflying.HPa
	}{
		{"zero", -273.15, 8.516012143884676e+57},
		{"freeze", 0.0, 6.1078},
		{"boil", 100.0, 1021.9383170856194},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saturationVaporPressure(tt.temperature); got != tt.want {
				t.Errorf("saturationVaporPressure() = %g, want %g", got, tt.want)
			}
		})
	}
}
