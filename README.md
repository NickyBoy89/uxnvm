# uxnvm

An implementation of the [Uxn Virtual Machine](https://wiki.xxiivv.com/site/uxn.html), written in Go

# Running a ROM

`uxnvm <rom.rom>`

# Building

1. `go build`
2. `./uxnvm`

# Supported Features

All instuctions in the [Uxn instruction set](https://wiki.xxiivv.com/site/uxntal_reference.html) are supported, but some [Varvara](https://wiki.xxiivv.com/site/varvara.html) devices are not supported yet. Currently implemented are:

* [x] System
* [x] Console
* [ ] Screen
* [ ] Audio
* [ ] MIDI
* [ ] Controller
* [ ] Mouse
* [ ] File
* [ ] Datetime
* [x] Empty
* [x] Reserved
* [x] Reserved
