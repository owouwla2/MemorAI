package main

import (
	"flag"
	"fmt"
	"os"

	"memory-analyzer/startup"
)

func runStartup(args []string) {
	if len(args) == 0 {
		printStartupUsage()
		return
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "list", "ls":
		runStartupList(rest)
	case "enable", "on":
		runStartupSet(rest, true)
	case "disable", "off":
		runStartupSet(rest, false)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n\n", sub)
		printStartupUsage()
	}
}

func runStartupList(args []string) {
	fs := flag.NewFlagSet("startup list", flag.ExitOnError)
	allCols := fs.Bool("verbose", false, "显示完整命令行")
	fs.Parse(args)

	items, err := startup.EnumerateStartupItems()
	if err != nil {
		fatal("枚举失败: %v", err)
	}
	isAdmin := startup.IsAdmin()

	fmt.Println()
	if isAdmin {
		fmt.Println(paint("✓ Running as Administrator (all items modifiable)", cGreen))
	} else {
		fmt.Println(paint("⚠ Not running as Administrator (system-level items locked)", cYellow))
	}
	fmt.Println(paint(hr(70), cGray))

	enabled := 0
	disabled := 0
	for _, it := range items {
		if it.Enabled {
			enabled++
		} else {
			disabled++
		}
	}
	fmt.Printf("  Total: %d · ", len(items))
	fmt.Printf("%s · %s\n",
		paint(fmt.Sprintf("Enabled: %d", enabled), cGreen),
		paint(fmt.Sprintf("Disabled: %d", disabled), cDim),
	)
	fmt.Println()

	header := formatRow(
		[]string{"ID", "STATUS", "NAME", "SCOPE", "SOURCE"},
		[]int{18, 9, 30, 7, 10},
	)
	fmt.Println(paint(header, cBold))
	fmt.Println(paint(hr(80), cGray))

	for _, it := range items {
		status := paint("ENABLED ", cGreen)
		if !it.Enabled {
			status = paint("disabled", cDim)
		}
		scope := it.Scope
		if it.Scope == startup.ScopeSystem {
			scope = paint(it.Scope, cYellow)
		}
		row := formatRow(
			[]string{
				it.ID,
				status,
				truncate(it.Name, 30),
				scope,
				it.Source,
			},
			[]int{18, 9, 30, 7, 10},
		)
		fmt.Println(row)
		if *allCols {
			fmt.Printf("    %s\n", paint(truncate(it.Command, 76), cDim))
		}
	}
	fmt.Println()
	fmt.Println(paint("Tip:", cDim), "memorai-cli startup disable <ID>   /   startup enable <ID>")
	fmt.Println()
}

func runStartupSet(args []string, enable bool) {
	if len(args) == 0 {
		fatal("缺少 ID。运行: memorai-cli startup list")
	}
	id := args[0]
	if err := startup.SetEnabled(id, enable); err != nil {
		fatal("操作失败: %v", err)
	}
	action := "enabled"
	color := cGreen
	if !enable {
		action = "disabled"
		color = cYellow
	}
	fmt.Printf(paint("✓ %s startup item %s\n", color), action, id)
}

func printStartupUsage() {
	fmt.Println(`startup - 自启项管理

USAGE:
    memorai-cli startup <subcommand> [args]

SUBCOMMANDS:
    list                   列出所有自启项
    list --verbose         同上，显示完整命令行
    enable <id>            启用某项
    disable <id>           禁用某项

EXAMPLES:
    memorai-cli startup list
    memorai-cli startup disable a1b2c3d4e5f6...
    memorai-cli startup enable  a1b2c3d4e5f6...

NOTE:
    系统级 (HKLM) 自启项需要以管理员身份运行。`)
}
