package main

import (
	"os"
	"path/filepath"
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
	os.Mkdir(filepath.Join(dir, "gotemp-one"), 0o755)
	os.WriteFile(filepath.Join(dir, "gotemp-two"), []byte("x"), 0o644)
	os.WriteFile(filepath.Join(dir, "other"), []byte("x"), 0o644)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range entries {
		if len(e.Name()) >= 8 && e.Name()[:8] == "gotemp-" {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("应找到 2 个 gotemp-* 项, got %d", count)
	}
}
