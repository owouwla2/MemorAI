package memory

import "memory-analyzer/sysstat"

// AnalysisReport 完整的内存分析报告（提供给前端和AI）
type AnalysisReport struct {
	System    *SystemMemoryInfo `json:"system"`
	CPU       *sysstat.CPUInfo  `json:"cpu"`        // CPU 状态
	Groups    []*SoftwareGroup  `json:"groups"`     // 按软件聚合
	TopProcs  []*ProcessInfo    `json:"topProcs"`   // 单进程 Top N
	Total     int               `json:"total"`      // 总进程数
}

// Analyze 执行一次完整采集
// topN: 单进程列表只返回 Top N 个，<=0 表示返回全部
func Analyze(topN int) (*AnalysisReport, error) {
	sys, err := GetSystemMemory()
	if err != nil {
		return nil, err
	}

	procs, err := ListProcesses()
	if err != nil {
		return nil, err
	}

	groups := GroupBySoftware(procs, sys.Total)

	top := procs
	if topN > 0 && len(procs) > topN {
		top = procs[:topN]
	}

	// CPU（失败不致命）
	cpuInfo, _ := sysstat.GetCPU()

	return &AnalysisReport{
		System:   sys,
		CPU:      cpuInfo,
		Groups:   groups,
		TopProcs: top,
		Total:    len(procs),
	}, nil
}
