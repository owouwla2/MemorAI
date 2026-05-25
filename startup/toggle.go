package startup

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

var (
	ErrItemNotFound  = errors.New("自启项不存在")
	ErrAccessDenied  = errors.New("权限不足，可能需要以管理员身份运行")
)

// SetEnabled 启用或禁用一个自启项
// 不删除原始项，只通过 StartupApproved 注册表项控制（同任务管理器）
func SetEnabled(id string, enable bool) error {
	item, err := findItem(id)
	if err != nil {
		return err
	}

	// 找到对应的 approved 路径
	root, approvedPath, valueName, err := resolveApprovedTarget(item)
	if err != nil {
		return err
	}

	// 构造 12 字节值
	val := make([]byte, 12)
	if enable {
		val[0] = 0x02
	} else {
		val[0] = 0x03
		// 写入禁用时间（FILETIME, 100ns intervals from 1601-01-01）
		ft := timeToFILETIME(time.Now())
		binary.LittleEndian.PutUint64(val[4:], ft)
	}

	// 打开/创建 approved 键
	k, _, err := registry.CreateKey(root, approvedPath, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		if isAccessDenied(err) {
			return fmt.Errorf("%w: 修改 %s 需要管理员权限", ErrAccessDenied, approvedPath)
		}
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	if err := k.SetBinaryValue(valueName, val); err != nil {
		if isAccessDenied(err) {
			return fmt.Errorf("%w: 写入 %s 需要管理员权限", ErrAccessDenied, approvedPath)
		}
		return fmt.Errorf("写入注册表失败: %w", err)
	}
	return nil
}

// resolveApprovedTarget 根据 item 找到正确的 approved 注册表位置和值名
func resolveApprovedTarget(it *StartupItem) (registry.Key, string, string, error) {
	if it.Source == SourceRegistry {
		// 从 it.Location（如 "HKCU\Software\..."）反查
		for _, loc := range runLocations {
			locFull := rootName(loc.root) + `\` + loc.path
			if strings.EqualFold(locFull, it.Location) {
				return loc.root, loc.approved, it.Name, nil
			}
		}
		return 0, "", "", fmt.Errorf("无法识别注册表位置: %s", it.Location)
	}
	if it.Source == SourceFolder {
		for _, loc := range folderLocations() {
			if strings.EqualFold(loc.folder, it.Location) {
				return loc.approvedRoot, loc.approvedPath, it.Name, nil
			}
		}
		return 0, "", "", fmt.Errorf("无法识别启动文件夹: %s", it.Location)
	}
	return 0, "", "", fmt.Errorf("未知的来源: %s", it.Source)
}

// isAccessDenied 简单判断是否权限错误
func isAccessDenied(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "access is denied") ||
		strings.Contains(s, "拒绝访问")
}

// timeToFILETIME 将 time.Time 转为 Windows FILETIME（uint64, 100ns intervals since 1601-01-01）
func timeToFILETIME(t time.Time) uint64 {
	// FILETIME 起点
	const epochDiff = 11644473600 // 秒：1970 - 1601
	sec := uint64(t.Unix() + epochDiff)
	ns := uint64(t.Nanosecond())
	return sec*10000000 + ns/100
}
