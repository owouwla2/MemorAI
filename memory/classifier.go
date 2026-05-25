package memory

import (
	"sort"
	"strings"
)

// SoftwareCategory 软件分类
type SoftwareCategory string

const (
	CategoryAntivirus   SoftwareCategory = "杀毒安全"
	CategoryBrowser     SoftwareCategory = "浏览器"
	CategoryIM          SoftwareCategory = "通讯工具"
	CategoryDev         SoftwareCategory = "开发工具"
	CategoryGame        SoftwareCategory = "游戏平台"
	CategoryVirtualize  SoftwareCategory = "虚拟化"
	CategorySystem      SoftwareCategory = "系统服务"
	CategoryDriver      SoftwareCategory = "驱动/硬件"
	CategoryOffice      SoftwareCategory = "办公软件"
	CategoryMedia       SoftwareCategory = "影音媒体"
	CategoryOther       SoftwareCategory = "其他"
)

// rule 一条进程识别规则
type rule struct {
	// 显示名（如 "Chrome浏览器"）
	displayName string
	// 分类
	category SoftwareCategory
	// 匹配模式：进程名（小写）匹配任一即命中
	patterns []string
	// 匹配模式是否使用前缀匹配（默认是包含匹配）
	prefixOnly bool
}

// 规则库（顺序很重要：越具体的越靠前）
var rules = []rule{
	// 杀毒
	{"卡巴斯基", CategoryAntivirus, []string{"avp.exe", "avpui.exe", "kavfsslp", "kavfs"}, false},
	{"火绒安全", CategoryAntivirus, []string{"hipstray", "wsctrl", "usysdiag", "hipsmain", "hipsdaemon"}, false},
	{"360安全卫士", CategoryAntivirus, []string{"360safe", "360tray", "zhudongfangyu", "360sd"}, false},
	{"Windows Defender", CategoryAntivirus, []string{"msmpeng.exe", "nissrv.exe", "securityhealthservice", "smartscreen.exe"}, false},
	{"腾讯电脑管家", CategoryAntivirus, []string{"qqpcrtp", "qqpctray", "qqpcmgr"}, false},

	// 浏览器
	{"Google Chrome", CategoryBrowser, []string{"chrome.exe"}, false},
	{"Microsoft Edge", CategoryBrowser, []string{"msedge.exe", "msedgewebview2.exe"}, false},
	{"Firefox", CategoryBrowser, []string{"firefox.exe"}, false},
	{"Brave", CategoryBrowser, []string{"brave.exe"}, false},
	{"Opera", CategoryBrowser, []string{"opera.exe"}, false},
	{"360浏览器", CategoryBrowser, []string{"360se.exe", "360chrome.exe"}, false},
	{"QQ浏览器", CategoryBrowser, []string{"qqbrowser.exe"}, false},

	// 通讯
	{"微信", CategoryIM, []string{"wechat.exe", "wechatappex.exe", "weixin.exe"}, false},
	{"QQ", CategoryIM, []string{"qq.exe", "qqprotect.exe", "qqexternal.exe"}, false},
	{"钉钉", CategoryIM, []string{"dingtalk.exe"}, false},
	{"飞书", CategoryIM, []string{"feishu.exe", "lark.exe"}, false},
	{"Microsoft Teams", CategoryIM, []string{"teams.exe", "ms-teams.exe"}, false},
	{"Telegram", CategoryIM, []string{"telegram.exe"}, false},
	{"Discord", CategoryIM, []string{"discord.exe"}, false},
	{"Slack", CategoryIM, []string{"slack.exe"}, false},

	// 开发工具
	{"VS Code", CategoryDev, []string{"code.exe", "code - insiders.exe"}, false},
	{"JetBrains IDE", CategoryDev, []string{"idea64.exe", "pycharm64.exe", "goland64.exe", "webstorm64.exe", "clion64.exe", "rider64.exe", "phpstorm64.exe"}, false},
	{"Visual Studio", CategoryDev, []string{"devenv.exe"}, false},
	{"Android Studio", CategoryDev, []string{"studio64.exe"}, false},
	{"Sublime Text", CategoryDev, []string{"sublime_text.exe"}, false},
	{"Cursor", CategoryDev, []string{"cursor.exe"}, false},
	{"Node.js", CategoryDev, []string{"node.exe"}, false},
	{"Python", CategoryDev, []string{"python.exe", "python3.exe", "pythonw.exe"}, false},
	{"Java/JVM", CategoryDev, []string{"java.exe", "javaw.exe"}, false},
	{"Git", CategoryDev, []string{"git.exe", "git-bash.exe"}, false},

	// 游戏平台
	{"Steam", CategoryGame, []string{"steam.exe", "steamwebhelper.exe", "steamservice.exe"}, false},
	{"Epic Games", CategoryGame, []string{"epicgameslauncher.exe", "epicwebhelper.exe", "epicgamesservice.exe"}, false},
	{"Battle.net", CategoryGame, []string{"battle.net.exe", "agent.exe"}, false},
	{"Origin/EA", CategoryGame, []string{"origin.exe", "eadesktop.exe", "easteamproxy.exe"}, false},
	{"Ubisoft Connect", CategoryGame, []string{"upc.exe", "uplay.exe"}, false},
	{"GOG Galaxy", CategoryGame, []string{"galaxyclient.exe"}, false},
	{"WeGame", CategoryGame, []string{"wegame.exe", "tenio"}, false},

	// 虚拟化（重点关注）
	{"WSL2/Hyper-V 虚拟机", CategoryVirtualize, []string{"vmmem", "vmmemwsl", "vmwp.exe"}, false},
	{"Docker Desktop", CategoryVirtualize, []string{"docker desktop.exe", "com.docker.backend.exe", "com.docker.service", "dockerd.exe"}, false},
	{"VMware", CategoryVirtualize, []string{"vmware.exe", "vmware-vmx.exe"}, false},
	{"VirtualBox", CategoryVirtualize, []string{"virtualbox.exe", "virtualboxvm.exe"}, false},

	// 办公软件
	{"Microsoft Office", CategoryOffice, []string{"winword.exe", "excel.exe", "powerpnt.exe", "outlook.exe", "onenote.exe"}, false},
	{"WPS Office", CategoryOffice, []string{"wps.exe", "et.exe", "wpp.exe", "wpscloudsvr.exe", "ksomisc.exe"}, false},
	{"Adobe Acrobat", CategoryOffice, []string{"acrobat.exe", "acrord32.exe"}, false},
	{"Notion", CategoryOffice, []string{"notion.exe"}, false},
	{"Obsidian", CategoryOffice, []string{"obsidian.exe"}, false},

	// 影音
	{"网易云音乐", CategoryMedia, []string{"cloudmusic.exe"}, false},
	{"QQ音乐", CategoryMedia, []string{"qqmusic.exe"}, false},
	{"Spotify", CategoryMedia, []string{"spotify.exe"}, false},
	{"PotPlayer", CategoryMedia, []string{"potplayermini64.exe", "potplayer.exe"}, false},
	{"VLC", CategoryMedia, []string{"vlc.exe"}, false},
	{"哔哩哔哩", CategoryMedia, []string{"bilibili.exe"}, false},

	// 系统
	{"Windows系统服务(svchost)", CategorySystem, []string{"svchost.exe"}, false},
	{"系统进程", CategorySystem, []string{"system", "registry", "smss.exe", "csrss.exe", "wininit.exe", "services.exe", "lsass.exe", "winlogon.exe", "fontdrvhost.exe", "dwm.exe", "lsm.exe", "memory compression"}, false},
	{"资源管理器", CategorySystem, []string{"explorer.exe"}, false},
	{"Windows Search", CategorySystem, []string{"searchindexer.exe", "searchprotocolhost.exe", "searchfilterhost.exe", "searchapp.exe"}, false},
	{"任务管理器", CategorySystem, []string{"taskmgr.exe"}, false},
	{"Shell体验", CategorySystem, []string{"shellexperiencehost.exe", "startmenuexperiencehost.exe", "searchhost.exe", "runtimebroker.exe", "applicationframehost.exe", "sihost.exe", "ctfmon.exe", "taskhostw.exe"}, false},
	{"Windows更新/安装", CategorySystem, []string{"trustedinstaller.exe", "tiworker.exe", "wuauclt.exe", "wuauserv.exe", "usoclient.exe"}, false},

	// 驱动/硬件
	{"NVIDIA驱动", CategoryDriver, []string{"nvidia", "nvcontainer.exe", "nvdisplay.container.exe", "nvsphelper64.exe"}, false},
	{"AMD驱动", CategoryDriver, []string{"radeonsoftware.exe", "amdrsserv.exe", "atieclxx.exe"}, false},
	{"Intel驱动", CategoryDriver, []string{"igfxem.exe", "igfxhk.exe", "igfxtray.exe"}, false},
	{"Realtek音频", CategoryDriver, []string{"ravbg64.exe", "rtkauduservice64.exe", "rtkngui64.exe"}, false},
	{"Logitech外设", CategoryDriver, []string{"lghub.exe", "logioptionsplus", "lcore.exe"}, false},
}

// SoftwareGroup 一个软件的合并统计
type SoftwareGroup struct {
	Name        string           `json:"name"`        // 显示名
	Category    SoftwareCategory `json:"category"`    // 分类
	ProcessCount int             `json:"processCount"`
	TotalRSS    uint64           `json:"totalRSS"`
	TotalMB     float64          `json:"totalMB"`
	TotalPct    float64          `json:"totalPct"`    // 占系统总内存百分比
	Processes   []*ProcessInfo   `json:"processes"`   // 该软件下的所有进程
}

// classifyProcess 把一个进程归到一个软件名下
// 返回 (displayName, category, matched)
func classifyProcess(name string) (string, SoftwareCategory, bool) {
	lower := LowerName(name)
	for _, r := range rules {
		for _, pat := range r.patterns {
			if r.prefixOnly {
				if strings.HasPrefix(lower, pat) {
					return r.displayName, r.category, true
				}
			} else {
				if strings.Contains(lower, pat) {
					return r.displayName, r.category, true
				}
			}
		}
	}
	return "", "", false
}

// GroupBySoftware 把进程列表按软件聚合
func GroupBySoftware(procs []*ProcessInfo, totalMem uint64) []*SoftwareGroup {
	groups := make(map[string]*SoftwareGroup)

	for _, p := range procs {
		display, cat, matched := classifyProcess(p.Name)
		if !matched {
			// 未识别的进程，单独一组（用进程名作为显示名）
			display = p.Name
			cat = CategoryOther
		}

		g, ok := groups[display]
		if !ok {
			g = &SoftwareGroup{
				Name:      display,
				Category:  cat,
				Processes: make([]*ProcessInfo, 0, 2),
			}
			groups[display] = g
		}
		g.ProcessCount++
		g.TotalRSS += p.MemoryRSS
		g.Processes = append(g.Processes, p)
	}

	// 转切片
	result := make([]*SoftwareGroup, 0, len(groups))
	for _, g := range groups {
		g.TotalMB = bytesToMB(g.TotalRSS)
		if totalMem > 0 {
			g.TotalPct = float64(g.TotalRSS) / float64(totalMem) * 100
		}
		result = append(result, g)
	}

	// 按总内存降序
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalRSS > result[j].TotalRSS
	})

	return result
}
