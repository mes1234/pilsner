package stream_test

import (
	"pilsner/internal/stream"
	"testing"
)

func TestHashingFunctions(t *testing.T) {

	hash := stream.NewHash(0)

	hash.Update(make([]byte, 4))

}
