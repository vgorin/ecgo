// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

/*
 * Bit Architecture compatability patch
 *     by vgorin
 *
 * This code is required to support a situation when
 * C code and Go code is compiled for different bit architectures
 *
 * TODO: We must compile both C and Go code in the same bit architecture!
 */

package cauchy

/*
#include "sizeof.h"
*/
import "C"

import "unsafe"
import "fmt"

// sizes in bytes of different data types
// sizeofint refers to C.int
// sizeofint16, sizeofint32, sizeofint64 refer to golang int8, int16, int32, int64
var sizeofint, sizeofint8, sizeofint16, sizeofint32, sizeofint64 int

func init() {
	sizeofint = int(C.sizeofint())
	sizeofint8 = int(unsafe.Sizeof(*new(int8)))
	sizeofint16 = int(unsafe.Sizeof(*new(int16)))
	sizeofint32 = int(unsafe.Sizeof(*new(int32)))
	sizeofint64 = int(unsafe.Sizeof(*new(int64)))
}

// erasures converts []int to *C.int
func erasures(missing []int) *C.int {
	// if architecture is the same then just return as is
	if sizeofint == int(unsafe.Sizeof(missing[0])) {
		return (*C.int)(unsafe.Pointer(&missing[0]))
	} else {
		// if not – determine architecture and convert
		switch sizeofint {
		case sizeofint8:
			int8array := make([]int8, len(missing))
			for i, v := range missing {
				int8array[i] = int8(v)
			}
			return (*C.int)(unsafe.Pointer(&int8array[0]))
		case sizeofint16:
			int16array := make([]int16, len(missing))
			for i, v := range missing {
				int16array[i] = int16(v)
			}
			return (*C.int)(unsafe.Pointer(&int16array[0]))
		case sizeofint32:
			int32array := make([]int32, len(missing))
			for i, v := range missing {
				int32array[i] = int32(v)
			}
			return (*C.int)(unsafe.Pointer(&int32array[0]))
		case sizeofint64:
			int64array := make([]int64, len(missing))
			for i, v := range missing {
				int64array[i] = int64(v)
			}
			return (*C.int)(unsafe.Pointer(&int64array[0]))
		default:
			panic(fmt.Sprintf("not standart sizeof(int): %d", sizeofint))
		}
	}
}
