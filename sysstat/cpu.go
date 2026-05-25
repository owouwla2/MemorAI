package sysstat

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUInfo CPU 信息
type CPUInfo struct {
	UsedPercent float64 `json:"usedPercent"` // 总占用率
	Cores       int     `json:"cores"`       // 物理核心数
	LogicalCores int    `json:"logicalCores"`// 逻辑核心数
	ModelName   string  `json:"modelName"`
	MHz         float64 `json:"mhz"`
	// 各核心占用率（可选展示）
	PerCore []float64 `json:"perCore"`
}

// GetCPU 获取一次 CPU 状态
// 第一次调用会有 200ms 阻塞（gopsutil 需要采样窗口计算占用率）
func GetCPU() (*CPUInfo, error) {
	info := &CPUInfo{}

	// 总占用率（非阻塞 0 = 直接返回上次的差值；首次会得到 0，所以这里给 200ms）
	all, err := cpu.Percent(200*time.Millisecond, false)
	if err == nil && len(all) > 0 {
		info.UsedPercent = all[0]
	}

	// 各核心占用率
	per, err := cpu.Percent(0, true)
	if err == nil {
		info.PerCore = per
	}

	// 物理/逻辑核心数
	if phy, err := cpu.Counts(false); err == nil {
		info.Cores = phy
	}
	if log, err := cpu.Counts(true); err == nil {
		info.LogicalCores = log
	}

	// 型号、频率
	if cs, err := cpu.Info(); err == nil && len(cs) > 0 {
		info.ModelName = cs[0].ModelName
		info.MHz = cs[0].Mhz
	}

	return info, nil
}

// GetCPUQuick 不阻塞的快速采样（用于循环采样场景）
// 调用前需要先用 cpu.Percent(0, false) "预热" 一下
func GetCPUQuick() float64 {
	all, err := cpu.Percent(0, false)
	if err != nil || len(all) == 0 {
		return 0
	}
	return all[0]
}

// WarmupCPU 预热 CPU 计数器（让下次的 0-interval 调用能拿到正确的值）
// 在采样开始前调用
func WarmupCPU() {
	_, _ = cpu.Percent(0, false)
}
