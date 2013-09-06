package xorcoding

import "unsafe"

// sizeof(int64)
const int64_size = 8

func setInt(bytes []byte, offset, value int) {
	*(*int)(unsafe.Pointer(&bytes[offset])) = value
}

func getInt(bytes []byte, offset int) int {
	return *(*int)(unsafe.Pointer(&bytes[offset]))
}

func setInt64(bytes []byte, offset int, value int64) {
	*(*int64)(unsafe.Pointer(&bytes[offset])) = value
}

func getInt64(bytes []byte, offset int) int64 {
	return *(*int64)(unsafe.Pointer(&bytes[offset]))
}
