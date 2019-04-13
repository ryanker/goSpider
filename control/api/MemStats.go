package api

import (
	"fmt"
	"runtime"
	"time"

	"../../lib/misc"

	"github.com/xiuno/gin"
)

type MemStats struct {
	Uptime       string // 程序运行时间
	NumCPU       int    // 当前可用的 CPU 数
	NumGoroutine int    // 当前存在的 GO 程数

	// 基本信息
	Alloc      string // 当前使用内存
	TotalAlloc string // 累积使用内存
	Sys        string // 系统使用内存
	Lookups    uint64 // 指针查找次数
	Mallocs    uint64 // 内存分配次数
	Frees      uint64 // 内存释放次数

	// Heap 基本信息
	HeapAlloc    string // 当前 Heap 内存使用量
	HeapSys      string // Heap 内存占用量
	HeapIdle     string // Heap 内存空闲量
	HeapInuse    string // 正在使用的 Heap 内存
	HeapReleased string // 被释放的 Heap 内存
	HeapObjects  uint64 // Heap 对象数量

	// 其他信息
	StackInuse  string // 启动 Stack 使用量
	StackSys    string // 系统分配的 Stack 内存
	MSpanInuse  string // MSpan 结构内存使用量
	MSpanSys    string // 系统分配的 MSpan 结构内存
	MCacheInuse string // MCache 结构内存使用量
	MCacheSys   string // 系统分配的 MCache 结构内存
	BuckHashSys string // 系统分配的剖析哈希表内存
	GCSys       string // 系统分配的 GC 元数据内存
	OtherSys    string // 其它被分配的系统内存

	// 垃圾回收信息
	NextGC       string // 下次 GC 内存回收量
	LastGC       string // 距离上次 GC 时间
	PauseTotalNs string // GC 暂停时间总量
	PauseNs      string // 上次 GC 暂停时间
	NumGC        uint32 // GC 执行次数
}

// 程序运行时间
var Uptime = time.Now()

// 内存信息
func MemStatsInfo(c *gin.Context) {
	ms := new(runtime.MemStats)
	runtime.ReadMemStats(ms)

	mem := MemStats{}
	mem.Uptime = Uptime.Format("2006-01-02 15:04:05")
	mem.NumCPU = runtime.NumCPU()
	mem.NumGoroutine = runtime.NumGoroutine()

	// 基本信息
	mem.Alloc = misc.HumanSize(ms.Alloc)
	mem.TotalAlloc = misc.HumanSize(ms.TotalAlloc)
	mem.Sys = misc.HumanSize(ms.Sys)
	mem.Lookups = ms.Lookups
	mem.Mallocs = ms.Mallocs
	mem.Frees = ms.Frees

	// Heap 基本信息
	mem.HeapAlloc = misc.HumanSize(ms.HeapAlloc)
	mem.HeapSys = misc.HumanSize(ms.HeapSys)
	mem.HeapIdle = misc.HumanSize(ms.HeapIdle)
	mem.HeapInuse = misc.HumanSize(ms.HeapInuse)
	mem.HeapReleased = misc.HumanSize(ms.HeapReleased)
	mem.HeapObjects = ms.HeapObjects

	// 其他信息
	mem.StackInuse = misc.HumanSize(ms.StackInuse)
	mem.StackSys = misc.HumanSize(ms.StackSys)
	mem.MSpanInuse = misc.HumanSize(ms.MSpanInuse)
	mem.MSpanSys = misc.HumanSize(ms.MSpanSys)
	mem.MCacheInuse = misc.HumanSize(ms.MCacheInuse)
	mem.MCacheSys = misc.HumanSize(ms.MCacheSys)
	mem.BuckHashSys = misc.HumanSize(ms.BuckHashSys)
	mem.GCSys = misc.HumanSize(ms.GCSys)
	mem.OtherSys = misc.HumanSize(ms.OtherSys)

	// 垃圾回收信息
	ts := float64(time.Second)
	mem.NextGC = misc.HumanSize(ms.NextGC)
	mem.LastGC = fmt.Sprintf("%.4f S", float64(time.Now().UnixNano()-int64(ms.LastGC))/ts)
	mem.PauseTotalNs = fmt.Sprintf("%.4f S", float64(ms.PauseTotalNs)/ts)
	mem.PauseNs = fmt.Sprintf("%.6f S", float64(ms.PauseNs[(ms.NumGC+255)%256])/ts)
	mem.NumGC = ms.NumGC

	c.Message("0", "success", mem)
}
