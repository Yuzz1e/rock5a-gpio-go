package gpio

// IOC base addresses for RK3588 (from TRM).
const (
	basePMU1IOC     = 0xFD5F0000
	basePMU2IOC     = 0xFD5F4000
	baseVCCIO14IOC  = 0xFD5F9000
	baseVCCIO35IOC  = 0xFD5FA000
	baseVCCIO2IOC   = 0xFD5FB000
	baseVCCIO6IOC   = 0xFD5FC000
	baseEMMCIOC     = 0xFD5FD000
)

// Pull register offsets per port (within each IOC).
const (
	offGPIO0A = 0x0020
	offGPIO0B = 0x0024 // pins 0-3 on PMU1
	offGPIO0B2 = 0x0028 // pins 4-7 on PMU2
	offGPIO0C = 0x002C
	offGPIO0D = 0x0030
	offGPIO1A = 0x0110
	offGPIO1B = 0x0114
	offGPIO1C = 0x0118
	offGPIO1D = 0x011C
	offGPIO2A = 0x0120
	offGPIO2B = 0x0124
	offGPIO2C = 0x0128
	offGPIO2D = 0x012C
	offGPIO3A = 0x0130
	offGPIO3B = 0x0134
	offGPIO3C = 0x0138
	offGPIO3D = 0x013C
	offGPIO4A = 0x0140
	offGPIO4B = 0x0144
	offGPIO4C = 0x0148
	offGPIO4D = 0x014C // VCCIO2_IOC, not VCCIO6
)

// GetPullReg returns the physical base address and register offset for the
// pull control register of the given (bank, port, pin).
// Bank 0-4, port 'A'-'D', pin 0-7. GPIO0 Bank B is split: pins 0-3 use PMU1,
// pins 4-7 use PMU2. GPIO4 Port D uses VCCIO2_IOC (0xFD5FB000).
func GetPullReg(bank int, port rune, pin int) (baseAddr uint32, offset uint32, ok bool) {
	if pin < 0 || pin > 7 || bank < 0 || bank > 4 {
		return 0, 0, false
	}
	switch port {
	case 'A', 'a':
		// fall through to port index 0
	case 'B', 'b':
	case 'C', 'c':
	case 'D', 'd':
	default:
		return 0, 0, false
	}

	portIdx := portToIndex(port)
	if portIdx < 0 {
		return 0, 0, false
	}

	switch bank {
	case 0:
		switch portIdx {
		case 0: // A
			return basePMU1IOC, offGPIO0A, true
		case 1: // B: pins 0-3 -> PMU1/0x0024, pins 4-7 -> PMU2/0x0028
			if pin <= 3 {
				return basePMU1IOC, offGPIO0B, true
			}
			return basePMU2IOC, offGPIO0B2, true
		case 2: // C
			return basePMU2IOC, offGPIO0C, true
		case 3: // D
			return basePMU2IOC, offGPIO0D, true
		}
	case 1:
		switch portIdx {
		case 0:
			return baseVCCIO14IOC, offGPIO1A, true
		case 1:
			return baseVCCIO14IOC, offGPIO1B, true
		case 2:
			return baseVCCIO14IOC, offGPIO1C, true
		case 3:
			return baseVCCIO14IOC, offGPIO1D, true
		}
	case 2:
		switch portIdx {
		case 0:
			return baseEMMCIOC, offGPIO2A, true
		case 1:
			return baseVCCIO35IOC, offGPIO2B, true
		case 2:
			return baseVCCIO35IOC, offGPIO2C, true
		case 3:
			return baseEMMCIOC, offGPIO2D, true
		}
	case 3:
		switch portIdx {
		case 0:
			return baseVCCIO35IOC, offGPIO3A, true
		case 1:
			return baseVCCIO35IOC, offGPIO3B, true
		case 2:
			return baseVCCIO35IOC, offGPIO3C, true
		case 3:
			return baseVCCIO35IOC, offGPIO3D, true
		}
	case 4:
		switch portIdx {
		case 0:
			return baseVCCIO6IOC, offGPIO4A, true
		case 1:
			return baseVCCIO6IOC, offGPIO4B, true
		case 2:
			return baseVCCIO6IOC, offGPIO4C, true
		case 3:
			// GPIO4 Port D uses VCCIO2_IOC, not VCCIO6.
			return baseVCCIO2IOC, offGPIO4D, true
		}
	}
	return 0, 0, false
}

func portToIndex(port rune) int {
	switch port {
	case 'A', 'a':
		return 0
	case 'B', 'b':
		return 1
	case 'C', 'c':
		return 2
	case 'D', 'd':
		return 3
	}
	return -1
}
