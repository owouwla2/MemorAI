package memory

import (
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMemoryInfo 系统内存概览
type SystemMemoryInfo struct {
	Total       uint64  `json:"total"`        // 总内存（字节）
	Available   uint64  `json:"available"`    // 可用内存
	Used        uint64  `json:"used"`         // 已用内存
	UsedPercent float64 `json:"usedPercent"`  // 占用率
	Free        uint64  `json:"free"`         // 空闲
	// 格式化后的可读字符串（GB）
	TotalGB     float64 `json:"totalGB"`
	UsedGB      float64 `json:"usedGB"`
	AvailableGB float64 `json:"availableGB"`

	// 主机信息
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platformVersion"`
	UptimeSeconds   uint64 `json:"uptimeSeconds"`
}

// GetSystemMemory 获取系统内存信息
func GetSystemMemory() (*SystemMemoryInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	info := &SystemMemoryInfo{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
		Free:        v.Free,
		TotalGB:     bytesToGB(v.Total),
		UsedGB:      bytesToGB(v.Used),
		AvailableGB: bytesToGB(v.Available),
	}

	// 主机信息（失败不算致命错误）
	if h, err := host.Info(); err == nil {
		info.Hostname = h.Hostname
		info.OS = h.OS
		info.Platform = h.Platform
		info.PlatformVersion = h.PlatformVersion
		info.UptimeSeconds = h.Uptime
	}

	return info, nil
}

func bytesToGB(b uint64) float64 {
	return float64(b) / 1024.0 / 1024.0 / 1024.0
}

func bytesToMB(b uint64) float64 {
	return float64(b) / 1024.0 / 1024.0
}
