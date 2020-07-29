package geos

/*
#include "geos.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Reads the WKT serialization and produces geometries
type wktDecoder struct {
	r *C.GEOSWKTReader
}

// Creates a new WKT decoder, can be nil if initialization in the C API fails
func newWktDecoder() *wktDecoder {
	r := cGEOSWKTReader_create()
	if r == nil {
		return nil
	}
	d := &wktDecoder{r}
	runtime.SetFinalizer(d, (*wktDecoder).destroy)
	return d
}

// decode decodes the WKT string and returns a geometry
func (d *wktDecoder) decode(wkt string) (*Geometry, error) {
	cstr := C.CString(wkt)
	defer C.free(unsafe.Pointer(cstr))
	g := cGEOSWKTReader_read(d.r, cstr)
	if g == nil {
		return nil, Error()
	}
	gg := geomFromPtr(g)
	runtime.KeepAlive(d)
	return gg, nil
}

func (d *wktDecoder) destroy() {
	// XXX: mutex
	cGEOSWKTReader_destroy(d.r)
	d.r = nil
}

type wktEncoder struct {
	w *C.GEOSWKTWriter
}

func newWktEncoder() *wktEncoder {
	w := cGEOSWKTWriter_create()
	if w == nil {
		return nil
	}
	e := &wktEncoder{w}
	runtime.SetFinalizer(e, (*wktEncoder).destroy)
	return e
}

// Encode returns a string that is the geometry encoded as WKT
func (e *wktEncoder) encode(g *Geometry) (string, error) {
	cstr := cGEOSWKTWriter_write(e.w, g.g)
	defer C.free(unsafe.Pointer(cstr))
	if cstr == nil {
		return "", Error()
	}
	ret := C.GoString(str)
	runtime.KeepAlive(e)
	return ret, nil
}

func (e *wktEncoder) destroy() {
	// XXX: mutex
	cGEOSWKTWriter_destroy(e.w)
	e.w = nil
}
