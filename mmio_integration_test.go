//go:build linux && arm64

package gpio

import "testing"

func TestSetPull_GPIO1B_Pin3_Integration(t *testing.T) {
	base, off, ok := GetPullReg(1, 'B', 3)
	if !ok {
		t.Skip("platform not supported")
	}
	before, err := readPullRegMMIO(base, off)
	if err != nil {
		t.Skip(err)
	}
	if err := SetPull(1, 'B', 3, PullUp); err != nil {
		t.Fatal(err)
	}
	after, err := readPullRegMMIO(base, off)
	if err != nil {
		t.Fatal(err)
	}
	if before == after {
		t.Fatalf("register unchanged: %#x", after)
	}
	pe := (after>>6)&1 == 1
	ps := (after>>7)&1 == 1
	if !pe || !ps {
		t.Fatalf("pin3 not pull-up: PE=%d PS=%d reg=%#x", pe, ps, after)
	}
}
