package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testString = "This is a test"
const expectedHash = uint(611517686)

func TestHash(t *testing.T) {
	adler32 := NewAdler32(4)
	adler32.Write([]byte(testString))
	hash := adler32.Hash()
	assert.Equal(t, expectedHash, hash)
}

func TestRollIn(t *testing.T) {
	adler32 := NewAdler32(4)
	stringBytes := []byte(testString)
	var hash uint
	for len(stringBytes) > 0 {
		hash = adler32.RollIn(stringBytes[0])
		stringBytes = stringBytes[1:]
	}
	assert.Equal(t, expectedHash, hash)
}

func TestRollOut(t *testing.T) {
	adler32 := NewAdler32(4)
	updatedString := "a" + testString
	stringBytes := []byte(updatedString)
	var hash uint
	adler32.Write(stringBytes)
	adler32.Hash()
	hash, removed := adler32.RollOut()
	assert.Equal(t, hash, expectedHash)
	assert.Equal(t, removed, uint8([]byte("a")[0]))
}
