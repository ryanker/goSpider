package api

import (
	"syscall"

	"../../lib/misc"

	"github.com/xiuno/gin"
)

type Disk struct {
	All  string // 总大小
	Used string // 使用大小
	Free string // 剩余大小
}

// 磁盘信息，暂不考虑 windows 系统。参考: https://www.jianshu.com/p/f3d31f84d95d
func DiskInfo(c *gin.Context) {
	disk := Disk{}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("./", &fs)
	if err == nil {
		disk.All = misc.HumanSize(fs.Blocks * uint64(fs.Bsize))
		disk.Free = misc.HumanSize(fs.Bfree * uint64(fs.Bsize))
		disk.Used = misc.HumanSize((fs.Blocks - fs.Bfree) * uint64(fs.Bsize))
	}
	c.Message("0", "success", disk)
}
