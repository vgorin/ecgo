// Copyright 2013-2014 Vasiliy Gorin.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Original Jerasure C/C++ code –
// Copyright 2007 James S. Plank
// See copyright notice inside *.c, *.h files

package cauchy

import "testing"
import "bytes"

func TestCauchyEncode(t *testing.T) {
	p := NewEncoderParams(10, 5, 8, 256)
	data_block := make([]byte, 10000)
	chunks, length := CauchyEncode(data_block, p)
	t.Logf("chunks length: %d;\nlength: %d\n", len(chunks), length)
	if length != len(data_block) {
		t.Fail()
	}
}

func TestCauchyDecode(t *testing.T) {
	p := NewEncoderParams(10, 5, 8, 256)
	data_block := []byte("In information theory, an erasure code is a forward error correction (FEC) code for the binary erasure channel, which transforms a message of k symbols into a longer message (code word) with n symbols such that the original message can be recovered from a subset of the n symbols. The fraction r = k/n is called the code rate, the fraction k’/k, where k’ denotes the number of symbols required for recovery, is called reception efficiency.")
	chunks, length := CauchyEncode(data_block, p)
	t.Logf("chunks length: %d;\nlength: %d\n", len(chunks), length)
	if length != len(data_block) {
		t.Fail()
	}
	chunks[0] = nil
	chunks[3] = nil
	chunks[5] = nil
	chunks[9] = nil
	chunks[13] = nil
	original_block, err := CauchyDecode(chunks, p, length)
	if err != nil {
		t.Error(err)
	}
	if i := bytes.Compare(data_block, original_block); i != 0 {
		t.Errorf("decoded and original blocks missmatch: %d", i)
	}
}
