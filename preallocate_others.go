//+build !windows

package ox

import (
	"io"
	"os"
	"syscall"

	"github.com/detailyang/go-fallocate"
	"github.com/pkg/errors"
)

var SIMULATE_FALLOCATE_NOT_SUPPORTED = false

func Fallocate(file *os.File, offset int64, length int64) error {
	if SIMULATE_FALLOCATE_NOT_SUPPORTED {
		return syscall.ENOTSUP
	} else {
		return fallocate.Fallocate(file, offset, length)
	}
}

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

	err = Fallocate(f, currentSize, remaining)
	if err != nil {
		if errors.Is(err, syscall.ENOTSUP) {
			// as of July 2020, we've seen this error condition happpen
			// on Linux with NTFS, eCryptFS, and ZFS partitions. ext* are fine.

			// this is the slower fallback:
			_, err = f.Seek(currentSize, io.SeekStart)
			if err != nil {
				return errors.Wrapf(err, "while pre-allocating %v bytes with fallback", size)
			}

			_, err = io.Copy(f, io.LimitReader(&zeroReader{}, remaining))
			if err != nil {
				return errors.Wrapf(err, "while pre-allocating %v bytes with fallback", size)
			}

			return nil
		}

		if err != nil {
			return errors.Wrapf(err, "while pre-allocating %v bytes with fallocate", size)
		}
	}

	return nil
}

type zeroReader struct{}

var _ io.Reader = (*zeroReader)(nil)

func (zr *zeroReader) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		p[i] = 0
	}
	return len(p), nil
}
