// go-temp 是一个零依赖的临时文件/目录管理小工具。
// 子命令：
//   new           在系统临时目录创建一个临时文件，打印路径
//   newdir        创建一个临时目录，打印路径
//   clean <目录>   删除指定临时目录及其内容（需确认路径在 os.TempDir 内）
// 例：
//   go run . new
//   go run . newdir
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
	case "clean":
		if len(os.Args) < 3 {
			fmt.Println("用法: go-temp clean <目录>")
			os.Exit(1)
		}
		target := os.Args[2]
		// 安全限制：只允许清理位于系统临时目录下的路径，避免误删。
		tmp := os.TempDir()
		abs, err := filepath.Abs(target)
		if err != nil || !strings.HasPrefix(abs, tmp) {
			fmt.Fprintf(os.Stderr, "拒绝：只能清理 %s 下的目录\n", tmp)
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

func usage() {
	fmt.Println("用法:")
	fmt.Println("  new            创建临时文件，打印路径")
	fmt.Println("  newdir         创建临时目录，打印路径")
	fmt.Println("  clean <目录>    删除临时目录（仅限系统临时目录内）")
}
