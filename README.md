# PDF Merge Universal

一个轻量、原生的 Windows PDF 合并工具。支持拖拽添加、自然排序、手动调整顺序、CSV 导入导出，以及资源管理器右键菜单。

![截图](/ScreenShot.png)

## 功能

- 拖入多个 PDF 或整个文件夹
- 按文件名自然排序、反转顺序、上移和下移
- 自动统计 PDF 数量和页数
- 自动生成公共前缀文件名
- 导出/导入 CSV 顺序表
- 合并结果默认保存到首个 PDF 所在目录
- 支持跨目录添加文件
- 支持安装 Windows PDF 右键菜单
- 单实例运行，连续右键添加文件时合并到同一窗口
- 跟随 Windows 系统字体和显示缩放

## 下载与使用

在仓库右侧的 **Releases** 下载最新版本，解压后运行 `PdfMergeUniversal.exe`。

本工具使用 [qpdf](https://github.com/qpdf/qpdf) 完成 PDF 合并。qpdf 未包含在本仓库中，可采用以下任一方式配置：

1. 将 `qpdf.exe` 放在程序同目录；
2. 将 qpdf 的 `bin` 目录加入系统 `PATH`；
3. 点击右上角 `···` → `配置 qpdf...`，手动选择 `qpdf.exe`。

添加文件并调整顺序后，点击“合并 PDF”或按 `Enter` 即可。

## 快捷键

| 快捷键 | 功能 |
| --- | --- |
| `Enter` / `Ctrl+Enter` | 开始合并 |
| `Ctrl+O` | 添加 PDF |
| `Ctrl+A` | 全选列表 |
| `Delete` | 删除选中项 |
| `Alt+↑` / `Alt+↓` | 上移 / 下移 |

## 从源码构建

需要 Go 1.25 或更高版本：

```bash
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build \
  -buildvcs=false -trimpath \
  -ldflags="-H=windowsgui -s -w" \
  -o PdfMergeUniversal.exe .
```

仓库内的 `rsrc.syso` 包含程序图标资源。生成新的 Windows 图标资源时，可使用 [`rsrc`](https://github.com/akavel/rsrc)：

```bash
rsrc -ico PdfMergeUniversal.ico -o rsrc.syso
```

每次提交都会由 GitHub Actions 自动构建 Windows EXE；推送 `v*` 标签会自动创建 Release。

## 系统要求

- Windows 7 或更高版本
- qpdf

## 许可证

本项目采用 [MIT License](LICENSE)。qpdf 是独立项目，遵循其自身许可证。
