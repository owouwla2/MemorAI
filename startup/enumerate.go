package startup

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// 注册表 Run 路径列表
type runLocation struct {
	root      registry.Key
	path      string // Run 键路径
	approved  string // StartupApproved\Run 路径
	scope     string
	needAdmin bool
}

var runLocations = []runLocation{
	{
		root:      registry.CURRENT_USER,
		path:      `Software\Microsoft\Windows\CurrentVersion\Run`,
		approved:  `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`,
		scope:     ScopeUser,
		needAdmin: false,
	},
	{
		root:      registry.LOCAL_MACHINE,
		path:      `Software\Microsoft\Windows\CurrentVersion\Run`,
		approved:  `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`,
		scope:     ScopeSystem,
		needAdmin: true,
	},
	{
		root:      registry.LOCAL_MACHINE,
		path:      `Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`,
		approved:  `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32`,
		scope:     ScopeSystem,
		needAdmin: true,
	},
}

// 启动文件夹来源
type folderLocation struct {
	folder           string
	approvedRoot     registry.Key
	approvedPath     string
	scope            string
	needAdmin        bool
}

func folderLocations() []folderLocation {
	userStartup := filepath.Join(os.Getenv("APPDATA"), `Microsoft\Windows\Start Menu\Programs\Startup`)
	commonStartup := filepath.Join(os.Getenv("ProgramData"), `Microsoft\Windows\Start Menu\Programs\Startup`)
	return []folderLocation{
		{
			folder:       userStartup,
			approvedRoot: registry.CURRENT_USER,
			approvedPath: `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\StartupFolder`,
			scope:        ScopeUser,
			needAdmin:    false,
		},
		{
			folder:       commonStartup,
			approvedRoot: registry.LOCAL_MACHINE,
			approvedPath: `Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\StartupFolder`,
			scope:        ScopeSystem,
			needAdmin:    true,
		},
	}
}

// makeID 生成稳定的唯一ID
func makeID(parts ...string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// readApprovedStatus 读取某项的 StartupApproved 状态
// approvedRoot/approvedPath 指向 ...\StartupApproved\Run 或 \StartupFolder
// valueName 是注册表值名（与 Run 键里的值名一致）
// 返回 (启用?, 是否找到记录)
func readApprovedStatus(root registry.Key, path, valueName string) (bool, bool) {
	k, err := registry.OpenKey(root, path, registry.READ)
	if err != nil {
		// approved 子键不存在，视为"未记录"，按启用处理
		return true, false
	}
	defer k.Close()

	val, _, err := k.GetBinaryValue(valueName)
	if err != nil {
		return true, false
	}
	if len(val) == 0 {
		return true, false
	}
	// 第一字节: 02 启用, 03 禁用 (其余 0x06 等也算启用)
	return val[0] != 0x03, true
}

// EnumerateStartupItems 列出所有自启项
func EnumerateStartupItems() ([]*StartupItem, error) {
	items := make([]*StartupItem, 0, 32)

	// 1) 注册表 Run 键
	for _, loc := range runLocations {
		k, err := registry.OpenKey(loc.root, loc.path, registry.READ)
		if err != nil {
			continue // 不存在/无权限就跳过
		}
		names, _ := k.ReadValueNames(0)
		for _, name := range names {
			cmd, _, err := k.GetStringValue(name)
			if err != nil {
				continue
			}
			enabled, _ := readApprovedStatus(loc.root, loc.approved, name)
			items = append(items, &StartupItem{
				ID:        makeID(loc.scope, "registry", loc.path, name),
				Name:      name,
				Command:   cmd,
				Source:    SourceRegistry,
				Location:  rootName(loc.root) + `\` + loc.path,
				Scope:     loc.scope,
				Enabled:   enabled,
				NeedAdmin: loc.needAdmin,
			})
		}
		k.Close()
	}

	// 2) 启动文件夹
	for _, loc := range folderLocations() {
		entries, err := os.ReadDir(loc.folder)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			// 通常是 .lnk
			full := filepath.Join(loc.folder, name)
			cmd := full
			// 解析 .lnk 比较麻烦（需要COM），这里直接显示路径
			enabled, _ := readApprovedStatus(loc.approvedRoot, loc.approvedPath, name)
			items = append(items, &StartupItem{
				ID:        makeID(loc.scope, "folder", loc.folder, name),
				Name:      name,
				Command:   cmd,
				Source:    SourceFolder,
				Location:  loc.folder,
				Scope:     loc.scope,
				Enabled:   enabled,
				NeedAdmin: loc.needAdmin,
			})
		}
	}

	// 排序：按启用/禁用 + 名称
	sort.Slice(items, func(i, j int) bool {
		if items[i].Enabled != items[j].Enabled {
			return items[i].Enabled
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

// rootName 注册表根名称（用于显示）
func rootName(r registry.Key) string {
	switch r {
	case registry.CURRENT_USER:
		return "HKCU"
	case registry.LOCAL_MACHINE:
		return "HKLM"
	case registry.CLASSES_ROOT:
		return "HKCR"
	case registry.USERS:
		return "HKU"
	}
	return "HK?"
}

// findItem 根据 ID 重新查找 item（避免前端传入错误）
func findItem(id string) (*StartupItem, error) {
	items, err := EnumerateStartupItems()
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if it.ID == id {
			return it, nil
		}
	}
	return nil, ErrItemNotFound
}
