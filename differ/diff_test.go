package differ

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkChange(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "This is a Rolling hashes file difference algorithm. It should check for changes in file and text"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)

	assert.Equal(t, deltas[1].deleted, true)
	assert.Equal(t, string(deltas[2].updatedLiterals), "g hashes file difference")
	assert.Equal(t, deltas[2].startIndex, 17)
	assert.Equal(t, deltas[2].endIndex, 32)
}

func TestChunkDeletion(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "This is a Rolling hash file diff algorithm. It should check for changes"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)

	assert.Equal(t, deltas[4].startIndex, 64)
	assert.Equal(t, deltas[4].endIndex, 80)
	assert.Equal(t, string(deltas[4].updatedLiterals), "changes")
	assert.Equal(t, deltas[5].deleted, true)
}
func TestChunkAddition(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text. This is written in a way to detect addition to the text"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)

	assert.Equal(t, deltas[5].startIndex, 80)
	assert.Equal(t, deltas[5].endIndex, 96)
	assert.Equal(t, string(deltas[5].updatedLiterals), "and text. This is written in a way to detect addition to the text")
}

func TestTextWithNoChanges(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)

	assert.Equal(t, len(deltas), 0)
}

func TestWithFirstChunkChange(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "The a Rolling hash file diff algorithm. It should check for changes in file and text"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)

	assert.Equal(t, deltas[0].deleted, true)
	assert.Equal(t, deltas[1].startIndex, 1)
	assert.Equal(t, deltas[1].endIndex, 16)
	assert.Equal(t, string(deltas[1].updatedLiterals), "The a Rollin")
}

func TestAllChunksChange(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "This is a different text and it is different from all the chunks above"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)

	assert.Equal(t, deltas[0].startIndex, 0)
	assert.Equal(t, deltas[0].endIndex, 16)
	assert.Equal(t, string(deltas[0].updatedLiterals), "This is a different text and it is different from all the chunks above")
	assert.Equal(t, deltas[1].deleted, true)
	assert.Equal(t, deltas[2].deleted, true)
	assert.Equal(t, deltas[3].deleted, true)
	assert.Equal(t, deltas[4].deleted, true)
	assert.Equal(t, deltas[5].deleted, true)
}
