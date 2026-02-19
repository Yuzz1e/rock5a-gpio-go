// Package gpio provides RK3588 (ROCK 5A) GPIO pull control via MMIO and
// value read/write via sysfs. All comments in this package are in English.
package gpio

// PullMode represents the pull resistor configuration for a GPIO pin.
// Rockchip IOC: PE (Pull Enable) and PS (Pull Select) control each pin.
// PS=1 -> pull-up, PS=0 -> pull-down; PE=1 -> enabled, PE=0 -> floating.
type PullMode int

const (
	// PullUp enables internal pull-up resistor.
	PullUp PullMode = iota
	// PullDown enables internal pull-down resistor.
	PullDown
	// Floating disables pull resistor (PE=0).
	Floating
)

// pullBitsForPin returns the 32-bit data and mask values for modifying
// the pull configuration of pin n (0-7) in a single register.
// Rockchip: Data bits - PE at bit[2n], PS at bit[2n+1];
// Mask bits - PE at bit[2n+16], PS at bit[2n+17].
// Lower 16 bits = data, upper 16 bits = write-enable mask.
func pullBitsForPin(pin int, mode PullMode) (data, mask uint32) {
	if pin < 0 || pin > 7 {
		return 0, 0
	}
	peBit := uint(2 * pin)
	psBit := uint(2*pin + 1)
	peMaskBit := peBit + 16
	psMaskBit := psBit + 16

	// We need to set both PE and PS in data, and set both mask bits to 1
	// so the hardware updates those bits.
	mask = (1 << peMaskBit) | (1 << psMaskBit)

	switch mode {
	case PullUp:
		data = (1 << peBit) | (1 << psBit) // PE=1, PS=1
	case PullDown:
		data = (1 << peBit) | (0 << psBit) // PE=1, PS=0
	case Floating:
		data = (0 << peBit) | (0 << psBit) // PE=0, PS=don't care
	default:
		data = 0
	}
	return data, mask
}
