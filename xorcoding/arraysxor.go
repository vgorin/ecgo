package xorcoding

// arraysXor is unsafe function, data_blocks lengths, coding_block length
// and block_length must be equal and multiple of 8
func arraysXor(data_blocks [][]byte, coding_block []byte, block_length int) {
	var xor_long int64 = 0

	for i := 0; i < block_length; i += 8 {
		xor_long = 0
		for j := range data_blocks {
			xor_long ^= getInt64(data_blocks[j], i)
		}
		setInt64(coding_block, i, xor_long)
	}
}
