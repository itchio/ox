package ox_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itchio/ox"
	"github.com/stretchr/testify/assert"
)

func doTestPreallocate(t *testing.T) {
	assert := assert.New(t)
	f, err := ioutil.TempFile("", "")
	must(err)
	defer f.Close()
	defer os.Remove(f.Name())

	assertSize := func(expected int64) {
		s, err := f.Stat()
		must(err)

		assert.Equal(expected, s.Size())
	}
	assertSize(0)

	must(ox.Preallocate(f, 2048))

	assertSize(2048)

	_, err = f.Seek(0, io.SeekStart)
	must(err)

	_, err = f.Write([]byte("hello"))
	must(err)

	assertSize(2048)

	must(ox.Preallocate(f, 4096))
	assertSize(4096)

	buf := make([]byte, 5)
	n, err := f.ReadAt(buf, 0)
	must(err)
	assert.Equal(5, n)

	assert.Equal("hello", string(buf))
}

func Test_Preallocate(t *testing.T) {
	t.Logf("With fallocate...")
	ox.SIMULATE_FALLOCATE_NOT_SUPPORTED = true
	doTestPreallocate(t)

	t.Logf("Without fallocate...")
	ox.SIMULATE_FALLOCATE_NOT_SUPPORTED = false
	doTestPreallocate(t)
}

func must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
