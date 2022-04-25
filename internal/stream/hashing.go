package stream

import (
	"encoding/binary"
	"hash/fnv"
)

type Hasher interface {
	Update(newContent []byte)
}

type hash struct {
	value uint64
}

func NewHash(value uint64) *hash {
	return &hash{value: value}
}

func (h *hash) Update(newContent []byte) {
	hasher := fnv.New64a()

	buf := make([]byte, 4)

	binary.PutUvarint(buf, h.value)

	concat := append(buf, newContent...)

	hasher.Write(concat)

	h.value = hasher.Sum64()
}
