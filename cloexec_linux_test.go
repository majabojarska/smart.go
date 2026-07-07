package smart

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// TestOpenSetsCloexec verifies that device descriptors are opened with the
// close-on-exec flag so they are not leaked to child processes across execve.
//
// OpenNVMe only opens the path O_RDONLY and returns without issuing any ioctl,
// so it can be exercised against a regular file. OpenSata/OpenScsi share the
// identical unix.Open call but reject a regular file at their inquiry ioctl;
// their descriptors are checked against real devices in the QEMU integration
// tests (see examples/*_linux_test.go).
func TestOpenSetsCloexec(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "fakedev")
	require.NoError(t, os.WriteFile(path, nil, 0o600))

	dev, err := OpenNVMe(path)
	require.NoError(t, err)
	defer dev.Close()

	flags, err := unix.FcntlInt(uintptr(dev.fd), unix.F_GETFD, 0)
	require.NoError(t, err)
	require.NotZero(t, flags&unix.FD_CLOEXEC, "device fd must have FD_CLOEXEC set")
}
