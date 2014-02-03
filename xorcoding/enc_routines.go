package xorcoding

func CalcPaddedCapacity(original_length int, b byte) (padded_capacity int) {
	// data block length must to be multiple of int64_size * b,
	var multiplier int = int64_size * int(b)

	// pad data block if necessary
	if r := original_length % multiplier; r != 0 {
		return original_length + multiplier - r
	}

	return original_length
}

// encInit pads original data block (returning poiter to padded one), inits n, original_length, chunk_length
func encInitValues(data_block []byte, b byte) (n byte, original_length, chunk_length int, padded_data_block []byte) {
	// total number of chunks
	n = b + 1

	// save original block length
	original_length = len(data_block)

	// data block length must to be multiple of int64_size * b,
	var multiplier int = int64_size * int(b)

	// pad data block if necessary
	padded_data_block = data_block
	if r := original_length % multiplier; r != 0 {
		// s stands for sparse block of length r
		s := make([]byte, multiplier-r)
		padded_data_block = append(data_block, s...)
	}

	// calculate chunk length
	chunk_length = len(padded_data_block) / int(b)

//	log.Debug("b: %d; n: %d; original_length: %d; multiplier: %d; padded length: %d; chunk_length: %d\n", b, n, original_length, multiplier, len(padded_data_block), chunk_length)

	return n, original_length, chunk_length, padded_data_block
}

// encHeaders creates headers for chunks (they contain b, chunk number, original block length)
func encCreateHeaders(b, n byte, original_length, chunk_length int) (headers [][]byte) {
	// create headers
	headers = make([][]byte, n)
	var i byte
	for i = 0; i < n; i++ {
		// initialize
		headers[i] = make([]byte, fhead_len, fhead_len + chunk_length)
		// write encoding info
		headers[i][0] = b
		headers[i][1] = i
	}
//	log.Debug("headers:\t%v\n", headers)
	// write chunk length info
	for i = 0; i < n; i++ {
		// first b-1 chunks has equal length
		setInt(headers[i], 2, original_length)
	}
//	log.Debug("headers:\t%v\n", headers)

	return headers
}

// encChunks creates stripped chunks from data block
func encCreateChunks(data_block []byte, b, n byte, chunk_length int) (chunks [][]byte) {
	// virtually split data_block into b slices
	// the last slice is for coding chunk
	chunks = make([][]byte, n)

	var i byte
	for i = 0; i < b; i++ {
		chunks[i] = data_block[int(i)*chunk_length : int(i+1)*chunk_length]
	}
	// allocate space for coding block
	chunks[b] = make([]byte, chunk_length)
//	log.Debug("chunks:\t%v\n", chunks)

	// write coding block
	arraysXor(chunks[:b], chunks[b], chunk_length)
//	log.Debug("chunks:\t%v\n", chunks)

	return chunks
}

func encAppendChunks(n byte, headers, chunks [][]byte) (full_chunks [][]byte) {
	var i byte
	for i = 0; i < n; i++ {
		chunks[i] = append(headers[i], chunks[i]...)
	}
//	log.Debug("result:\t%v\n", chunks)

	return chunks
}
