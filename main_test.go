package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeAbs(t *testing.T) {
	tmp := t.TempDir()
	// 目录内允许
	inner := filepath.Join(tmp, "gotemp-abc")
	if got, err := safeAbs(inner, tmp); err != nil || got != filepath.Clean(inner) {
		t.Fatalf("应允许目录内路径, got=%q err=%v", got, err)
	}
	// 根本身允许
	if _, err := safeAbs(tmp, tmp); err != nil {
		t.Fatalf("应允许根本身: %v", err)
	}
	// 形如 /tmpfoo 的近似路径要拒绝
	sibling := filepath.Join(filepath.Dir(tmp), filepath.Base(tmp)+"evil")
	if _, err := safeAbs(sibling, tmp); err == nil {
		t.Fatal("近似路径应被拒绝")
	}
	// 完全外部拒绝
	if _, err := safeAbs("/etc", tmp); err == nil {
		t.Fatal("/etc 应被拒绝")
	}
}

func TestCmdList(t *testing.T) {
	dir := t.TempDir()
	// 创建两个 gotemp-* 与一个无关项
	if err := os.Mkdir(filepath.Join(dir, "gotemp-one"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "gotemp-two"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "other"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	// 真正调用 cmdList，并捕获它的 stdout 来校验
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	cmdList(dir)
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	count := 0
	for _, l := range lines {
		if strings.Contains(l, "gotemp-") {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("cmdList 应找到 2 个 gotemp-* 项, got %d\n输出:\n%s", count, buf.String())
	}
}
