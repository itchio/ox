package winox_test

import (
	"testing"

	"github.com/itchio/ox/winox"
	"github.com/stretchr/testify/assert"
)

func Test_GetFolderPath(t *testing.T) {
	type tcase struct {
		name string
		typ  winox.FolderType
	}

	cases := []tcase{
		tcase{name: "appData", typ: winox.FolderTypeAppData},
		tcase{name: "localAppData", typ: winox.FolderTypeLocalAppData},
		tcase{name: "profile", typ: winox.FolderTypeProfile},
		tcase{name: "startMenu", typ: winox.FolderTypeStartMenu},
		tcase{name: "programs", typ: winox.FolderTypePrograms},
	}

	for _, cas := range cases {
		t.Run(cas.name, func(t *testing.T) {
			s, err := winox.GetFolderPath(cas.typ)
			assert.NoError(t, err)
			assert.NotEmpty(t, s)
		})
	}

	{
		s, err := winox.GetFolderPath(winox.FolderType(-1))
		assert.Error(t, err)
		assert.Empty(t, s)
	}
}
