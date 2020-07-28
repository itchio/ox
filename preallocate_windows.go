//+build windows

package ox

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

// Note: this does nothing on Windows
var SIMULATE_FALLOCATE_NOT_SUPPORTED = false

// Reserve `size` bytes of space for f, in the
// quickest way possible. f must be opened with O_RDWR.
func Preallocate(f *os.File, size int64) error {
	_, err := f.Seek(size, io.SeekStart)
	if err != nil {
		return errors.WithStack(err)
	}

	err = windows.SetEndOfFile(windows.Handle(f.Fd()))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
