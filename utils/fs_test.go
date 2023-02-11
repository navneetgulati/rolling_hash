package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f := New(16)
	assert.Equal(t, 16, f.chunkSize)
}

func TestOpen(t *testing.T) {
	filePath := "../mock_data/testData.txt"
	f := New(16)
	reader, err := f.Open(filePath)
	assert.NotNil(t, reader)
	assert.Nil(t, err)
}

func TestNotEnoughChunks(t *testing.T) {
	filePath := "../mock_data/testData.txt"
	f := New(1500)
	reader, err := f.Open(filePath)
	assert.Nil(t, reader)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("file is not of appropriate size of generate enough chunks"), err)
}

func TestErrorInOpeningFile(t *testing.T) {
	filePath := "fake_path"
	f := New(16)
	reader, err := f.Open(filePath)
	assert.NotNil(t, err)
	assert.Nil(t, reader)
}
