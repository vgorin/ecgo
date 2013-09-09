package xorcoding

/*
void region_xor(
	char *r1,
	char *r2,
	char *r3,
	int nbytes)
{
	long *l1;
	long *l2;
	long *l3;
	long *ltop;
	char *ctop;

	ctop = r1 + nbytes;
	ltop = (long *) ctop;
	l1 = (long *) r1;
	l2 = (long *) r2;
	l3 = (long *) r3;

	while (l1 < ltop) {
		*l3 = ((*l1)  ^ (*l2));
		l1++;
		l2++;
		l3++;
	}
}
*/
import "C"

import "unsafe"

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

// arrayXor is unsafe function, array1, array2, xor lengths
// and length must be equal and multiple of int64_size
func arrayXor(array1, array2, xor []byte, length int) {
	C.region_xor(
		(*C.char)(unsafe.Pointer(&array1[0])),
		(*C.char)(unsafe.Pointer(&array2[0])),
		(*C.char)(unsafe.Pointer(&xor[0])),
		C.int(length))
}
