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
	procSwitchToThisWindow       = moduser32.NewProc("SwitchToThisWindow")
	procFindWindow               = moduser32.NewProc("FindWindowW")
	procSendInput                = moduser32.NewProc("SendInput")
	procGetWindowRect            = moduser32.NewProc("GetWindowRect")
	procGetSystemMetrics         = moduser32.NewProc("GetSystemMetrics")
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
		procIsWindowVisible.Addr(),
		1,
		uintptr(hwnd),
		0,
		0,
	)

	return ret != 0
}

func SwitchToThisWindow(
	hwnd syscall.Handle,
	altTab bool,
) {
	altTabInt := 0
	if altTab {
		altTabInt = 1
	}

	syscall.Syscall(
		procSwitchToThisWindow.Addr(),
		2,
		uintptr(hwnd),
		uintptr(altTabInt),
		0,
	)
}

func FindWindow(cls string, win string) (syscall.Handle, error) {
	r0, _, e1 := syscall.Syscall(
		procFindWindow.Addr(), 2,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(cls))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(win))),
		0,
	)
	var err error
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return syscall.Handle(r0), err
}

func FindWindowByClass(cls string) (syscall.Handle, error) {
	r0, _, e1 := syscall.Syscall(
		procFindWindow.Addr(), 2,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(cls))),
		0,
		0,
	)
	var err error
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return syscall.Handle(r0), err
}

type INPUT struct {
	Type  uint32
	Input MOUSEINPUT
}

var (
	INPUT_MOUSE    uint32 = 0
	INPUT_KEYBOARD uint32 = 1
	INPUT_HARDWARE uint32 = 2
)

func (i *INPUT) SetKeyboardInput(ki KEYBDINPUT) {
	i.Type = INPUT_KEYBOARD
	var mip = (&i.Input)
	*mip = MOUSEINPUT{}

	var p = (*KEYBDINPUT)((unsafe.Pointer)(mip))
	*p = ki
}

func (i *INPUT) SetMouseInput(mi MOUSEINPUT) {
	i.Type = INPUT_MOUSE
	var mip = (&i.Input)
	*mip = mi
}

type MOUSEINPUT struct {
	X         int32
	Y         int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo *uint32
}

type KEYBDINPUT struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo *uint32
}

type HARDWAREINPUT struct {
	UMsg   uint32
	ParamL int16
	ParamH int16
}

var (
	KEYEVENTF_EXTENDEDKEY uint32 = 0x0001
	KEYEVENTF_KEYUP       uint32 = 0x0002
	KEYEVENTF_UNICODE     uint32 = 0x0004
	KEYEVENTF_SCANCODE    uint32 = 0x0008
)

var (
	MOUSEEVENTF_ABSOLUTE        uint32 = 0x8000
	MOUSEEVENTF_HWHEEL          uint32 = 0x1000
	MOUSEEVENTF_MOVE            uint32 = 0x0001
	MOUSEEVENTF_MOVE_NOCOALESCE uint32 = 0x2000
	MOUSEEVENTF_LEFTDOWN        uint32 = 0x0002
	MOUSEEVENTF_LEFTUP          uint32 = 0x0004
	MOUSEEVENTF_RIGHTDOWN       uint32 = 0x0008
	MOUSEEVENTF_RIGHTUP         uint32 = 0x0010
	MOUSEEVENTF_MIDDLEDOWN      uint32 = 0x0020
	MOUSEEVENTF_MIDDLEUP        uint32 = 0x0040
	MOUSEEVENTF_VIRTUALDESK     uint32 = 0x4000
	MOUSEEVENTF_WHEEL           uint32 = 0x0800
	MOUSEEVENTF_XDOWN           uint32 = 0x0080
	MOUSEEVENTF_XUP             uint32 = 0x0100
)

func SendMouseInput(mi MOUSEINPUT) (err error) {
	var i INPUT
	i.SetMouseInput(mi)
	return SendInput(i)
}

func SendKeyboardInput(ki KEYBDINPUT) (err error) {
	var i INPUT
	i.SetKeyboardInput(ki)
	return SendInput(i)
}

func SendInput(input INPUT) (err error) {
	r0, _, e1 := syscall.Syscall(
		procSendInput.Addr(), 3,
		1,
		uintptr(unsafe.Pointer(&input)),
		unsafe.Sizeof(INPUT{}),
	)
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func GetWindowRect(hwnd syscall.Handle) (rect RECT, err error) {
	r0, _, e1 := syscall.Syscall(
		procGetWindowRect.Addr(), 2,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)),
		0,
	)
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

var (
	SM_CXSCREEN int = 0
	SM_CYSCREEN     = 1
)

func GetSystemMetrics(nIndex int) (ret int) {
	r0, _, _ := syscall.Syscall(
		procGetSystemMetrics.Addr(), 1,
		uintptr(nIndex),
		0, 0,
	)
	ret = int(r0)
	return
}
