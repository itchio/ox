package osex

import (
	"os"
	"syscall"

	"github.com/itchio/ox/syscallex"
)

// ProcAttr holds the attributes that will be applied to a new process
// started by StartProcess.
type ProcAttr struct {
	// If Dir is non-empty, the child changes into the directory before
	// creating the process.
	Dir string
	// If Env is non-nil, it gives the environment variables for the
	// new process in the form returned by Environ.
	// If it is nil, the result of Environ will be used.
	Env []string
	// Files specifies the open files inherited by the new process. The
	// first three entries correspond to standard input, standard output, and
	// standard error. An implementation may support additional entries,
	// depending on the underlying operating system. A nil entry corresponds
	// to that file being closed when the process starts.
	Files []*os.File

	// Operating system-specific process creation attributes.
	// Note that setting this field means that your program
	// may not execute properly or even compile on some
	// operating systems.
	Sys *syscallex.SysProcAttr
}

// StartProcess starts a new process with the program, arguments and attributes
// specified by name, argv and attr.
//
// StartProcess is a low-level interface. The os/exec package provides
// higher-level interfaces.
//
// If there is an error, it will be of type *PathError.
func StartProcessWithLogon(name string, argv []string, username string, domain string, password string, attr *ProcAttr) (*os.Process, error) {
	return startProcessWithLogon(name, argv, username, domain, password, attr)
}

func startProcessWithLogon(name string, argv []string, username string, domain string, password string, attr *ProcAttr) (p *os.Process, err error) {
	// If there is no SysProcAttr (ie. no Chroot or changed
	// UID/GID), double-check existence of the directory we want
	// to chdir into. We can make the error clearer this way.
	if attr != nil && attr.Sys == nil && attr.Dir != "" {
		if _, err := os.Stat(attr.Dir); err != nil {
			pe := err.(*os.PathError)
			pe.Op = "chdir"
			return nil, pe
		}
	}

	sysattr := &syscallex.ProcAttr{
		Dir: attr.Dir,
		Env: attr.Env,
		Sys: attr.Sys,
	}
	if sysattr.Env == nil {
		sysattr.Env = os.Environ()
	}
	for _, f := range attr.Files {
		sysattr.Files = append(sysattr.Files, f.Fd())
	}

	pid, h, e := syscallex.StartProcessWithLogon(name, argv, username, domain, password, sysattr)
	if e != nil {
		return nil, &os.PathError{"fork/exec", name, e}
	}

	// Historically we returned *os.Process by constructing os.Process{Pid: pid}
	// and writing the private `handle` field via reflection/unsafe so Wait/Kill
	// used the CreateProcessWithLogon handle directly. Go 1.24 changed
	// os.Process internals, so that hack is no longer valid. Re-open by PID
	// using the supported API and close our raw handle unless the caller
	// explicitly requested it through SysProcAttr.ProcessHandle.
	p, err = os.FindProcess(pid)
	if err != nil {
		syscall.CloseHandle(syscall.Handle(h))
		return nil, &os.PathError{"findprocess", name, err}
	}

	if attr == nil || attr.Sys == nil || attr.Sys.ProcessHandle == 0 {
		syscall.CloseHandle(syscall.Handle(h))
	}
	return p, nil
}
