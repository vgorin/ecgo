package jerasurego

import "testing"
import "bytes"

func TestCauchyEncode(t *testing.T) {
	data_block := make([]byte, 10000)
	chunks, meta := CauchyEncode(data_block)
	t.Logf("chunks length: %v;\nmeta: %v\n", len(chunks), meta)
	if meta.Length != int64(len(data_block)) {
		t.Fail()
	}
}

func TestCauchyDecode(t *testing.T) {
	data_block := []byte("In information theory, an erasure code is a forward error correction (FEC) code for the binary erasure channel, which transforms a message of k symbols into a longer message (code word) with n symbols such that the original message can be recovered from a subset of the n symbols. The fraction r = k/n is called the code rate, the fraction k’/k, where k’ denotes the number of symbols required for recovery, is called reception efficiency.")
	chunks, meta := CauchyEncode(data_block)
	t.Logf("chunks length: %v;\nmeta: %v\n", len(chunks), meta)
	if meta.Length != int64(len(data_block)) {
		t.Fail()
	}
	chunks[1] = nil
	chunks[3] = nil
	chunks[5] = nil
	chunks[9] = nil
	chunks[13] = nil
	original_block, err := CauchyDecode(chunks, meta)
	if err != nil {
		t.Error(err)
	}
	if i := bytes.Compare(data_block, original_block); i != 0 {
		t.Errorf("decoded and original blocks missmatch: %d", i)
	}
}
