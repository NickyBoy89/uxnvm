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

func (d *Device) DeviceRead8(x byte) byte {
	return d.DeviceInput(d, x&0x0f)
}

func (d *Device) DeviceRead16(x byte) uint16 {
	return uint16(d.DeviceRead8(x)<<8) + uint16(d.DeviceInput(d, (x+1)&0x0f))
}
