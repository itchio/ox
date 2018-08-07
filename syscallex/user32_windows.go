package syscallex

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	moduser32 = windows.NewLazySystemDLL("user32.dll")

	procEnumWindows              = moduser32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = moduser32.NewProc("GetWindowThreadProcessId")
	procSetForegroundWindow      = moduser32.NewProc("SetForegroundWindow")
	procShowWindow               = moduser32.NewProc("ShowWindow")
	procIsWindowVisible          = moduser32.NewProc("IsWindowVisible")
	procGetWindowLongPtr         = moduser32.NewProc("GetWindowLongPtr")
)

// GetWindowLong and GetWindowLongPtr constants
const (
	GWL_EXSTYLE     = -20
	GWL_STYLE       = -16
	GWL_WNDPROC     = -4
	GWLP_WNDPROC    = -4
	GWL_HINSTANCE   = -6
	GWLP_HINSTANCE  = -6
	GWL_HWNDPARENT  = -8
	GWLP_HWNDPARENT = -8
	GWL_ID          = -12
	GWLP_ID         = -12
	GWL_USERDATA    = -21
	GWLP_USERDATA   = -21
)

func EnumWindows(
	cb uintptr,
	lparam uintptr,
) (err error) {
	r1, _, e1 := syscall.Syscall(
		procEnumWindows.Addr(),
		2,
		cb,
		lparam,
		0,
	)
	if r1 == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowThreadProcessId(
	hwnd syscall.Handle,
	pProcessId *uint32,
) uint32 {
	r1, _, _ := syscall.Syscall(
		procGetWindowThreadProcessId.Addr(),
		2,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pProcessId)),
		0,
	)
	return uint32(r1)
}

func SetForegroundWindow(
	hwnd syscall.Handle,
) (err error) {
	r1, _, e1 := syscall.Syscall(
		procSetForegroundWindow.Addr(),
		1,
		uintptr(hwnd),
		0,
		0,
	)
	if r1 == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ShowWindow(
	hwnd syscall.Handle,
	flags int,
) (err error) {
	r1, _, e1 := syscall.Syscall(
		procShowWindow.Addr(),
		2,
		uintptr(hwnd),
		uintptr(flags),
		0,
	)
	if r1 == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func IsWindowVisible(
	hwnd syscall.Handle,
) bool {
	ret, _, _ := syscall.Syscall(
		procGetWindowLongPtr.Addr(),
		1,
		uintptr(hwnd),
		0,
		0,
	)

	return ret != 0
}

func GetWindowLongPtr(
	hwnd syscall.Handle,
	index int32,
) uintptr {
	ret, _, _ := syscall.Syscall(
		procGetWindowLongPtr.Addr(),
		2,
		uintptr(hwnd),
		uintptr(index),
		0,
	)

	return ret
}
