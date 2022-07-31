package goflying

import (
	"fmt"
	"strings"

	"github.com/kidoman/embd"
)

// I2CBus wraps an I2CBus interface with debugging support
type I2CBus struct {
	I2CBus embd.I2CBus
}

//goland:noinspection GoStandardMethods
func (b *I2CBus) ReadByte(addr byte) (byte, error) {
	if !Debugging {
		return b.I2CBus.ReadByte(addr)
	}

	bytes, err := b.ReadBytes(addr, 1)
	if len(bytes) < 1 {
		return 0, err
	}

	return bytes[0], err
}

func (b *I2CBus) ReadBytes(addr byte, num int) ([]byte, error) {
	if !Debugging {
		return b.I2CBus.ReadBytes(addr, num)
	}

	Logger.Printf("i2c: reading %d bytes from address 0x%02X", num, addr)

	bytes, err := b.I2CBus.ReadBytes(addr, num)
	if len(bytes) < 1 {
		return bytes, err
	}

	Logger.Printf("i2c: received from 0x%02X: %s", addr, debugBytes(bytes))

	return bytes, err
}

//goland:noinspection GoStandardMethods
func (b *I2CBus) WriteByte(addr, value byte) error {
	if !Debugging {
		return b.I2CBus.WriteByte(addr, value)
	}

	return b.WriteBytes(addr, []byte{value})
}

func (b *I2CBus) WriteBytes(addr byte, value []byte) error {
	if !Debugging {
		return b.I2CBus.WriteBytes(addr, value)
	}

	Logger.Printf("i2c: writing to 0x%02X: %s", addr, debugBytes(value))

	return b.I2CBus.WriteBytes(addr, value)
}

func (b *I2CBus) ReadFromReg(addr, reg byte, value []byte) error {
	if !Debugging {
		return b.I2CBus.ReadFromReg(addr, reg, value)
	}

	Logger.Printf("i2c: reading %d bytes from address 0x%02X at 0x%02X", len(value), addr, reg)

	err := b.I2CBus.ReadFromReg(addr, reg, value)

	Logger.Printf("i2c: received from 0x%02X at 0x%02X: %s", addr, reg, debugBytes(value))

	return err
}

func (b *I2CBus) ReadByteFromReg(addr, reg byte) (byte, error) {
	buf := make([]byte, 1)
	if err := b.ReadFromReg(addr, reg, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (b *I2CBus) ReadWordFromReg(addr, reg byte) (uint16, error) {
	buf := make([]byte, 2)
	if err := b.ReadFromReg(addr, reg, buf); err != nil {
		return 0, err
	}
	return (uint16(buf[0]) << 8) | uint16(buf[1]), nil
}

func (b *I2CBus) WriteToReg(addr, reg byte, value []byte) error {
	if !Debugging {
		return b.I2CBus.WriteToReg(addr, reg, value)
	}

	Logger.Printf("i2c: writing to 0x%02X at 0x%02X: %s", addr, reg, debugBytes(value))

	return b.I2CBus.WriteToReg(addr, reg, value)
}

func (b *I2CBus) WriteByteToReg(addr, reg, value byte) error {
	if !Debugging {
		return b.I2CBus.WriteByteToReg(addr, reg, value)
	}

	Logger.Printf("i2c: writing to 0x%02X at 0x%02X: %s", addr, reg, debugBytes([]byte{value}))

	return b.I2CBus.WriteByteToReg(addr, reg, value)
}

func (b *I2CBus) WriteWordToReg(addr, reg byte, value uint16) error {
	if !Debugging {
		return b.I2CBus.WriteWordToReg(addr, reg, value)
	}

	Logger.Printf("i2c: writing to 0x%02X at 0x%02X: %s", addr, reg, debugBytes([]byte{byte(value >> 8), byte(value)}))

	return b.I2CBus.WriteWordToReg(addr, reg, value)
}

func (b *I2CBus) Close() error {
	return b.I2CBus.Close()
}

func debugBytes(bs []byte) string {
	bytes := make([]string, len(bs))

	for i, b := range bs {
		bytes[i] = fmt.Sprintf("0x%02X", b)
	}

	return fmt.Sprintf("[%s]", strings.Join(bytes, " "))
}
