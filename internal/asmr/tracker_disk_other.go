//go:build !windows

package asmr

import "golang.org/x/sys/unix"

// freeDiskBytes 返回 path 所在卷的剩余字节数（D6 磁盘保护，非 Windows 平台）。
func freeDiskBytes(path string) (uint64, error) {
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return 0, err
	}
	return st.Bavail * uint64(st.Bsize), nil
}
