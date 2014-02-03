package xorcoding

import "errors"

func decReadHeaders(chunks [][]byte) (b, n byte, original_length, chunk_length int) {
	// base parts number, number of chunks required
	b = chunks[0][0]
	// total number of chunks
	n = b + 1
	// original data_block length
	original_length = getInt(chunks[0], 2)
	// one chunk legth, used for missing block recovery and joining chunks
	chunk_length = len(chunks[0]) - fhead_len

//	log.Debug("b: %d; n: %d; original_length: %d; chunk_length: %d\n", b, n, original_length, chunk_length)

	return b, n, original_length, chunk_length
}

func decSortChunks(n byte, chunks [][]byte) (sorted [][]byte) {
	sorted = make([][]byte, n)

	for i := range chunks {
		sorted[chunks[i][1]] = chunks[i][fhead_len:]
	}

//	log.Debug("unsorted:\t%v", chunks)
//	log.Debug("sorted:\t%v", sorted)

	return sorted
}

func decCheckIntegrity(b byte, sorted [][]byte) (nil_chunk byte, err error) {
	var counter byte = 0
	nil_chunk = b
	for i := range sorted {
		if sorted[i] != nil {
			counter++
		} else {
			nil_chunk = byte(i)
		}
	}
	if counter < b {
//		log.Error("not enough chunks to decode data block; required: %d; found: %d", b, counter)
		err = errors.New("not enough chunks to decode data block")
	}

//	log.Debug("counter: %d; nil_chunk: %d\n", counter, nil_chunk)

	return
}

// decRecover recovers missing chunk in sorted chunk array
func decRecover(chunks, sorted [][]byte, b, nil_chunk byte, chunk_length int) {
	var i byte
	if nil_chunk != b {
//		log.Debug("recovering of missing data chunk required")
		stripped := make([][]byte, b)
		for i = 0; i < b; i++ {
			stripped[i] = chunks[i][fhead_len:]
		}
//		log.Debug("striped:\t%v", stripped)
		sorted[nil_chunk] = make([]byte, chunk_length)
		arraysXor(stripped, sorted[nil_chunk], chunk_length)
//		log.Debug("sorted:\t%v", sorted)
	}
}

func decJoinChunks(sorted [][]byte, b, chunk_length int) (data_block []byte) {
	data_block = make([]byte, 0, b*chunk_length)
	for i := 0; i < b; i++ {
		data_block = append(data_block, sorted[i]...)
	}

//	log.Debug("data_block (%d bytes):\t%v", len(data_block), data_block)

	return data_block
}
