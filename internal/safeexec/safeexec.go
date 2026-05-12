// Package safeexec validates binaries that are about to be invoked under
// elevated privileges (sudo / osascript "administrator privileges"). It
// rejects paths in user-writable locations or owned by foreign uids, closing
// the obvious privilege-escalation hole where a malicious binary in $PATH
// inherits root via mkdev's helper invocations.
package safeexec

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// VerifyBinPath rejects bin paths that aren't safe to invoke under sudo.
// On macOS we require the binary to be a regular file owned by root or the
// current uid, with no group/world write bit set. Symlinks are followed.
func VerifyBinPath(bin string) error {
	resolved, err := filepath.EvalSymlinks(bin)
	if err != nil {
		return fmt.Errorf("safeexec: resolve %s: %w", bin, err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return fmt.Errorf("safeexec: stat %s: %w", resolved, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("safeexec: %s is not a regular file", resolved)
	}
	if info.Mode().Perm()&0o022 != 0 {
		return fmt.Errorf("safeexec: %s is group/world writable; refusing to invoke under sudo", resolved)
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := stat.Uid
		if uid != 0 && int(uid) != os.Getuid() {
			return fmt.Errorf("safeexec: %s owned by uid %d (not root or current user); refusing to invoke under sudo", resolved, uid)
		}
	}
	return nil
}
