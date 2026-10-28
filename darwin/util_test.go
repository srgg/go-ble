package darwin

import (
	"testing"

	"github.com/go-ble/ble"
	"github.com/tinygo-org/cbgo"
)

// Each input is a CoreBluetooth UUID string as cbgo parses it.
func TestBleUUIDKeepsGoBleForm(t *testing.T) {
	cases := []struct {
		cb   string
		want ble.UUID
	}{
		{"FF30", ble.MustParse("ff30")},
		{"DAF51202", ble.MustParse("daf51202-0000-1000-8000-00805f9b34fb")},
		{"6E400001-B5A3-F393-E0A9-E50E24DCCA9E", ble.MustParse("6e400001-b5a3-f393-e0a9-e50e24dcca9e")},
	}
	for _, c := range cases {
		u, err := cbgo.ParseUUID(c.cb)
		if err != nil {
			t.Fatalf("cbgo.ParseUUID(%q): %v", c.cb, err)
		}
		if got := bleUUID(u); !got.Equal(c.want) {
			t.Errorf("bleUUID(%s) = %s, want %s", c.cb, got, c.want)
		}
	}
}
