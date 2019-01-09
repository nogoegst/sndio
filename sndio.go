package sndio

/*
#cgo LDFLAGS: -lsndio
#include <sndio.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Device struct {
	hdl *C.struct_sio_hdl
}

const (
	Play      = C.SIO_PLAY
	Record    = C.SIO_REC
	AnyDevice = C.SIO_DEVANY
)

func Open(name string, mode int, noblock bool) (*Device, error) {
	Cnoblock := C.int(0)
	if noblock {
		Cnoblock = C.int(1)
	}
	hdl := C.sio_open(C.CString(name), C.uint(mode), Cnoblock)
	if hdl == nil {
		return nil, errors.New("sndio: unable to open unit")
	}
	device := &Device{
		hdl: hdl,
	}
	return device, nil
}

func (d *Device) Close() error {
	C.sio_close(d.hdl)
	return nil
}

type Parameters struct {
	PlayChans    int
	RecordChans  int
	Bits         int
	Signed       bool
	LittleEndian bool
}

func (d *Device) SetParameters(p *Parameters) (*Parameters, error) {
	par := &C.struct_sio_par{}
	C.sio_initpar(par)

	par.pchan = C.uint(p.PlayChans)
	par.rchan = C.uint(p.RecordChans)
	par.bits = C.uint(p.Bits)
	par.sig = C.uint(0)
	if p.Signed {
		par.sig = C.uint(1)
	}
	par.le = C.uint(0)
	if p.LittleEndian {
		par.le = C.uint(1)
	}

	e := C.sio_setpar(d.hdl, par)
	if e == 0 {
		return nil, errors.New("sio_setpar failed")
	}
	e = C.sio_getpar(d.hdl, par)
	if e == 0 {
		return nil, errors.New("sio_getpar failed")
	}

	// TODO: check parameters
	return nil, nil
}

func (d *Device) Start() error {
	e := C.sio_start(d.hdl)
	if e == 0 {
		return errors.New("start failed")
	}
	return nil
}

func (d *Device) Stop() error {
	e := C.sio_stop(d.hdl)
	if e == 0 {
		return errors.New("stop failed")
	}
	return nil
}

func (d *Device) Read(p []byte) (int, error) {
	n := C.sio_read(d.hdl, unsafe.Pointer(&p[0]), C.size_t(len(p)))
	return int(n), nil
}

func (d *Device) Write(p []byte) (int, error) {
	n := C.sio_write(d.hdl, unsafe.Pointer(&p[0]), C.size_t(len(p)))
	return int(n), nil
}
