package execas

import (
	"strings"
	"syscall"
	"testing"

	"github.com/itchio/ox/syscallex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunWithoutSysProcAttr(t *testing.T) {
	cmd := Command("cmd", "/c", "exit 0")

	err := cmd.Run()
	require.NoError(t, err)
	require.NotNil(t, cmd.ProcessState)
	assert.True(t, cmd.ProcessState.Success())
}

func TestOutputWithoutSysProcAttr(t *testing.T) {
	cmd := Command("cmd", "/c", "echo hello")

	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "hello", strings.TrimSpace(string(output)))
}

func TestRunWithSysProcAttrExposesHandles(t *testing.T) {
	sys := &syscallex.SysProcAttr{}
	cmd := Command("cmd", "/c", "exit 0")
	cmd.SysProcAttr = sys

	err := cmd.Run()
	require.NoError(t, err)
	require.NotZero(t, sys.ProcessHandle)
	require.NotZero(t, sys.ThreadHandle)

	assert.NoError(t, syscall.CloseHandle(sys.ThreadHandle))
	assert.NoError(t, syscall.CloseHandle(sys.ProcessHandle))
}
