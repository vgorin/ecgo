package jerasurego

/*
#include "cauchy.h"
#include "galois.h"
#include "jerasure.h"
#include "liberation.h"
#include "reed_sol.h"
*/
import "C"

import "unsafe"

type encoder_params struct {
	b,
	r,
	n,
	word_size byte
	block_size uint
}

func NewEncoderParams(b, r, word_size byte, block_size uint) encoder_params {
	if b < 1 {
		panic("b < 1")
	}
	if b + r > 255 {
		panic("b + r > 255")
	}
	if word_size < 1 {
		panic("word_size < 1")
	}
	if (1 << word_size) < uint(b + r) {
		panic("(1 << word_size) < b + r")
	}
	if block_size < uint(word_size) {
		panic("block_size < word_size")
	}
	return encoder_params{
		b:          b,
		r:          r,
		n:          b + r,
		word_size:  word_size,
		block_size: block_size,
	}
}

type CauchyEncoder struct {
	p *encoder_params
	k,
	m,
	w,
	block_size C.int
	bitmatrix *C.int
}

func NewCauchyEncoder(p encoder_params) *CauchyEncoder {
	k := C.int(p.b)
	m := C.int(p.r)
	w := C.int(p.word_size)
	matrix := C.cauchy_good_general_coding_matrix(k, m, w)
	bitmatrix := C.jerasure_matrix_to_bitmatrix(k, m, w, matrix)

	return &CauchyEncoder{
		p:          &p,
		k:          k,
		m:          m,
		w:          w,
		block_size: C.int(p.block_size),
		bitmatrix:  bitmatrix,
	}
}

// Encode takes a data_block as input and produces base + parity chunks from it
// Original data block can then be recovered from any of base chunks
func (e *CauchyEncoder) Encode(data_block []byte) (chunks [][]byte, length int) {
	// save original block length
	original_length := len(data_block)

	b := int(e.p.b)
	n := int(e.p.n)
	word_size := int(e.p.word_size)
	block_size := int(e.p.block_size)

	// data block length must to be multiple of base, word size, block size
	var multiplier int = b * word_size * block_size

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
	C.jerasure_bitmatrix_encode(e.k, e.m, e.w, e.bitmatrix, data_ptrs, coding_ptrs, C.int(chunk_length), e.block_size)
	return chunks, original_length
}

func GetCauchyEncoder(p encoder_params) *CauchyEncoder {
	return DefaultCauchyEncoderCache.Get(p)
}

func CauchyEncode(data_block []byte, p encoder_params) (chunks [][]byte, length int) {
	return GetCauchyEncoder(p).Encode(data_block)
}
