package xorcoding

import "testing"
import "bytes"

func TestXorEncode(t *testing.T) {
	data_block := []byte{128, 1, 2, 3, 4, 5, 6, 7, 64, 1, 2, 3, 4}
	chunks := XorEncode(data_block, 2)
	expected := [][]byte{
		{2, 0, 13, 0, 0, 0, 0, 0, 0, 0, 128, 1, 2, 3, 4, 5, 6, 7,},
		{2, 1, 13, 0, 0, 0, 0, 0, 0, 0,  64, 1, 2, 3, 4, 0, 0, 0,},
		{2, 2, 13, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 5, 6, 7,},
	}
	if len(chunks) != len(expected) {
		t.Fatalf("chunks length expected is %d, but found %d", len(expected), len(chunks))
	}
	for i := range chunks {
		if bytes.Compare(expected[i], chunks[i]) != 0 {
			t.Fatalf("bytes mismatch in chunks[%d]: expected %v but found %v", i, expected[i], chunks[i])
		}
	}
}

func TestIntegrity(t *testing.T) {
	data_block := []byte("This is a test for XOR-based erasure coding module")
	chunks := XorEncode(data_block, 2)
	chunks = chunks[1:]
	recovered_block, err := XorDecode(chunks)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data_block, recovered_block) != 0 {
		t.Fatalf("recovered block doesn't much original data block")
	}
}

func BenchmarkXorEncode(b *testing.B) {
	// encode data block using b = 2 (renamed as bp)
	block_size := b.N
	bp := byte(2)
	data_block := make([]byte, block_size)
	XorEncode(data_block, bp)
}
