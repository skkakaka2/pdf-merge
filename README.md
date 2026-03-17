# PDF 两页合并工具

将PDF文件的每两页合并成一页（上下排列）。

## 安装

```bash
go build -o pdf-merger .
```

## 使用方法

### 处理单个文件
```bash
./pdf-merger -input /path/to/file.pdf -output /path/to/output
```

### 批量处理目录
```bash
./pdf-merger -input /path/to/pdf/folder -output /path/to/output
```

## 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| -input | 输入PDF文件或目录路径 | 当前目录 |
| -output | 输出目录路径 | ./output |
| -workers | 并发处理数 | 4 |

## 示例

```bash
# 处理当前目录下所有PDF文件
./pdf-merger

# 指定输入输出目录
./pdf-merger -input ./pdfs -output ./merged
```

## 输出

合并后的文件名格式：`原文件名_merged.pdf`