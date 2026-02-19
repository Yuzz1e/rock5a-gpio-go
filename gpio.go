package gpio

import "errors"

// SetPull sets the pull resistor mode for the given GPIO pin via MMIO.
// Bank 0-4, port 'A'-'D', pin 0-7. Uses /dev/mem and bypasses the OS driver.
// GPIO0 Bank B is automatically split: pins 0-3 use PMU1, pins 4-7 use PMU2.
// GPIO4 Port D uses VCCIO2_IOC. Requires root or CAP_SYS_RAWIO.
func SetPull(bank int, port rune, pin int, mode PullMode) error {
	if pin < 0 || pin > 7 {
		return errors.New("pin must be 0-7")
	}
	baseAddr, offset, ok := GetPullReg(bank, port, pin)
	if !ok {
		return errors.New("invalid bank/port combination")
	}
	data, mask := pullBitsForPin(pin, mode)
	return writePullReg(baseAddr, offset, data, mask)
}
