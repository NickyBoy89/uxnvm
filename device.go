package main

func NilDei(d *Device, port byte) byte {
	return d.Data[port]
}

// A Device represents an external device connected to a Uxn CPU
type Device struct {
	// A link to the parent virtual machine, because devices can define addresses
	// in main memory to call when they are changed
	u *Uxn
	// A device has 16 IO ports (0x00-0x0f) that can be written to and read from
	Data [16]byte
	// DeviceInput defines what happens when a byte is read from a device
	DeviceInput func(d *Device, port byte) byte
	// DeviceOut defines what happens when a byte is written to a device
	DeviceOut func(d *Device, port byte)
}

// DeviceWrite8 writes a single byte to the device at a given port
// The functionality behind writing to a device is entirely defined by the
// implementation of the device in question.
// For example, writing to the system state port `0x0f` with any data byte will
// halt the machine, and with this function, this is called as
// `DeviceWrite8(0x0f, 0x01)`
func (d *Device) DeviceWrite8(port, data byte) {
	d.Data[port&0x0f] = data
	d.DeviceOut(d, port&0x0f)
}

// DeviceWrite16 writes a single short to the device at a given port
func (d *Device) DeviceWrite16(port byte, data uint16) {
	d.DeviceWrite8(port, byte(data>>8))
	d.DeviceWrite8(port+1, byte(data))
}

// DeviceRead8 reads a single byte from the device at a given port
func (d *Device) DeviceRead8(port byte) byte {
	return d.DeviceInput(d, port&0x0f)
}

// DeviceRead16 reads a single short from the device at a given port
func (d *Device) DeviceRead16(port byte) uint16 {
	return uint16(d.DeviceRead8(port)<<8) + uint16(d.DeviceInput(d, (port+1)&0x0f))
}
