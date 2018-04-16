package macox_test

import (
	"testing"

	"github.com/itchio/ox/macox"
	"github.com/stretchr/testify/assert"
)

func Test_GetExecutablePath(t *testing.T) {
	s, err := macox.GetExecutablePath("/Applications/TextEdit.app")
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func Test_GetLibraryPath(t *testing.T) {
	s, err := macox.GetLibraryPath()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func Test_GetApplicationSupportPath(t *testing.T) {
	s, err := macox.GetApplicationSupportPath()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}
