package bme280

import "testing"

func TestConfig_FilterCoefficient(t *testing.T) {
	tests := []struct {
		name string
		b    Config
		want FilterCoefficient
	}{
		{"filter coefficient off", 0x00, FilterCoefficientOff},
		{"filter coefficient 2", 0x04, FilterCoefficient2},
		{"filter coefficient 4", 0x08, FilterCoefficient4},
		{"filter coefficient 8", 0x0C, FilterCoefficient8},
		{"filter coefficient 16", 0x10, FilterCoefficient16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.FilterCoefficient(); got != tt.want {
				t.Errorf("FilterCoefficient() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestConfig_InactiveDuration(t *testing.T) {
	tests := []struct {
		name string
		b    Config
		want InactiveDuration
	}{
		{"inactive duration 0.5ms", 0x00, InactiveDuration0_5ms},
		{"inactive duration 62.5ms", 0x20, InactiveDuration62_5ms},
		{"inactive duration 125ms", 0x40, InactiveDuration125ms},
		{"inactive duration 250ms", 0x60, InactiveDuration250ms},
		{"inactive duration 500ms", 0x80, InactiveDuration500ms},
		{"inactive duration 1000ms", 0xA0, InactiveDuration1000ms},
		{"inactive duration 10ms", 0xC0, InactiveDuration10ms},
		{"inactive duration 20ms", 0xE0, InactiveDuration20ms},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.InactiveDuration(); got != tt.want {
				t.Errorf("InactiveDuration() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestConfig_SetFilterCoefficient(t *testing.T) {
	tests := []struct {
		name string
		b    Config
		new  FilterCoefficient
		want Config
	}{
		{"set filter coefficient to off", 0xFF, FilterCoefficientOff, 0xE3},
		{"set filter coefficient to 2", 0xFF, FilterCoefficient2, 0xE7},
		{"set filter coefficient to 4", 0xFF, FilterCoefficient4, 0xEB},
		{"set filter coefficient to 8", 0xFF, FilterCoefficient8, 0xEF},
		{"set filter coefficient to 16", 0xFF, FilterCoefficient16, 0xF3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SetFilterCoefficient(tt.new); got != tt.want {
				t.Errorf("SetFilterCoefficient() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestConfig_SetInactiveDuration(t *testing.T) {
	tests := []struct {
		name string
		b    Config
		new  InactiveDuration
		want Config
	}{
		{"set inactive duration to 0.5ms", 0xFF, InactiveDuration0_5ms, 0x1F},
		{"set inactive duration to 62.5ms", 0xFF, InactiveDuration62_5ms, 0x3F},
		{"set inactive duration to 125ms", 0xFF, InactiveDuration125ms, 0x5F},
		{"set inactive duration to 250ms", 0xFF, InactiveDuration250ms, 0x7F},
		{"set inactive duration to 500ms", 0xFF, InactiveDuration500ms, 0x9F},
		{"set inactive duration to 1000ms", 0xFF, InactiveDuration1000ms, 0xBF},
		{"set inactive duration to 10ms", 0xFF, InactiveDuration10ms, 0xDF},
		{"set inactive duration to 20ms", 0x1F, InactiveDuration20ms, 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SetInactiveDuration(tt.new); got != tt.want {
				t.Errorf("SetInactiveDuration() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlMeas_PressureOversampling(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlMeas
		want PressureOversampling
	}{
		{"skip pressure oversampling", 0x00, PressureOversamplingSkipped},
		{"pressure oversampling x1", 0x04, PressureOversampling1x},
		{"pressure oversampling x2", 0x08, PressureOversampling2x},
		{"pressure oversampling x4", 0x0C, PressureOversampling4x},
		{"pressure oversampling x8", 0x30, PressureOversampling8x},
		{"pressure oversampling x16", 0x34, PressureOversampling16x},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.PressureOversampling(); got != tt.want {
				t.Errorf("PressureOversampling() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlMeas_Mode(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlMeas
		want Mode
	}{
		{"sleep mode", 0x00, ModeSleep},
		{"forced mode", 0x01, ModeForced},
		{"normal mode", 0x03, ModeNormal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Mode(); got != tt.want {
				t.Errorf("Mode() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlMeas_SetPressureOversampling(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlMeas
		new  PressureOversampling
		want CtrlMeas
	}{
		{"set pressure oversampling to skipped", 0xFF, PressureOversamplingSkipped, 0xE3},
		{"set pressure oversampling to 1x", 0xFF, PressureOversampling1x, 0xE7},
		{"set pressure oversampling to 2x", 0xFF, PressureOversampling2x, 0xEB},
		{"set pressure oversampling to 4x", 0xFF, PressureOversampling4x, 0xEF},
		{"set pressure oversampling to 8x", 0xFF, PressureOversampling8x, 0xF3},
		{"set pressure oversampling to 16x", 0xFF, PressureOversampling16x, 0xF7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SetPressureOversampling(tt.new); got != tt.want {
				t.Errorf("SetPressureOversampling() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlMeas_SetMode(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlMeas
		new  Mode
		want CtrlMeas
	}{
		{"set to sleep mode", 0xFF, ModeSleep, 0xFC},
		{"set to forced mode", 0xFF, ModeForced, 0xFD},
		{"set to normal mode", 0xFC, ModeNormal, 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SetMode(tt.new); got != tt.want {
				t.Errorf("SetMode() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlMeas_SetTemperatureOversampling(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlMeas
		new  TemperatureOversampling
		want CtrlMeas
	}{
		{"set temperature oversampling to skipped", 0xFF, TemperatureOversamplingSkipped, 0x1F},
		{"set temperature oversampling to 1x", 0xFF, TemperatureOversampling1x, 0x3F},
		{"set temperature oversampling to 2x", 0xFF, TemperatureOversampling2x, 0x5F},
		{"set temperature oversampling to 4x", 0xFF, TemperatureOversampling4x, 0x7F},
		{"set temperature oversampling to 8x", 0xFF, TemperatureOversampling8x, 0x9F},
		{"set temperature oversampling to 16x", 0xFF, TemperatureOversampling16x, 0xBF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SetTemperatureOversampling(tt.new); got != tt.want {
				t.Errorf("SetTemperatureOversampling() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlMeas_TemperatureOversampling(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlMeas
		want TemperatureOversampling
	}{
		{"skip temperature oversampling", 0x00, TemperatureOversamplingSkipped},
		{"temperature oversampling x1", 0x20, TemperatureOversampling1x},
		{"temperature oversampling x2", 0x40, TemperatureOversampling2x},
		{"temperature oversampling x4", 0x60, TemperatureOversampling4x},
		{"temperature oversampling x8", 0x80, TemperatureOversampling8x},
		{"temperature oversampling x16", 0xA0, TemperatureOversampling16x},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.TemperatureOversampling(); got != tt.want {
				t.Errorf("TemperatureOversampling() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestFilterCoefficient_Value(t *testing.T) {
	tests := []struct {
		name string
		b    FilterCoefficient
		want int
	}{
		{"filter coefficient off", 0x00, 0},
		{"filter coefficient 2", 0x01, 2},
		{"filter coefficient 4", 0x02, 4},
		{"filter coefficient 8", 0x03, 8},
		{"filter coefficient 16", 0x04, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Value(); got != tt.want {
				t.Errorf("Value() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCtrlHum_HumidityOversampling(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlHum
		want HumidityOversampling
	}{
		{"skip humidity oversampling", 0x00, HumidityOversamplingSkipped},
		{"humidity oversampling x1", 0x01, HumidityOversampling1x},
		{"humidity oversampling x2", 0x02, HumidityOversampling2x},
		{"humidity oversampling x4", 0x03, HumidityOversampling4x},
		{"humidity oversampling x8", 0x04, HumidityOversampling8x},
		{"humidity oversampling x16", 0x05, HumidityOversampling16x},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HumidityOversampling(); got != tt.want {
				t.Errorf("HumidityOversampling() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestCtrlHum_SetHumidityOversampling(t *testing.T) {
	tests := []struct {
		name string
		b    CtrlHum
		new  HumidityOversampling
		want CtrlHum
	}{
		{"set humidity oversampling to skipped", 0xFF, HumidityOversamplingSkipped, 0xF8},
		{"set humidity oversampling to 1x", 0xFF, HumidityOversampling1x, 0xF9},
		{"set humidity oversampling to 2x", 0xFF, HumidityOversampling2x, 0xFA},
		{"set humidity oversampling to 4x", 0xFF, HumidityOversampling4x, 0xFB},
		{"set humidity oversampling to 8x", 0xFF, HumidityOversampling8x, 0xFC},
		{"set humidity oversampling to 16x", 0xFF, HumidityOversampling16x, 0xFD},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SetHumidityOversampling(tt.new); got != tt.want {
				t.Errorf("SetHumidityOversampling() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func TestHumidityOversampling_Value(t *testing.T) {
	tests := []struct {
		name string
		b    HumidityOversampling
		want int
	}{
		{"humidity oversampling skipped", HumidityOversamplingSkipped, 0},
		{"humidity oversampling 1x", HumidityOversampling1x, 1},
		{"humidity oversampling 2x", HumidityOversampling2x, 2},
		{"humidity oversampling 4x", HumidityOversampling4x, 4},
		{"humidity oversampling 8x", HumidityOversampling8x, 8},
		{"humidity oversampling 16x", HumidityOversampling16x, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Value(); got != tt.want {
				t.Errorf("Value() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestInactiveDuration_Milliseconds(t *testing.T) {
	tests := []struct {
		name string
		d    InactiveDuration
		want float64
	}{
		{"inactive duration 0.5ms", InactiveDuration0_5ms, 0.5},
		{"inactive duration 62.5ms", InactiveDuration62_5ms, 62.5},
		{"inactive duration 125ms", InactiveDuration125ms, 125},
		{"inactive duration 250ms", InactiveDuration250ms, 250},
		{"inactive duration 500ms", InactiveDuration500ms, 500},
		{"inactive duration 1000ms", InactiveDuration1000ms, 1000},
		{"inactive duration 10ms", InactiveDuration10ms, 10},
		{"inactive duration 20ms", InactiveDuration20ms, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Milliseconds(); got != tt.want {
				t.Errorf("Milliseconds() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestPressureOversampling_Value(t *testing.T) {
	tests := []struct {
		name string
		b    PressureOversampling
		want int
	}{
		{"pressure oversampling skipped", PressureOversamplingSkipped, 0},
		{"pressure oversampling 1x", PressureOversampling1x, 1},
		{"pressure oversampling 2x", PressureOversampling2x, 2},
		{"pressure oversampling 4x", PressureOversampling4x, 4},
		{"pressure oversampling 8x", PressureOversampling8x, 8},
		{"pressure oversampling 16x", PressureOversampling16x, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Value(); got != tt.want {
				t.Errorf("Value() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTemperatureOversampling_Value(t *testing.T) {
	tests := []struct {
		name string
		b    TemperatureOversampling
		want int
	}{
		{"temperature oversampling skipped", TemperatureOversamplingSkipped, 0},
		{"temperature oversampling 1x", TemperatureOversampling1x, 1},
		{"temperature oversampling 2x", TemperatureOversampling2x, 2},
		{"temperature oversampling 4x", TemperatureOversampling4x, 4},
		{"temperature oversampling 8x", TemperatureOversampling8x, 8},
		{"temperature oversampling 16x", TemperatureOversampling16x, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Value(); got != tt.want {
				t.Errorf("Value() = %d, want %d", got, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	tests := []struct {
		name  string
		b     byte
		size  int
		shift int
		want  byte
	}{
		{"case 01", 0x01, 1, 0, 0x01},
		{"case 02", 0x02, 1, 1, 0x01},
		{"case 03", 0x04, 1, 2, 0x01},
		{"case 04", 0x08, 1, 3, 0x01},
		{"case 05", 0x10, 1, 4, 0x01},
		{"case 06", 0x20, 1, 5, 0x01},
		{"case 07", 0x40, 1, 6, 0x01},
		{"case 08", 0x80, 1, 7, 0x01},
		{"case 09", 0x03, 2, 0, 0x03},
		{"case 10", 0x06, 2, 1, 0x03},
		{"case 11", 0x0C, 2, 2, 0x03},
		{"case 12", 0x18, 2, 3, 0x03},
		{"case 13", 0x30, 2, 4, 0x03},
		{"case 14", 0x60, 2, 5, 0x03},
		{"case 15", 0xC0, 2, 6, 0x03},
		{"case 16", 0x07, 3, 0, 0x07},
		{"case 17", 0x0E, 3, 1, 0x07},
		{"case 18", 0x1C, 3, 2, 0x07},
		{"case 19", 0x38, 3, 3, 0x07},
		{"case 20", 0x70, 3, 4, 0x07},
		{"case 21", 0xE0, 3, 5, 0x07},
		{"case 22", 0x0F, 4, 0, 0x0F},
		{"case 23", 0x1E, 4, 1, 0x0F},
		{"case 24", 0x3C, 4, 2, 0x0F},
		{"case 25", 0x78, 4, 3, 0x0F},
		{"case 26", 0xF0, 4, 4, 0x0F},
		{"case 27", 0x1F, 5, 0, 0x1F},
		{"case 28", 0x3E, 5, 1, 0x1F},
		{"case 29", 0x7C, 5, 2, 0x1F},
		{"case 30", 0xF8, 5, 3, 0x1F},
		{"case 31", 0x3F, 6, 0, 0x3F},
		{"case 32", 0x7E, 6, 1, 0x3F},
		{"case 33", 0xFC, 6, 2, 0x3F},
		{"case 34", 0x7F, 7, 0, 0x7F},
		{"case 35", 0xFE, 7, 1, 0x7F},
		{"case 36", 0xFF, 8, 0, 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValue(tt.b, tt.size, tt.shift); got != tt.want {
				t.Errorf("getValue() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}

func Test_oversamplingValue(t *testing.T) {
	tests := []struct {
		name string
		b    byte
		want int
	}{
		{"0", 0x00, 0},
		{"1", 0x01, 1},
		{"2", 0x02, 2},
		{"4", 0x03, 4},
		{"8", 0x04, 8},
		{"16", 0x05, 16},
		{"32", 0x06, 32},
		{"64", 0x07, 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oversamplingValue(tt.b); got != tt.want {
				t.Errorf("oversamplingValue() = %d, want %d", got, tt.want)
			}
		})
	}
}

func Test_setValue(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		shift int
		value byte
		want  byte
	}{
		{"case 01", 1, 0, 0x01, 0x01},
		{"case 02", 1, 1, 0x01, 0x02},
		{"case 03", 1, 2, 0x01, 0x04},
		{"case 04", 1, 3, 0x01, 0x08},
		{"case 05", 1, 4, 0x01, 0x10},
		{"case 06", 1, 5, 0x01, 0x20},
		{"case 07", 1, 6, 0x01, 0x40},
		{"case 08", 1, 7, 0x01, 0x80},
		{"case 09", 2, 0, 0x03, 0x03},
		{"case 10", 2, 1, 0x03, 0x06},
		{"case 11", 2, 2, 0x03, 0x0C},
		{"case 12", 2, 3, 0x03, 0x18},
		{"case 13", 2, 4, 0x03, 0x30},
		{"case 14", 2, 5, 0x03, 0x60},
		{"case 15", 2, 6, 0x03, 0xC0},
		{"case 16", 3, 0, 0x07, 0x07},
		{"case 17", 3, 1, 0x07, 0x0E},
		{"case 18", 3, 2, 0x07, 0x1C},
		{"case 19", 3, 3, 0x07, 0x38},
		{"case 20", 3, 4, 0x07, 0x70},
		{"case 21", 3, 5, 0x07, 0xE0},
		{"case 22", 4, 0, 0x0F, 0x0F},
		{"case 23", 4, 1, 0x0F, 0x1E},
		{"case 24", 4, 2, 0x0F, 0x3C},
		{"case 25", 4, 3, 0x0F, 0x78},
		{"case 26", 4, 4, 0x0F, 0xF0},
		{"case 27", 5, 0, 0x1F, 0x1F},
		{"case 28", 5, 1, 0x1F, 0x3E},
		{"case 29", 5, 2, 0x1F, 0x7C},
		{"case 30", 5, 3, 0x1F, 0xF8},
		{"case 31", 6, 1, 0x3F, 0x7E},
		{"case 32", 6, 2, 0x3F, 0xFC},
		{"case 33", 6, 0, 0x3F, 0x3F},
		{"case 34", 7, 0, 0x7F, 0x7F},
		{"case 35", 7, 1, 0x7F, 0xFE},
		{"case 36", 8, 0, 0xFF, 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setValue(0x00, tt.size, tt.shift, tt.value); got != tt.want {
				t.Errorf("setValue() = 0x%02X, want 0x%02X", got, tt.want)
			}
		})
	}
}
