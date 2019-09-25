//+build !windows

package ox

import (
	"io"
	"os"

	"github.com/detailyang/go-fallocate"
	"github.com/pkg/errors"
)

// Reserve `size` bytes of space for f, in the
// quickest way possible. f must be opened with O_RDWR.
func Preallocate(f *os.File, size int64) error {
	currentSize, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return errors.WithStack(err)
	}

	remaining := size - currentSize
	if remaining <= 0 {
		return nil
	}

	err = fallocate.Fallocate(f, currentSize, remaining)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
