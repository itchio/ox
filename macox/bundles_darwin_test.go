package macox_test

import (
	"testing"

	"github.com/itchio/ox/macox"
	"github.com/stretchr/testify/assert"
)

func Test_GetLibraryPath(t *testing.T) {
	{
		s, err := macox.GetExecutablePath("/Applications/TextEdit.app")
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	}
	{
		s, err := macox.GetLibraryPath()
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	}
	{
		s, err := macox.GetApplicationSupportPath()
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	}
}
