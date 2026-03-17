package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func main() {
	inputDir := flag.String("input", ".", "输入PDF文件或目录路径")
	outputDir := flag.String("output", "./output", "输出目录路径")
	workers := flag.Int("workers", 4, "并发处理数")
	flag.Parse()

	if err := run(*inputDir, *outputDir, *workers); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func run(inputPath, outputDir string, workers int) error {
	log.SetDefaultLoggers()

	var pdfFiles []string

	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("无法访问路径 %s: %w", inputPath, err)
	}

	if info.IsDir() {
		files, err := filepath.Glob(filepath.Join(inputPath, "*.pdf"))
		if err != nil {
			return fmt.Errorf("扫描PDF文件失败: %w", err)
		}
		pdfFiles = files
	} else {
		pdfFiles = []string{inputPath}
	}

	if len(pdfFiles) == 0 {
		return fmt.Errorf("未找到PDF文件")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	total := len(pdfFiles)
	var completed int32
	var failed int32

	jobs := make(chan string, len(pdfFiles))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pdfFile := range jobs {
				err := mergePdfPages(pdfFile, outputDir)
				mu.Lock()
				done := atomic.AddInt32(&completed, 1)
				mu.Unlock()

				if err != nil {
					atomic.AddInt32(&failed, 1)
					fmt.Printf("[%d/%d] 失败: %s - %v\n", done, total, filepath.Base(pdfFile), err)
				} else {
					fmt.Printf("[%d/%d] 完成: %s\n", done, total, filepath.Base(pdfFile))
				}
			}
		}()
	}

	for _, pdfFile := range pdfFiles {
		jobs <- pdfFile
	}
	close(jobs)

	wg.Wait()

	fmt.Printf("\n处理完成! 成功: %d, 失败: %d\n", total-int(failed), failed)
	return nil
}

func mergePdfPages(inputPath, outputDir string) error {
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ".pdf")
	outputPath := filepath.Join(outputDir, baseName+"_merged.pdf")

	nup := model.DefaultNUpConfig()
	nup.Grid = &types.Dim{Width: 1, Height: 2}
	nup.Enforce = true

	conf := model.NewDefaultConfiguration()

	if err := api.NUpFile([]string{inputPath}, outputPath, nil, nup, conf); err != nil {
		return fmt.Errorf("合并页面失败: %w", err)
	}

	return nil
}
