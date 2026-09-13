//go:build windows

package asmr

import "golang.org/x/sys/windows"

// freeDiskBytes 返回 path 所在卷的剩余字节数（D6 磁盘保护）。
func freeDiskBytes(path string) (uint64, error) {
	var free, total, avail uint64
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	if err := windows.GetDiskFreeSpaceEx(p, &free, &total, &avail); err != nil {
		return 0, err
	}
	return free, nil
}
