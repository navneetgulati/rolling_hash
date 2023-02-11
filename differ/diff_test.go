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

	assert.Equal(t, true, deltas[1].deleted)
	assert.Equal(t, "g hashes file difference", string(deltas[2].updatedLiterals))
	assert.Equal(t, 17, deltas[2].startIndex)
	assert.Equal(t, 32, deltas[2].endIndex)
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

	assert.Equal(t, 64, deltas[4].startIndex)
	assert.Equal(t, 80, deltas[4].endIndex)
	assert.Equal(t, "changes", string(deltas[4].updatedLiterals))
	assert.Equal(t, true, deltas[5].deleted)
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

	assert.Equal(t, 80, deltas[5].startIndex)
	assert.Equal(t, 96, deltas[5].endIndex)
	assert.Equal(t, "and text. This is written in a way to detect addition to the text", string(deltas[5].updatedLiterals))
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

	assert.Equal(t, 0, len(deltas))
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

	assert.Equal(t, true, deltas[0].deleted)
	assert.Equal(t, 1, deltas[1].startIndex)
	assert.Equal(t, 16, deltas[1].endIndex)
	assert.Equal(t, "The a Rollin", string(deltas[1].updatedLiterals))
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

	assert.Equal(t, 0, deltas[0].startIndex)
	assert.Equal(t, 16, deltas[0].endIndex)
	assert.Equal(t, "This is a different text and it is different from all the chunks above", string(deltas[0].updatedLiterals))
	assert.Equal(t, true, deltas[1].deleted)
	assert.Equal(t, true, deltas[2].deleted)
	assert.Equal(t, true, deltas[3].deleted)
	assert.Equal(t, true, deltas[4].deleted)
	assert.Equal(t, true, deltas[5].deleted)
}

func TestPrettifyDelta(t *testing.T) {
	txt1 := "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2 := "The a Rolling hash file diff algorithm. It should check for changes in file and text"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas := differInstance.GenerateDelta(signatures, buffReader2)
	prettyDelta := PrettifyDelta(deltas)
	assert.Equal(t, true, prettyDelta[0].deleted)
	assert.Equal(t, 1, prettyDelta[1].startIndex)
	assert.Equal(t, 16, prettyDelta[1].endIndex)
	assert.Equal(t, "The a Rollin", prettyDelta[1].updatedLiterals)
}
