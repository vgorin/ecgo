// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

/*
 * Bitmatrix-based Cauchy Reed-Solomon encoding related routines
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

import "unsafe"

// encoder_params stores parameters required for Encoder initialization
type encoder_params struct {
	b,
	r,
	n,
	word_size int
	packet_size int
}

// NewEncoderParams constructs encoder_params structure, checking validity of parameters specified
func NewEncoderParams(b, r, word_size, packet_size int) encoder_params {
	if b < 1 {
		panic("b < 1")
	}
	if r < 0 {
		panic("r < 0")
	}
	if b+r > 255 {
		panic("b + r > 255")
	}
	if word_size < 1 {
		panic("word_size < 1")
	}
	if 1<<uint(word_size) < b+r {
		panic("1 << word_size < b + r")
	}
	if packet_size < 1 {
		panic("packet_size < 1")
	}
	if packet_size < int(word_size) {
		panic("packet_size < word_size")
	}
	return encoder_params{
		b:           b,
		r:           r,
		n:           b + r,
		word_size:   word_size,
		packet_size: packet_size,
	}
}

// CauchyEncoder reprsents a Bitmatrix-based Cauchy Reed-Solomon Encoder
type CauchyEncoder struct {
	p *encoder_params
	k,
	m,
	w,
	packet_size C.int
	bitmatrix *C.int
}

// NewCauchyEncoder creates new CauchyEncoder, initialized with parameters from p
func NewCauchyEncoder(p encoder_params) *CauchyEncoder {
	k := C.int(p.b)
	m := C.int(p.r)
	w := C.int(p.word_size)
	matrix := C.cauchy_good_general_coding_matrix(k, m, w)
	bitmatrix := C.jerasure_matrix_to_bitmatrix(k, m, w, matrix)

	return &CauchyEncoder{
		p:           &p,
		k:           k,
		m:           m,
		w:           w,
		packet_size: C.int(p.packet_size),
		bitmatrix:   bitmatrix,
	}
}

// Encode takes a data_block as input and produces base + parity chunks from it
// Original data block can then be recovered from any of base chunks
func (e *CauchyEncoder) Encode(data_block []byte) (chunks [][]byte, length int) {
	// save original block length
	original_length := len(data_block)

	// extract required parameters
	b := e.p.b
	n := e.p.n
	word_size := e.p.word_size
	packet_size := e.p.packet_size
	k := e.k
	m := e.m
	w := e.w
	bitmatrix := e.bitmatrix

	// data block length must to be multiple of base, word size, block size
	var multiplier int = b * word_size * packet_size

	// pad data block if necessary
	if r := original_length % multiplier; r != 0 {
		// sparse block of length r
		// TODO: implement PKCS7 padding
		s := make([]byte, multiplier-r)
		data_block = append(data_block, s...)
	}

	// calculate chunk length
	chunk_length := len(data_block) / int(b)

	// virtually split data_block into k slices
	// the last m slices are for coding chunks
	chunks = make([][]byte, n)
	// aux structure, used to pass chunks into codec
	pointers := make([]*byte, n)

	var i int
	for i = 0; i < b; i++ {
		chunks[i] = data_block[i*chunk_length : (i+1)*chunk_length]
		pointers[i] = &chunks[i][0]
	}
	// allocate space for coding blocks
	for i = b; i < n; i++ {
		chunks[i] = make([]byte, chunk_length)
		pointers[i] = &chunks[i][0]
	}

	// write coding c
	data_ptrs := (**C.char)(unsafe.Pointer(&pointers[:b][0]))
	coding_ptrs := (**C.char)(unsafe.Pointer(&pointers[b:][0]))
	C.jerasure_bitmatrix_encode(k, m, w, bitmatrix, data_ptrs, coding_ptrs, C.int(chunk_length), e.packet_size)
	return chunks, original_length
}

// GetCauchyEncoder retrieves an Encoder from DefaultCauchyEncoderCache
func GetCauchyEncoder(p encoder_params) *CauchyEncoder {
	return DefaultCauchyEncoderCache.Get(p)
}

// CauchyEncode performs Encode on the default Encoder
func CauchyEncode(data_block []byte, p encoder_params) (chunks [][]byte, length int) {
	return GetCauchyEncoder(p).Encode(data_block)
}
