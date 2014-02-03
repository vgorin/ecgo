package xorcoding

//import "os"
//import "code.google.com/p/vitess/go/relog"

// number of bytes required to store encoding header length
const ehead_len = 2

// number of bytes required to store chunk length (64bit)
const lhead_len = 8

// full header length: encoding header + length header
const fhead_len = ehead_len + lhead_len

// logger
//var log = relog.New(os.Stdout, "", relog.INFO)

// XorEncode takes a data_block as input and produces b + 1 chunks from it
// Original data block can then be recovered from any of b chunks
func XorEncode(data_block []byte, b byte) (chunks [][]byte) {
	// 1. General preparation (init aux variables, pad data_block)
	n, original_length, chunk_length, data_block := encInitValues(data_block, b)

	// 2. Prepare 6-byte headers
	headers := encCreateHeaders(b, n, original_length, chunk_length)

	// 3. Calculate coding block contents
	chunks = encCreateChunks(data_block, b, n, chunk_length)

	// 4. Join headers with chunks
	chunks = encAppendChunks(n, headers, chunks)

	return chunks
}

// XorDecode recovers original data_block from chunks given
func XorDecode(chunks [][]byte) (data_block []byte, err error) {
	// 1. Read header
	b, n, original_length, chunk_length := decReadHeaders(chunks)

	// 2. Prepare chunks and check integrity
	sorted := decSortChunks(n, chunks)
	nil_chunk, err := decCheckIntegrity(b, sorted)
	if err != nil {
		return nil, err
	}

	// 3. Perform data recovery (if required) and 
	decRecover(chunks, sorted, b, nil_chunk, chunk_length)

	// 4. Join chunks into original data block
	data_block = decJoinChunks(sorted, int(b), chunk_length)

	// 5. Bring data block length to original
	return data_block[:original_length], nil
}
