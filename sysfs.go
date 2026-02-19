package gpio

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	sysfsGPIOBase = "/sys/class/gpio"
)

// GPIO represents an exported GPIO pin for value read/write via sysfs.
type GPIO struct {
	num    int
	valuePath  string
	dirPath   string
	exported  bool
}

// GPIONum returns the Linux GPIO number for (bank, port, pin).
// Port: A=0, B=1, C=2, D=3. Formula: bank*32 + port*8 + pin.
func GPIONum(bank, port, pin int) int {
	return bank*32 + port*8 + pin
}

// OpenGPIO exports the GPIO and returns a handle for reading/writing value.
// Port is 0 for A, 1 for B, 2 for C, 3 for D. If the GPIO is already exported,
// OpenGPIO still returns a handle using the existing export.
func OpenGPIO(bank, port, pin int) (*GPIO, error) {
	if bank < 0 || bank > 4 || port < 0 || port > 3 || pin < 0 || pin > 7 {
		return nil, os.ErrInvalid
	}
	num := GPIONum(bank, port, pin)
	valuePath := filepath.Join(sysfsGPIOBase, "gpio"+strconv.Itoa(num), "value")
	dirPath := filepath.Join(sysfsGPIOBase, "gpio"+strconv.Itoa(num), "direction")

	exportPath := filepath.Join(sysfsGPIOBase, "export")
	f, err := os.OpenFile(exportPath, os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	_, err = f.WriteString(strconv.Itoa(num))
	f.Close()
	// Ignore "device or resource busy" (already exported).
	if err != nil && !strings.Contains(err.Error(), "busy") && !os.IsExist(err) {
		return nil, err
	}

	return &GPIO{num: num, valuePath: valuePath, dirPath: dirPath, exported: true}, nil
}

// SetDirection sets the direction to "in" or "out".
func (g *GPIO) SetDirection(dir string) error {
	if dir != "in" && dir != "out" {
		dir = "out"
	}
	return os.WriteFile(g.dirPath, []byte(dir), 0)
}

// Read returns the value as "0" (LOW) or "1" (HIGH).
func (g *GPIO) Read() (string, error) {
	b, err := os.ReadFile(g.valuePath)
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(b))
	if s == "0" {
		return "0", nil
	}
	return "1", nil
}

// Write sets the output value. Accepts "0", "1", "low", "high" (case-insensitive).
func (g *GPIO) Write(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "0", "low":
		return os.WriteFile(g.valuePath, []byte("0"), 0)
	case "1", "high":
		return os.WriteFile(g.valuePath, []byte("1"), 0)
	default:
		return os.ErrInvalid
	}
}

// Close unexports the GPIO. Idempotent.
func (g *GPIO) Close() error {
	if !g.exported {
		return nil
	}
	exportPath := filepath.Join(sysfsGPIOBase, "unexport")
	f, err := os.OpenFile(exportPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	_, _ = f.WriteString(strconv.Itoa(g.num))
	f.Close()
	g.exported = false
	return nil
}
