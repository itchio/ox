package winox_test

import (
	"testing"

	"github.com/itchio/ox/winox"
	"github.com/stretchr/testify/assert"
)

func Test_GetFolderPath(t *testing.T) {
	{
		s, err := winox.GetFolderPath(winox.FolderTypeProfile)
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	}

	{
		s, err := winox.GetFolderPath(winox.FolderType(-1))
		assert.Error(t, err)
		assert.Empty(t, s)
	}
}
