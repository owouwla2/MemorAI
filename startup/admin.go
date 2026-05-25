package startup

import "golang.org/x/sys/windows"

// IsAdmin 检测当前进程是否为管理员
func IsAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0) // 当前进程 token
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}
