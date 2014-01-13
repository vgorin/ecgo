// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

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
import "runtime"
import "errors"
import "fmt"

// TODO: 1. read codec parameters from config
// TODO: 2. later on remove config and implement caching for matrices instead
const (
	Base      = 10
	Parity    = 5
	Total     = Base + Parity
	WordSize  = 8
	BlockSize = 256
)

// NumCPU - how many goroutines can be started for encoding/decoding
var NumCPU int

// base, parity & total (synonims for Base, Parity & Total)
var b, r, n int = Base, Parity, Total

// k stands for base, m stands for parity, w stands for word_size
var k, m, w, block_size C.int

// bitmatrix is a subject for caching
var bitmatrix *C.int

// TODO: implement PKCS5 padding instead of saving block length
type Metadata struct {
	// original object length
	// TODO: add b and r
	Length int64
}

func init() {
	NumCPU = runtime.NumCPU()
	k = C.int(Base)
	m = C.int(Parity)
	w = C.int(WordSize)
	block_size = C.int(BlockSize)
	matrix := C.cauchy_good_general_coding_matrix(k, m, w)
	bitmatrix = C.jerasure_matrix_to_bitmatrix(k, m, w, matrix)
}

// Wrappers

// CauchyEncode takes a data_block as input and produces k (base) + m (parity) chunks from it
// Original data block can then be recovered from any of k chunks
func CauchyEncode(data_block []byte) (chunks [][]byte, meta *Metadata) {
	// save original block length
	original_length := len(data_block)

	// data block length must to be multiple of k, word size,
	var multiplier int = b * WordSize * BlockSize

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
	C.jerasure_bitmatrix_encode(k, m, w, bitmatrix, data_ptrs, coding_ptrs, C.int(chunk_length), block_size)
	return chunks, &Metadata{Length: int64(original_length)}
}

// CauchyDecode recovers original data_block from chunks given
func CauchyDecode(chunks [][]byte, meta *Metadata) (data_block []byte, err error) {
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
	status := C.jerasure_bitmatrix_decode(k, m, w, bitmatrix, row_k_ones, erasures, data_ptrs, coding_ptrs, C.int(chunk_length), block_size)

	if status != 0 {
		return nil, errors.New(fmt.Sprintf("jerasure_bitmatrix_decode returned %d status code", status))
	}

	data_block = make([]byte, 0, chunk_length*b)
	for i := 0; i < b; i++ {
		data_block = append(data_block, chunks[i]...)
	}

	return data_block[:meta.Length], nil
}
