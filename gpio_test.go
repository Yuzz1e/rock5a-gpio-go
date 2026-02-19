package gpio

import (
	"testing"
)

func TestPullBitsForPin(t *testing.T) {
	// Pin 0: PE=bit0, PS=bit1; mask PE=bit16, PS=bit17.
	data, mask := pullBitsForPin(0, PullUp)
	if data != 3 || mask != (1<<16|1<<17) {
		t.Errorf("pin 0 PullUp: data=%d mask=%d", data, mask)
	}
	data, mask = pullBitsForPin(0, PullDown)
	if data != 1 || mask != (1<<16|1<<17) {
		t.Errorf("pin 0 PullDown: data=%d mask=%d", data, mask)
	}
	data, mask = pullBitsForPin(0, Floating)
	if data != 0 || mask != (1<<16|1<<17) {
		t.Errorf("pin 0 Floating: data=%d mask=%d", data, mask)
	}

	// Pin 3: PE=bit6, PS=bit7; mask bit22, bit23.
	data, mask = pullBitsForPin(3, PullUp)
	if data != (1<<6|1<<7) || mask != (1<<22|1<<23) {
		t.Errorf("pin 3 PullUp: data=%d mask=%d", data, mask)
	}

	// Out of range
	data, mask = pullBitsForPin(-1, PullUp)
	if data != 0 || mask != 0 {
		t.Errorf("pin -1 should return 0,0")
	}
	data, mask = pullBitsForPin(8, PullUp)
	if data != 0 || mask != 0 {
		t.Errorf("pin 8 should return 0,0")
	}
}

func TestGetPullReg_GPIO0_B_Split(t *testing.T) {
	// GPIO0 B pins 0-3 -> PMU1 (0xFD5F0000), offset 0x0024
	for _, pin := range []int{0, 1, 2, 3} {
		base, off, ok := GetPullReg(0, 'B', pin)
		if !ok || base != basePMU1IOC || off != offGPIO0B {
			t.Errorf("GPIO0 B pin %d: base=%08x off=%04x ok=%v", pin, base, off, ok)
		}
	}
	// GPIO0 B pins 4-7 -> PMU2 (0xFD5F4000), offset 0x0028
	for _, pin := range []int{4, 5, 6, 7} {
		base, off, ok := GetPullReg(0, 'B', pin)
		if !ok || base != basePMU2IOC || off != offGPIO0B2 {
			t.Errorf("GPIO0 B pin %d: base=%08x off=%04x ok=%v", pin, base, off, ok)
		}
	}
}

func TestGetPullReg_GPIO4_D_VCCIO2(t *testing.T) {
	// GPIO4 Port D must use VCCIO2_IOC (0xFD5FB000), offset 0x014C
	for pin := 0; pin <= 7; pin++ {
		base, off, ok := GetPullReg(4, 'D', pin)
		if !ok || base != baseVCCIO2IOC || off != offGPIO4D {
			t.Errorf("GPIO4 D pin %d: base=%08x off=%04x ok=%v", pin, base, off, ok)
		}
	}
}

func TestGetPullReg_Sanity(t *testing.T) {
	// GPIO0 A -> PMU1, 0x0020
	base, off, ok := GetPullReg(0, 'A', 0)
	if !ok || base != basePMU1IOC || off != offGPIO0A {
		t.Errorf("GPIO0 A: base=%08x off=%04x", base, off)
	}
	// GPIO1 A -> VCCIO1_4, 0x0110
	base, off, ok = GetPullReg(1, 'A', 0)
	if !ok || base != baseVCCIO14IOC || off != offGPIO1A {
		t.Errorf("GPIO1 A: base=%08x off=%04x", base, off)
	}
	// GPIO4 A -> VCCIO6, 0x0140 (not VCCIO2)
	base, off, ok = GetPullReg(4, 'A', 0)
	if !ok || base != baseVCCIO6IOC || off != offGPIO4A {
		t.Errorf("GPIO4 A: base=%08x off=%04x", base, off)
	}
}

func TestGPIONum(t *testing.T) {
	// bank 0, port A(0), pin 0 -> 0
	if n := GPIONum(0, 0, 0); n != 0 {
		t.Errorf("GPIONum(0,0,0)=%d", n)
	}
	// bank 0, port B(1), pin 0 -> 8
	if n := GPIONum(0, 1, 0); n != 8 {
		t.Errorf("GPIONum(0,1,0)=%d", n)
	}
	// bank 4, port D(3), pin 7 -> 32*4+3*8+7 = 159
	if n := GPIONum(4, 3, 7); n != 159 {
		t.Errorf("GPIONum(4,3,7)=%d", n)
	}
}
