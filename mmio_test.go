package gpio

import "testing"

func TestMergePullReg_Pin3PullUp(t *testing.T) {
	const oldLower = 0xFF57 // e.g. pin3 pull-down (PE=1 PS=0)
	data, writeMask, lowerMask := pullBitsForPin(3, PullUp)
	got := mergePullReg(oldLower, data, writeMask, lowerMask)
	wantLower := uint32(0xFFD7) // pin3 becomes pull-up
	if got&0xFFFF != wantLower {
		t.Fatalf("lower: got %#x want %#x", got&0xFFFF, wantLower)
	}
	if got&0xFFFF0000 != writeMask {
		t.Fatalf("writeMask: got %#x want %#x", got&0xFFFF0000, writeMask)
	}
}
