package startup

// StartupItem 一条自启项
type StartupItem struct {
	ID       string `json:"id"`       // 唯一ID，前端用来回传
	Name     string `json:"name"`     // 显示名（注册表值名或文件名）
	Command  string `json:"command"`  // 启动命令或快捷方式目标
	Source   string `json:"source"`   // "registry" 或 "folder"
	Location string `json:"location"` // 注册表完整路径或文件夹路径
	Scope    string `json:"scope"`    // "user" 或 "system"
	Enabled  bool   `json:"enabled"`
	// 是否需要管理员权限才能修改（HKLM/系统启动文件夹）
	NeedAdmin bool `json:"needAdmin"`
	// 关联进程内存（可选，调用方填）
	MemoryMB float64 `json:"memoryMB,omitempty"`
}

// approvedStatus 描述 StartupApproved 的 12字节值
//   启用: 02 00 00 00 + 8字节 FILETIME（全 0 即可）
//   禁用: 03 00 00 00 + 8字节 FILETIME（禁用时的时间）
type approvedStatus struct {
	Enabled bool
	// 原始值（如果存在）
	Raw []byte
}

const (
	SourceRegistry = "registry"
	SourceFolder   = "folder"
	ScopeUser      = "user"
	ScopeSystem    = "system"
)
