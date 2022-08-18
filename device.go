package main

// A Device represents an external device connected to a Uxn CPU
//
// It has 16 bytes of internal "IO" memory that can be read and written to with
// the `DeviceWrite` and `DeviceRead` set of methods
//
// The behavior of reading and writing data from the device is completely
// defined by the `ReadByte` and `WriteByte` functions
//
// Reference: https://wiki.xxiivv.com/site/varvara.html
type Device struct {
	// A link to the parent virtual machine, because devices can define addresses
	// in main memory to call when they are changed
	u *Uxn
	// A device has 16 IO ports (0x00-0x0f) that can be written to and read from
	Data [16]byte
	// ReadByte defines what happens when a byte is read from a device
	ReadByte func(d *Device, port byte) byte
	// WriteByte defines what happens when a byte is written to a device
	WriteByte func(d *Device, port byte)
}

// DeviceWrite8 writes a single byte to the device at a given port
// A `port` a byte where the first 4 bits are the device being accessed (Ex: 0x1)
// and the second 4 bits are the IO port being accessed in the device (Ex: 0x08)
//
// For example, writing to a device with port (0x18) means the "Write" port of
// the Console device
//
// For implementation purposes, since the device has already been accessed if this
// function has been called, the `port` variable only uses the last 4 bits
func (d *Device) DeviceWrite8(port, data byte) {
	d.Data[port&0x0f] = data
	d.WriteByte(d, port&0x0f)
}

// DeviceWrite16 writes a single short to the device at a given port
func (d *Device) DeviceWrite16(port byte, data uint16) {
	d.DeviceWrite8(port, byte(data>>8))
	d.DeviceWrite8(port+1, byte(data))
}

// DeviceRead8 reads a single byte from the device at a given port
func (d *Device) DeviceRead8(port byte) byte {
	return d.ReadByte(d, port&0x0f)
}

// DeviceRead16 reads a single short from the device at a given port
func (d *Device) DeviceRead16(port byte) uint16 {
	return uint16(d.DeviceRead8(port)<<8) + uint16(d.ReadByte(d, (port+1)&0x0f))
}
