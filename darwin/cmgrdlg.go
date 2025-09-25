// cmgrdlg.go: Implements the CentralManagerDelegate interface.  CoreBluetooth
// communicates events asynchronously via callbacks.  This file implements a
// synchronous interface by translating these callbacks into channel
// operations.

package darwin

import (
	"github.com/JuulLabs-OSS/cbgo"
	"github.com/go-ble/ble"
)

func (d *Device) CentralManagerDidUpdateState(cmgr cbgo.CentralManager) {
	d.evl.stateChanged.RxSignal(struct{}{})
}

func (d *Device) DidDiscoverPeripheral(cmgr cbgo.CentralManager, prph cbgo.Peripheral,
	advFields cbgo.AdvFields, rssi int) {

	// The Scan operation is happening in another goroutine. If a scan is still in progress,
	// a channel receive operation on d.advCh will give us a valid channel on which
	// we can report this advertisement. If the Scan operation is over or no scan is active,
	// d.advCh may be nil or closed, and we can return early.
	d.connLock.Lock()
	advCh := d.advCh
	d.connLock.Unlock()

	if advCh == nil {
		return
	}

	// Receive the actual advertisement channel from advCh.
	// If no scan is in progress, return early.
	var ch chan<- ble.Advertisement
	select {
	case ch = <-advCh:
		// got the channel, continue
	default:
		// no scan is currently waiting
		return
	}

	if ch == nil {
		return
	}

	// Build the advertisement struct.
	a := &adv{
		localName: advFields.LocalName,
		rssi:      int(rssi),
		mfgData:   advFields.ManufacturerData,
	}
	if advFields.Connectable != nil {
		a.connectable = *advFields.Connectable
	}
	if advFields.TxPowerLevel != nil {
		a.powerLevel = *advFields.TxPowerLevel
	}
	for _, u := range advFields.ServiceUUIDs {
		a.svcUUIDs = append(a.svcUUIDs, ble.UUID(u))
	}
	for _, sd := range advFields.ServiceData {
		a.svcData = append(a.svcData, ble.ServiceData{
			UUID: ble.UUID(sd.UUID),
			Data: sd.Data,
		})
	}
	a.peerUUID = ble.UUID(prph.Identifier())

	// Send the advertisement into the scan channel safely.
	// If the channel is closed or not ready, the select default prevents panic.
	select {
	case ch <- a:
		// delivered
	default:
		// scan channel closed or not ready, drop safely
	}
}

func (d *Device) DidConnectPeripheral(cmgr cbgo.CentralManager, prph cbgo.Peripheral) {
	fail := func(err error) {
		d.evl.connected.RxSignal(&eventConnected{
			err: err,
		})
	}

	c, err := newCentralConn(d, prph)
	if err != nil {
		fail(err)
	}

	d.evl.connected.RxSignal(&eventConnected{
		conn: c,
	})
}

func (d *Device) DidDisconnectPeripheral(cmgr cbgo.CentralManager, prph cbgo.Peripheral, err error) {
	c := d.findConn(ble.NewAddr(prph.Identifier().String()))
	if c != nil {
		close(c.done)
	}
}
