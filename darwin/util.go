package darwin

import (
	"bytes"

	"github.com/go-ble/ble"
	"github.com/tinygo-org/cbgo"
)

// The Bluetooth Base UUID less its 32-bit value, in go-ble's byte order.
var baseUUIDTail = ble.MustParse("00000000-0000-1000-8000-00805F9B34FB")[:12]

// bleUUID converts a UUID cbgo reports into go-ble's form. cbgo spells a 16-bit UUID as the
// 128-bit Base UUID it stands for; go-ble compares UUIDs byte for byte and spells one in 2 bytes.
// A 32-bit UUID has no 2-byte form and stays 128-bit.
func bleUUID(u cbgo.UUID) ble.UUID {
	if len(u) == 16 && bytes.Equal(u[:12], baseUUIDTail) && u[14] == 0 && u[15] == 0 {
		return ble.UUID{u[12], u[13]}
	}
	return ble.UUID(u)
}

func uuidSlice(uu []ble.UUID) [][]byte {
	us := [][]byte{}
	for _, u := range uu {
		us = append(us, ble.Reverse(u))
	}
	return us
}

func uuidStrWithDashes(s string) string {
	if len(s) != 32 {
		return s
	}

	// 01234567-89ab-cdef-0123-456789abcdef
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

func uuidsToCbgoUUIDs(uuids []ble.UUID) []cbgo.UUID {
	var cbuuids []cbgo.UUID

	for _, u := range uuids {
		cbuuids = append(cbuuids, cbgo.UUID(u))
	}

	return cbuuids
}
