package gpio

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

const (
	devMem  = "/dev/mem"
	pageSize = 4096
)

var (
	mmioMu     sync.Mutex
	devMemFd   int = -1
	mmapCache  = make(map[uint32][]byte) // page-aligned base -> mapped slice
)

// writePullReg writes the 32-bit pull register at (baseAddr + offset).
// Rockchip: upper 16 bits = write mask, lower 16 bits = data.
// We write (mask<<16)|data so only the masked bits are updated.
func writePullReg(baseAddr, offset uint32, data, mask uint32) error {
	mmioMu.Lock()
	defer mmioMu.Unlock()

	if devMemFd < 0 {
		fd, err := syscall.Open(devMem, syscall.O_RDWR|syscall.O_SYNC, 0)
		if err != nil {
			return err
		}
		devMemFd = fd
	}

	physAddr := baseAddr + offset
	pageStart := physAddr & ^uint32(pageSize-1)

	slice, ok := mmapCache[pageStart]
	if !ok {
		mapped, err := syscall.Mmap(devMemFd, int64(pageStart), pageSize,
			syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			return err
		}
		mmapCache[pageStart] = mapped
		slice = mapped
	}

	regOff := physAddr - pageStart
	if regOff+4 > pageSize {
		return errOutOfRange
	}
	regPtr := unsafe.Pointer(&slice[regOff])

	old := *(*uint32)(regPtr)
	oldLower := old & 0xFFFF
	// Rockchip: bits 31-16 are write-enable mask; only lower bits with mask 1 get updated.
	// newLower = (oldLower & ^mask) | (data & mask); we write (mask<<16) | newLower.
	newLower := (oldLower & ^mask) | (data & mask)
	*(*uint32)(regPtr) = (mask << 16) | newLower

	return nil
}

// errOutOfRange is used when register offset falls outside the mapped page.
var errOutOfRange = errors.New("register offset out of mapped page range")
