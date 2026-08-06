# go-temp

临时文件和目录管理，基于标准库的 `os.CreateTemp` / `os.MkdirTemp`。

```bash
go run . new                       # 建临时文件，打印路径
go run . newdir                    # 建临时目录，打印路径
go run . list                      # 列出临时目录里的 gotemp-* 项
go run . clean "C:\Users\ASUS\AppData\Local\Temp\gotemp-abc123"
```

`clean` 只能删系统临时目录里的内容，避免误删别的地方（用路径分隔符严格判定，`/tmpfoo` 不会误匹配 `/tmp`）。
