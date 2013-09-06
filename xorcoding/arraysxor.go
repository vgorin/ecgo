package xorcoding

// arraysXor is unsafe function, arrays lengths, xor length
// and length must be equal and multiple of int64_size
func arraysXor(arrays [][]byte, xor []byte, length int) {
	if len(arrays) == 2 {
		arrayXor(arrays[0], arrays[1], xor, length)
		return
	}
	var xor_long int64 = 0

	for i := 0; i < length; i += int64_size {
		xor_long = 0
		for j := range arrays {
			xor_long ^= getInt64(arrays[j], i)
		}
		setInt64(xor, i, xor_long)
	}
}

func arrayXor(array1, array2, xor []byte, length int) {
	for i := 0; i < length; i += int64_size {
		setInt64(xor, i, getInt64(array1, i) ^ getInt64(array2, i))
	}
}
