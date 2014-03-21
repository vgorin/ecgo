// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

/*
 * Bitmatrix-based Cauchy Reed-Solomon decoding related routines
 */

package cauchy

/*
#include "cauchy.h"
#include "galois.h"
#include "jerasure.h"
#include "liberation.h"
#include "reed_sol.h"
*/
import "C"

import "fmt"
import "errors"
import "unsafe"

// TODO: create decoding-related structures and cache

// Decode recovers original data_block from chunks given
func (e *CauchyEncoder) Decode(chunks [][]byte, length int) (data_block []byte, err error) {
	b := int(e.p.b)
	n := int(e.p.n)

	if len(chunks) != n {
		return nil, errors.New(fmt.Sprintf("chunks length must be %d (variable b, r, n is not implemented yet)", n))
	}
	var chunk_length int = -1
	row_k_ones := C.int(1)
	missing := make([]int, n)
	var j int = 0
	for i := range chunks {
		if chunks[i] == nil || len(chunks[i]) == 0 {
			missing[j] = i
			j++
		} else if chunk_length == -1 {
			chunk_length = len(chunks[i])
		}
	}
	missing[j] = -1
	j++

	for i := range chunks {
		if chunks[i] == nil || len(chunks[i]) == 0 {
			chunks[i] = make([]byte, chunk_length)
		}
	}

	erasures := erasures(missing[:j]) // (*C.int)(unsafe.Pointer(&missing[:j][0]))
	pointers := make([]*byte, n)
	for i := range chunks {
		pointers[i] = &chunks[i][0]
	}
	data_ptrs := (**C.char)(unsafe.Pointer(&pointers[:b][0]))
	coding_ptrs := (**C.char)(unsafe.Pointer(&pointers[b:][0]))
	status := C.jerasure_bitmatrix_decode(e.k, e.m, e.w, e.bitmatrix, row_k_ones, erasures, data_ptrs, coding_ptrs, C.int(chunk_length), e.packet_size)

	if status != 0 {
		return nil, errors.New(fmt.Sprintf("jerasure_bitmatrix_decode returned %d status code", status))
	}

	data_block = make([]byte, 0, chunk_length*b)
	for i := 0; i < b; i++ {
		data_block = append(data_block, chunks[i]...)
	}

	return data_block[:length], nil
}

// CauchyDecode performs Decode on the default Encoder 
func CauchyDecode(chunks [][]byte, p encoder_params, length int) (data_block []byte, err error) {
	return GetCauchyEncoder(p).Decode(chunks, length)
}

