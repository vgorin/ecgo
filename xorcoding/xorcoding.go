package xorcoding

import "unsafe"
import "log"
import "fmt"

// number of bytes required to store encoding header length
const ehead_len = 2

// number of bytes required to store chunk length (64bit)
const lhead_len = 8

// full header length: encoding header + length header
const fhead_len = ehead_len + lhead_len

// XorEncode takes a data_block as input and produces b + 1 chunks from it
// Original data block can then be recovered from any of b chunks
func XorEncode(data_block []byte, b byte) (chunks [][]byte) {
	// 1. General preparation (init aux variables, pad data_block)

	// total number of chunks
	n := b + 1

	// byte counter for loops over GF(2^8)
	var i byte

	// save original block length
	original_length := len(data_block)

	// data block length must to be multiple of int64_size * b,
	var multiplier int = lhead_len * int(b)

	// pad data block if necessary
	if r := original_length % multiplier; r != 0 {
		// s stands for sparse block of length r
		s := make([]byte, multiplier-r)
		data_block = append(data_block, s...)
	}

	// calculate chunk length
	chunk_length := len(data_block) / int(b)

	log.Printf("b: %d; n: %d; original_length: %d (%d bytes); multiplier: %d; padded length: %d; chunk_length: %d\n", b, n, original_length, unsafe.Sizeof(original_length), multiplier, len(data_block), chunk_length)

	// 2. Prepare 6-byte headers

	// create headers
	headers := make([][]byte, n)
	for i = 0; i < n; i++ {
		// initialize
		headers[i] = make([]byte, fhead_len)
		// write encoding info
		headers[i][0] = b
		headers[i][1] = i
	}
	log.Printf("headers:\t%v\n", headers)
	// write chunk length info
	for i = 0; i < n; i++ {
		// first b-1 chunks has equal length
		*(*uintptr)(unsafe.Pointer(&headers[i][2])) = *(*uintptr)(unsafe.Pointer(&original_length))
	}
	log.Printf("headers:\t%v\n", headers)

	// 3. Calculate coding block contents

	// virtually split data_block into b slices
	// the last slice is for coding chunk
	chunks = make([][]byte, n)
	for i = 0; i < b; i++ {
		chunks[i] = data_block[int(i)*chunk_length : int(i+1)*chunk_length]
	}
	// allocate space for coding block
	chunks[b] = make([]byte, chunk_length)
	log.Printf("chunks:\t%v\n", chunks)

	// write coding block
	ArraysXor(chunks[:b], chunks[b], chunk_length)
	log.Printf("chunks:\t%v\n", chunks)

	// 4. Join headers with chunks

	for i = 0; i < n; i++ {
		chunks[i] = append(headers[i], chunks[i]...)
	}
	log.Printf("result:\t%v\n", chunks)

	return chunks
}

// XorDecode recovers original data_block from chunks given
func XorDecode(chunks [][]byte) (data_block []byte) {
	// 1. Read header
	
	// base parts number, number of chunks required
	b := chunks[0][0]
	// total number of chunks
	n := b + 1
	// original data_block length
	original_length := *(*uintptr)(unsafe.Pointer(&chunks[0][2]))
	// one chunk legth, used for missing block recovery and joining chunks
	chunk_length := len(chunks[0]) - fhead_len

	log.Printf("b: %d; n: %d; original_length: %d (%d bytes); chunk_length: %d\n", b, n, original_length, unsafe.Sizeof(original_length), chunk_length)


	// 2. Prepare chunks (required) and check integrity (optional)

	// sort chunks
	sorted := make([][]byte, n)

	for i := range chunks {
		sorted[chunks[i][1]] = chunks[i][fhead_len:]
	}

	log.Printf("unsorted:\t%v", chunks)
	log.Printf("sorted:\t%v", sorted)

	// check integrity (optional)
	var counter byte = 0
	var nil_chunk byte = b
	for i := range sorted {
		if sorted[i] != nil {
			counter++
		} else {
			nil_chunk = byte(i)
		}
	}
	if counter < b {
		panic(fmt.Sprintf("not enough chunks to decode data block; required: %d; found: %d", b, counter))
	}

	log.Printf("counter: %d; nil_chunk: %d\n", counter, nil_chunk)


	// 3. Perform data recovery (if required) and join into original data block

	// recover missing data chunk if required
	var i byte
	if nil_chunk != b {
		log.Println("recovering of missing data chunk required")
		striped := make([][]byte, b)
		for i = 0; i < b; i++ {
			striped[i] = chunks[i][fhead_len:]
		}
		log.Printf("striped:\t%v", striped)
		sorted[nil_chunk] = make([]byte, chunk_length)
		ArraysXor(striped, sorted[nil_chunk], chunk_length)
		log.Printf("sorted:\t%v", sorted)
	}

	// join sorted chunks into original data block
	data_block = make([]byte, 0, int(b) * chunk_length)
	for i = 0; i < b; i++ {
		data_block = append(data_block, sorted[i]...)
	}
	
	log.Printf("data_block:\t%v", data_block)

	return data_block[:original_length]
}

