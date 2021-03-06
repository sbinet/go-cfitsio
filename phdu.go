package cfitsio

// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PrimaryHDU is the primary HDU
type PrimaryHDU struct {
	ImageHDU
}

// Name returns the value of the 'EXTNAME' Card (or "PRIMARY" if none)
func (hdu *PrimaryHDU) Name() string {
	card := hdu.header.Get("EXTNAME")
	if card == nil {
		return "PRIMARY"
	}
	return card.Value.(string)
}

// Version returns the value of the 'EXTVER' Card (or 1 if none)
func (hdu *PrimaryHDU) Version() int {
	card := hdu.header.Get("EXTVER")
	if card == nil {
		return 1
	}
	rv := reflect.ValueOf(card.Value)
	return int(rv.Int())
}

// newPrimaryHDU returns a new PrimaryHDU attached to file f.
func newPrimaryHDU(f *File, hdr Header) (HDU, error) {
	var err error
	hdu := &PrimaryHDU{
		ImageHDU{
			f:      f,
			header: hdr,
		},
	}
	return hdu, err
}

// NewPrimaryHDU creates a new PrimaryHDU with Header hdr in File f.
// It returns an error if f already has a Primary HDU.
func NewPrimaryHDU(f *File, hdr Header) (HDU, error) {
	var err error

	naxes := len(hdr.axes)
	c_naxes := C.int(naxes)
	slice := (*reflect.SliceHeader)((unsafe.Pointer(&hdr.axes)))
	c_axes := (*C.long)(unsafe.Pointer(slice.Data))
	c_status := C.int(0)

	C.fits_create_img(f.c, C.int(hdr.bitpix), c_naxes, c_axes, &c_status)
	if c_status > 0 {
		return nil, to_err(c_status)
	}

	for icard := range hdr.slice {
		card := &hdr.slice[icard]
		c_name := C.CString(card.Name)
		defer C.free(unsafe.Pointer(c_name))
		c_type := C.int(0)
		c_status := C.int(0)
		c_comm := C.CString(card.Comment)
		defer C.free(unsafe.Pointer(c_comm))
		var c_ptr unsafe.Pointer

		switch v := card.Value.(type) {
		case bool:
			c_type = C.TLOGICAL
			c_value := C.char(0) // 'F'
			if v {
				c_value = 1 // 'T'
			}
			c_ptr = unsafe.Pointer(&c_value)

		case byte:
			c_type = C.TBYTE
			c_ptr = unsafe.Pointer(&v)

		case uint16:
			c_type = C.TUSHORT
			c_ptr = unsafe.Pointer(&v)

		case uint32:
			c_type = C.TUINT
			c_ptr = unsafe.Pointer(&v)

		case uint64:
			c_type = C.TULONG
			c_ptr = unsafe.Pointer(&v)

		case uint:
			c_type = C.TULONG
			c_value := C.ulong(v)
			c_ptr = unsafe.Pointer(&c_value)

		case int8:
			c_type = C.TSBYTE
			c_ptr = unsafe.Pointer(&v)

		case int16:
			c_type = C.TSHORT
			c_ptr = unsafe.Pointer(&v)

		case int32:
			c_type = C.TINT
			c_ptr = unsafe.Pointer(&v)

		case int64:
			c_type = C.TLONG
			c_ptr = unsafe.Pointer(&v)

		case int:
			c_type = C.TLONG
			c_value := C.long(v)
			c_ptr = unsafe.Pointer(&c_value)

		case float32:
			c_type = C.TFLOAT
			c_ptr = unsafe.Pointer(&v)

		case float64:
			c_type = C.TDOUBLE
			c_ptr = unsafe.Pointer(&v)

		case complex64:
			c_type = C.TCOMPLEX
			c_ptr = unsafe.Pointer(&v) // FIXME: assumes same memory layout than C

		case complex128:
			c_type = C.TDBLCOMPLEX
			c_ptr = unsafe.Pointer(&v) // FIXME: assumes same memory layout than C

		case string:
			c_type = C.TSTRING
			c_value := C.CString(v)
			defer C.free(unsafe.Pointer(c_value))
			c_ptr = unsafe.Pointer(c_value)

		default:
			panic(fmt.Errorf("cfitsio: invalid card type (%T)", v))
		}

		C.fits_update_key(f.c, c_type, c_name, c_ptr, c_comm, &c_status)

		if c_status > 0 {
			return nil, to_err(c_status)
		}
	}

	if len(f.hdus) > 0 {
		return nil, fmt.Errorf("cfitsio: File has already a Primary HDU")
	}

	hdu, err := f.readHDU(0)
	if err != nil {
		return nil, err
	}
	f.hdus = append(f.hdus, hdu)

	return hdu, err
}

// EOF
