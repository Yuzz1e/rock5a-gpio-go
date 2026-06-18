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

// pullBitsForPin returns data (lower 16b), writeMask (upper 16b write-enable),
// and lowerMask (lower 16b bit mask for merge) for pin n (0-7).
// Rockchip: Data bits - PE at bit[2n], PS at bit[2n+1];
// write-enable - PE at bit[2n+16], PS at bit[2n+17].
func pullBitsForPin(pin int, mode PullMode) (data, writeMask, lowerMask uint32) {
	if pin < 0 || pin > 7 {
		return 0, 0, 0
	}
	peBit := uint(2 * pin)
	psBit := uint(2*pin + 1)
	writeMask = (1 << (peBit + 16)) | (1 << (psBit + 16))
	lowerMask = (1 << peBit) | (1 << psBit)

	switch mode {
	case PullUp:
		data = (1 << peBit) | (1 << psBit) // PE=1, PS=1
	case PullDown:
		data = 1 << peBit // PE=1, PS=0
	case Floating:
		data = 0 // PE=0, PS=don't care
	default:
		data = 0
	}
	return data, writeMask, lowerMask
}
