// go-temp 是临时文件/目录管理小工具。
// 子命令：
//   new           在系统临时目录创建一个临时文件，打印路径
//   newdir        创建一个临时目录，打印路径
//   list          列出系统临时目录里由本工具创建的 gotemp-* 项
//   clean <目录>   删除指定临时目录及其内容（仅限系统临时目录内）
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "new":
		f, err := os.CreateTemp("", "gotemp-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "创建失败: %v\n", err)
			os.Exit(1)
		}
		f.Close()
		fmt.Println(f.Name())
	case "newdir":
		d, err := os.MkdirTemp("", "gotemp-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "创建失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(d)
	case "list":
		cmdList(os.TempDir())
	case "clean":
		if len(os.Args) < 3 {
			fmt.Println("用法: go-temp clean <目录>")
			os.Exit(1)
		}
		abs, err := safeAbs(os.Args[2], os.TempDir())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := os.RemoveAll(abs); err != nil {
			fmt.Fprintf(os.Stderr, "删除失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("已清理 %s\n", abs)
	default:
		usage()
		os.Exit(1)
	}
}

// safeAbs 把 target 解析为绝对路径，并确认它位于 tmpRoot 内（含 tmpRoot 自身）。
// 用路径分隔符严格判定，避免 "/tmpfoo" 误匹配 "/tmp"。
func safeAbs(target, tmpRoot string) (string, error) {
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf("路径无效: %v", err)
	}
	root := filepath.Clean(tmpRoot)
	if abs == root {
		return abs, nil
	}
	prefix := root + string(os.PathSeparator)
	if !strings.HasPrefix(abs, prefix) {
		return "", fmt.Errorf("拒绝：只能清理 %s 下的目录", root)
	}
	return abs, nil
}

// cmdList 列出临时目录里前缀为 gotemp- 的文件/目录。
func cmdList(tmpRoot string) {
	entries, err := os.ReadDir(tmpRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取临时目录失败: %v\n", err)
		os.Exit(1)
	}
	found := 0
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "gotemp-") {
			fmt.Println(filepath.Join(tmpRoot, e.Name()))
			found++
		}
	}
	if found == 0 {
		fmt.Println("没有由 go-temp 创建的临时项")
	}
}

func usage() {
	fmt.Println("用法:")
	fmt.Println("  new            创建临时文件，打印路径")
	fmt.Println("  newdir         创建临时目录，打印路径")
	fmt.Println("  list           列出临时目录里的 gotemp-* 项")
	fmt.Println("  clean <目录>    删除临时目录（仅限系统临时目录内）")
}
