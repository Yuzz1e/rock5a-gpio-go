# rock5a-gpio-go

RK3588 (ROCK 5A) GPIO library: MMIO pull control and sysfs value read/write.

- **SetPull(bank, port, pin, mode)** — configures pull resistor via `/dev/mem` and `syscall.Mmap` (bypasses OS driver). Handles GPIO0 Bank B split (pins 0–3: PMU1, 4–7: PMU2) and GPIO4 Port D (VCCIO2_IOC).
- **OpenGPIO / Read / Write / Close** — export GPIO via sysfs and read/write HIGH/LOW on `value`.

**Requirements:** Root or `CAP_SYS_RAWIO` is required for `/dev/mem` (SetPull). MMIO accesses hardware directly; use with care.

## Example

```go
import "github.com/rock5a-gpio-go"

// Pull-up on GPIO0 B pin 2
_ = gpio.SetPull(0, 'B', 2, gpio.PullUp)

// Sysfs: export and drive high
pin, _ := gpio.OpenGPIO(0, 1, 2) // bank 0, port B(1), pin 2
_ = pin.SetDirection("out")
_ = pin.Write("high")
v, _ := pin.Read() // "1"
_ = pin.Close()
```
