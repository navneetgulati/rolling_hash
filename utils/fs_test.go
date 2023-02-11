package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f := New(16)
	assert.Equal(t, f.chunkSize, 16)
}

func TestOpen(t *testing.T) {
	filePath := "../mock_data/testData.txt"
	f := New(16)
	reader, err := f.Open(filePath)
	assert.NotNil(t, reader)
	assert.Nil(t, err)
}
