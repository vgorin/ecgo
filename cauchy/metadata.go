package cauchy

import ()

// Metadata stores all the required information to decode object
type Metadata struct {
	Base,
	WordSize,
	BlockSize,
	Length int
}

func (p *encoder_params) fill_metadata() Metadata {
	return Metadata{
		Base:      p.b,
		WordSize:  p.word_size,
		BlockSize: p.block_size,
	}
}
