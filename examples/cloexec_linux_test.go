package test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// requireCloexec asserts that the process holds an open descriptor for the
// given device path and that it has the close-on-exec flag set, so the raw
// device fd is not leaked to child processes across execve.
//
// The smart.go device types keep their fd unexported, so we locate it by
// scanning /proc/self/fd for the symlink that resolves to path.
func requireCloexec(t *testing.T, path string) {
	t.Helper()

	entries, err := os.ReadDir("/proc/self/fd")
	require.NoError(t, err)

	found := false
	for _, e := range entries {
		target, err := os.Readlink(filepath.Join("/proc/self/fd", e.Name()))
		if err != nil || target != path {
			continue
		}

		fd, err := strconv.Atoi(e.Name())
		require.NoError(t, err)

		flags, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
		require.NoError(t, err)
		require.NotZero(t, flags&unix.FD_CLOEXEC, "device fd for %s must have FD_CLOEXEC set", path)
		found = true
	}

	require.True(t, found, "no open fd found for %s", path)
}
