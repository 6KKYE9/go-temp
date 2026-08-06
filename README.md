# go-temp

零依赖的临时文件/目录管理小工具（基于 `os.CreateTemp` / `os.MkdirTemp`）。

## 功能

- `new`：创建临时文件，打印路径
- `newdir`：创建临时目录，打印路径
- `clean <目录>`：删除临时目录（仅限系统临时目录内，防止误删）

## 用法

```bash
go run . new
go run . newdir
go run . clean "C:\Users\ASUS\AppData\Local\Temp\gotemp-abc123"
```
