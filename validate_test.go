package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidFileName(t *testing.T) {
	assert.True(t, isValidFileName("filename.txt"))
	assert.False(t, isValidFileName("<.txt"))
}

func TestDoesFileExist(t *testing.T) {
	assert.True(t, doesFileExist("lock.go", "."))
	assert.False(t, doesFileExist("<.txt", "."))
}
