package ifile

import (
	"bufio"
	"os"
	"path/filepath"
)

func CreateFileRecursive(path string) (file *os.File, err error) {
	// 获取文件所在的目录
	dir := filepath.Dir(path)

	// 创建目录
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}

	// 创建文件
	file, err = os.Create(path)
	if err != nil {
		return
	}
	return
}

// ReadLine 按行读入内存，批量返回
func ReadLine(fileName string) (lines []string, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return
}
