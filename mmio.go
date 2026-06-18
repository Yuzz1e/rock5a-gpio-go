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

// mergePullReg computes the 32-bit value written to an IOC pull register.
func mergePullReg(oldLower, data, writeMask, lowerMask uint32) uint32 {
	newLower := (oldLower & ^lowerMask) | (data & lowerMask)
	return writeMask | newLower
}

// writePullReg writes the 32-bit pull register at (baseAddr + offset).
// Rockchip: upper 16 bits = write-enable mask, lower 16 bits = data.
// writeMask is already in bits [31:16]; do NOT shift it again.
func writePullReg(baseAddr, offset uint32, data, writeMask, lowerMask uint32) error {
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

	oldLower := *(*uint32)(regPtr) & 0xFFFF
	*(*uint32)(regPtr) = mergePullReg(oldLower, data, writeMask, lowerMask)

	return nil
}

// readPullRegMMIO reads the 32-bit pull register at (baseAddr + offset).
func readPullRegMMIO(baseAddr, offset uint32) (uint32, error) {
	mmioMu.Lock()
	defer mmioMu.Unlock()

	if devMemFd < 0 {
		fd, err := syscall.Open(devMem, syscall.O_RDWR|syscall.O_SYNC, 0)
		if err != nil {
			return 0, err
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
			return 0, err
		}
		mmapCache[pageStart] = mapped
		slice = mapped
	}

	regOff := physAddr - pageStart
	if regOff+4 > pageSize {
		return 0, errOutOfRange
	}
	return *(*uint32)(unsafe.Pointer(&slice[regOff])), nil
}

// errOutOfRange is used when register offset falls outside the mapped page.
var errOutOfRange = errors.New("register offset out of mapped page range")
