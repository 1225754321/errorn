package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// 定义命令行参数
	filePath      string
	directoryPath string
	anyDirFiles   bool
)

func init() {
	// 设置flag读取命令入参
	// 两个场景, 一个是直接运行二进制，一个是通过go generate调用
	// 优先考虑直接运行二进制的场景
	// 直接运行
	// 一个-f参数用于指定文件路径, 只调整这个文件路径下的Wrap调用
	// 一个-d参数用于指定文件夹路径, 递归调整这个文件夹下的每个go文件的Wrap调用
	// 两个优先级-f大于-d
	flag.StringVar(&filePath, "f", "", "用于指定文件路径, 只调整这个文件路径下的Wrap调用")
	flag.StringVar(&directoryPath, "d", "", "用于指定文件夹路径, 递归调整这个文件夹下的每个go文件的Wrap调用")
	// go generate场景
	// -a 参数用于指定当前模式, 为false则是只调整当前文件, 为true则是递归调整当前文件对应目录下的每个go文件
	flag.BoolVar(&anyDirFiles, "a", false, "默认为false, 用于指定当前模式, 为false则是只调整当前文件, 为true则是递归调整当前文件对应目录下的每个go文件")
}

func main() {
	flag.Parse()

	// 确定处理路径
	var targetPath string
	goFile := os.Getenv("GOFILE")

	if filePath != "" {
		targetPath = filePath
	} else if directoryPath != "" {
		targetPath = directoryPath
	} else if goFile != "" { // go generate 场景
		targetPath = goFile
		if anyDirFiles {
			targetPath = path.Dir(goFile) // 如果是go generate模式且只处理当前文件，则使用当前文件所在目录
		}
	} else {
		fmt.Println("Usage: go run main.go -f <file_path> | -d <directory_path>")
		fmt.Println("       or set GOFILE environment variable for go generate")
		os.Exit(1)
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		panic(fmt.Errorf("path error: %v", err))
	}

	if info.IsDir() {
		// 处理目录
		err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				processFile(path)
			}
			return nil
		})
		if err != nil {
			panic(fmt.Errorf("directory walk error: %v", err))
		}
	} else {
		// 处理单个文件
		processFile(targetPath)
	}
}
