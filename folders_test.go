package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDirectory(t *testing.T) {
	dir := listFolder("tests")
	assert.Equal(t, 2, len(dir))
	assert.Contains(t, dir, "file.txt")
	assert.Contains(t, dir, "subfolder/subfile.txt")
}
