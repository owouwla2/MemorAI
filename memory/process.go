package memory

import (
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 单个进程的内存信息
type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`        // 进程名（如 chrome.exe）
	ExePath    string  `json:"exePath"`     // 完整路径（可能为空）
	MemoryRSS  uint64  `json:"memoryRSS"`   // 物理内存占用（字节）
	MemoryMB   float64 `json:"memoryMB"`    // MB
	MemoryPct  float64 `json:"memoryPct"`   // 占系统总内存百分比
}

// ListProcesses 列举所有进程及内存占用，按内存降序
func ListProcesses() ([]*ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	sysMem, err := GetSystemMemory()
	if err != nil {
		return nil, err
	}
	totalMem := sysMem.Total

	results := make([]*ProcessInfo, 0, len(procs))
	for _, p := range procs {
		// 进程可能在枚举过程中退出，错误忽略
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		memInfo, err := p.MemoryInfo()
		if err != nil || memInfo == nil {
			continue
		}

		exePath, _ := p.Exe() // 失败给空串

		pi := &ProcessInfo{
			PID:       p.Pid,
			Name:      name,
			ExePath:   exePath,
			MemoryRSS: memInfo.RSS,
			MemoryMB:  bytesToMB(memInfo.RSS),
		}
		if totalMem > 0 {
			pi.MemoryPct = float64(memInfo.RSS) / float64(totalMem) * 100
		}
		results = append(results, pi)
	}

	// 按内存降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].MemoryRSS > results[j].MemoryRSS
	})

	return results, nil
}

// LowerName 返回小写的进程名（用于分类匹配）
func LowerName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
