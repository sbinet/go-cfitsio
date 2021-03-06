package cfitsio

// #include <complex.h>
// #include "go-cfitsio.h"
// #include "go-cfitsio-utils.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// Value is a value in a FITS table
type Value interface{}

// Column represents a column in a FITS table
type Column struct {
	Name    string  // column name, corresponding to ``TTYPE`` keyword
	Format  string  // column format, corresponding to ``TFORM`` keyword
	Unit    string  // column unit, corresponding to ``TUNIT`` keyword
	Null    string  // null value, corresponding to ``TNULL`` keyword
	Bscale  float64 // bscale value, corresponding to ``TSCAL`` keyword
	Bzero   float64 // bzero value, corresponding to ``TZERO`` keyword
	Display string  // display format, corresponding to ``TDISP`` keyword
	Dim     []int64 // column dimension corresponding to ``TDIM`` keyword
	Start   int64   // column starting position, corresponding to ``TBCOL`` keyword
	IsVLA   bool    // whether this is a variable length array
	Value   Value   // value at current row
}

// inferFormat infers the FITS format associated with a Column, according to its HDUType and Go type.
func (col *Column) inferFormat(htype HDUType) error {
	var err error
	if col.Format != "" {
		return nil
	}

	str := gotype2FITS(col.Value, htype)
	if str == "" {
		return fmt.Errorf("cfitsio: %v can not handle [%T]", htype, col.Value)
	}
	col.Format = str
	return err
}

// read reads the value at column number icol and row irow, into ptr.
// icol and irow are 0-based indices.
func (col *Column) read(f *File, icol int, irow int64, ptr interface{}) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index
	c_anynul := C.int(0)
	c_status := C.int(0)

	var value interface{}
	rv := reflect.ValueOf(ptr).Elem()
	rt := reflect.TypeOf(rv.Interface())

	switch rt.Kind() {
	case reflect.Bool:
		c_type = C.TLOGICAL
		c_value := C.char(0) // 'F'
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = c_value == 1

	case reflect.Uint8:
		c_type = C.TBYTE
		var c_value C.char
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = byte(c_value)

	case reflect.Uint16:
		c_type = C.TUSHORT
		var c_value C.ushort
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint16(c_value)

	case reflect.Uint32:
		c_type = C.TUINT
		var c_value C.uint
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint32(c_value)

	case reflect.Uint64:
		c_type = C.TULONG
		var c_value C.ulong
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint64(c_value)

	case reflect.Uint:
		c_type = C.TULONG
		var c_value C.ulong
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = uint(c_value)

	case reflect.Int8:
		c_type = C.TSBYTE
		var c_value C.char
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int8(c_value)

	case reflect.Int16:
		c_type = C.TSHORT
		var c_value C.short
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int16(c_value)

	case reflect.Int32:
		c_type = C.TINT
		var c_value C.int
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int32(c_value)

	case reflect.Int64:
		c_type = C.TLONG
		var c_value C.long
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int64(c_value)

	case reflect.Int:
		c_type = C.TLONG
		var c_value C.long
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = int(c_value)

	case reflect.Float32:
		c_type = C.TFLOAT
		var c_value C.float
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = float32(c_value)

	case reflect.Float64:
		c_type = C.TDOUBLE
		var c_value C.double
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = float64(c_value)

	case reflect.Complex64:
		c_type = C.TCOMPLEX
		var c_value C.complexfloat
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = complex(
			float32(C.crealf(c_value)),
			float32(C.cimagf(c_value)),
		)

	case reflect.Complex128:
		c_type = C.TDBLCOMPLEX
		var c_value C.complexdouble
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = complex(
			float64(C.creal(c_value)),
			float64(C.cimag(c_value)),
		)

	case reflect.String:
		c_type = C.TSTRING
		// FIXME: get correct maximum size from card
		c_value := C.CStringN(C.FLEN_FILENAME)
		defer C.free(unsafe.Pointer(c_value))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, c_ptr, &c_anynul, &c_status)
		value = C.GoString(c_value)

	case reflect.Array:
		c_len := C.LONGLONG(rt.Len())
		switch rt.Elem().Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			v := make([]bool, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint8:
			c_type = C.TBYTE
			v := make([]uint8, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint16:
			c_type = C.TUSHORT
			v := make([]uint16, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint32:
			c_type = C.TUINT
			v := make([]uint32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint64:
			c_type = C.TULONG
			v := make([]uint64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint:
			c_type = C.TULONG
			v := make([]uint, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int8:
			c_type = C.TSBYTE
			v := make([]int8, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int16:
			c_type = C.TSHORT
			v := make([]int16, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int32:
			c_type = C.TINT
			v := make([]int32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int64:
			c_type = C.TLONG
			v := make([]int64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int:
			c_type = C.TLONG
			v := make([]int, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Float32:
			c_type = C.TFLOAT
			v := make([]float32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Float64:
			c_type = C.TDOUBLE
			v := make([]float64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			v := make([]complex64, int(c_len), int(c_len)) // FIXME: assume same binary layout
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			v := make([]complex128, int(c_len), int(c_len)) // FIXME: assume same binary layout
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)
		default:
			panic(fmt.Errorf("invalid type [%T]", value))
		}

	case reflect.Slice:
		switch rt.Elem().Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]bool, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint8:
			c_type = C.TBYTE
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]uint8, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint16:
			c_type = C.TUSHORT
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]uint16, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint32:
			c_type = C.TUINT
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]uint32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint64:
			c_type = C.TULONG
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]uint64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Uint:
			c_type = C.TULONG
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]uint, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int8:
			c_type = C.TSBYTE
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]int8, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int16:
			c_type = C.TSHORT
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]int16, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int32:
			c_type = C.TINT
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]int32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int64:
			c_type = C.TLONG
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]int64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Int:
			c_type = C.TLONG
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]int, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Float32:
			c_type = C.TFLOAT
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]float32, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Float64:
			c_type = C.TDOUBLE
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]float64, int(c_len), int(c_len))
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]complex64, int(c_len), int(c_len)) // FIXME: assume same binary layout
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			c_len := C.long(0)
			c_off := C.long(0)
			C.fits_read_descript(f.c, c_icol, c_irow, &c_len, &c_off, &c_status)
			if c_status > 0 {
				err = to_err(c_status)
				return err
			}
			v := make([]complex128, int(c_len), int(c_len)) // FIXME: assume same binary layout
			value = v
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&v)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_read_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(c_len), nil, c_ptr, &c_anynul, &c_status)
		default:
			panic(fmt.Errorf("invalid type [%T]", value))
		}

	default:
		panic(fmt.Errorf("invalid type [%T]", value))
	}

	if c_status > 0 {
		err = to_err(c_status)
	}

	rv.Set(reflect.ValueOf(value))
	col.Value = value
	return err
}

// write writes the current value of this Column into file f at column icol and row irow.
// icol and irow are 0-based indices.
func (col *Column) write(f *File, icol int, irow int64, value interface{}) error {
	var err error

	c_type := C.int(0)
	c_icol := C.int(icol + 1)      // 0-based to 1-based index
	c_irow := C.LONGLONG(irow + 1) // 0-based to 1-based index
	c_status := C.int(0)

	rv := reflect.ValueOf(value)
	rt := reflect.TypeOf(value)

	switch rt.Kind() {
	case reflect.Bool:
		c_type = C.TLOGICAL
		c_value := C.char(0) // 'F'
		if value.(bool) {
			c_value = 1
		}
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint8:
		c_type = C.TBYTE
		c_value := C.char(value.(byte))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint16:
		c_type = C.TUSHORT
		c_value := C.ushort(value.(uint16))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint32:
		c_type = C.TUINT
		c_value := C.uint(value.(uint32))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint64:
		c_type = C.TULONG
		c_value := C.ulong(value.(uint64))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Uint:
		c_type = C.TULONG
		c_value := C.ulong(value.(uint))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int8:
		c_type = C.TSBYTE
		c_value := C.char(value.(int8))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int16:
		c_type = C.TSHORT
		c_value := C.short(value.(int16))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int32:
		c_type = C.TINT
		c_value := C.int(value.(int32))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int64:
		c_type = C.TLONG
		c_value := C.long(value.(int64))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Int:
		c_type = C.TLONG
		c_value := C.long(value.(int))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Float32:
		c_type = C.TFLOAT
		c_value := C.float(value.(float32))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Float64:
		c_type = C.TDOUBLE
		c_value := C.double(value.(float64))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Complex64:
		c_type = C.TCOMPLEX
		value := value.(complex64)
		c_ptr := unsafe.Pointer(&value) // FIXME: assumes same memory layout
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Complex128:
		c_type = C.TDBLCOMPLEX
		value := value.(complex128)
		c_ptr := unsafe.Pointer(&value) // FIXME: assumes same memory layout
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.String:
		c_type = C.TSTRING
		c_value := C.CString(value.(string))
		defer C.free(unsafe.Pointer(c_value))
		c_ptr := unsafe.Pointer(&c_value)
		C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, 1, c_ptr, &c_status)

	case reflect.Slice:
		switch rt.Elem().Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			value := value.([]bool)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint8:
			c_type = C.TBYTE
			value := value.([]uint8)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint16:
			c_type = C.TUSHORT
			value := value.([]uint16)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint32:
			c_type = C.TUINT
			value := value.([]uint32)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint64:
			c_type = C.TULONG
			value := value.([]uint64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Uint:
			c_type = C.TULONG
			value := value.([]uint)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int8:
			c_type = C.TSBYTE
			value := value.([]int8)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int16:
			c_type = C.TSHORT
			value := value.([]int16)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int32:
			c_type = C.TINT
			value := value.([]int32)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int64:
			c_type = C.TLONG
			value := value.([]int64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Int:
			c_type = C.TLONG
			value := value.([]int)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Float32:
			c_type = C.TFLOAT
			value := value.([]float32)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Float64:
			c_type = C.TDOUBLE
			value := value.([]float64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			value := value.([]complex64)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value))) // FIXME: assume same bin-layout
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			value := value.([]complex128)
			slice := (*reflect.SliceHeader)((unsafe.Pointer(&value))) // FIXME: assume same bin-layout
			c_ptr := unsafe.Pointer(slice.Data)
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, C.LONGLONG(slice.Len), c_ptr, &c_status)
		default:
			panic(fmt.Errorf("unhandled type '%T'", value))
		}

	case reflect.Array:
		//rp := reflect.PtrTo(rv.Type())
		c_ptr := unsafe.Pointer(rv.Pointer())
		c_len := C.LONGLONG(rt.Len())

		switch rt.Elem().Kind() {
		case reflect.Bool:
			c_type = C.TLOGICAL
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint8:
			c_type = C.TBYTE
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint16:
			c_type = C.TUSHORT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint32:
			c_type = C.TUINT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint64:
			c_type = C.TULONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Uint:
			c_type = C.TULONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int8:
			c_type = C.TSBYTE
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int16:
			c_type = C.TSHORT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int32:
			c_type = C.TINT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int64:
			c_type = C.TLONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Int:
			c_type = C.TLONG
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Float32:
			c_type = C.TFLOAT
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Float64:
			c_type = C.TDOUBLE
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Complex64:
			c_type = C.TCOMPLEX
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)

		case reflect.Complex128:
			c_type = C.TDBLCOMPLEX
			C.fits_write_col(f.c, c_type, c_icol, c_irow, 1, c_len, c_ptr, &c_status)
		default:
			panic(fmt.Errorf("unhandled type '%T'", value))
		}

	default:
		panic(fmt.Errorf("unhandled type '%T' (%v)", value, rt.Kind()))
	}

	if c_status > 0 {
		err = to_err(c_status)
	}

	return err
}

// EOF
