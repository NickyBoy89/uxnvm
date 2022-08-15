package main

func NilDei(d *Device, port byte) byte {
	return d.Data[port]
}

type Device struct {
	u           *Uxn
	Data        [16]byte
	DeviceInput func(d *Device, port byte) byte
	DeviceOut   func(d *Device, port byte)
}

func (d *Device) DeviceWrite8(x, y byte) {
	d.Data[x&0xf] = y
	d.DeviceOut(d, x&0x0f)
}

func (d *Device) DeviceWrite16(x byte, y uint16) {
	d.DeviceWrite8(x, byte(y>>8))
	d.DeviceWrite8(x+1, byte(y))
}
