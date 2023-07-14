package ienv

// 避免项目中使用相对路径，造成文件不存在的问题
import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

var rootDir string
var confDir string
var dataDir string
var logDir string

func RootDir() string {
	return rootDir
}

func ConfDir() string {
	return confDir
}
func DataDir() string {
	return dataDir
}

func LogDir() string {
	return logDir
}

// init方法不需要调用，程序执行前会自动执行，在main函数之前运行
func init() {
	rootDir = autoDetect()
	confDir = path.Join(rootDir, "conf")
	dataDir = path.Join(rootDir, "data")
	logDir = path.Join(rootDir, "log")
}

var errNotFound = errors.New("cannot found")

// findDirMatch 在指定目录下，向其父目录查找对应的文件是否存在
// 若存在，则返回匹配到的路径
func findDirMatch(baseDir string, fileNames []string) (dir string, err error) {
	currentDir := baseDir
	for i := 0; i < 20; i++ {
		for _, fileName := range fileNames {
			depsPath := filepath.Join(currentDir, fileName)
			if _, err1 := os.Stat(depsPath); !os.IsNotExist(err1) {
				return currentDir, nil
			}
		}

		currentDir = filepath.Dir(currentDir)

		if currentDir == "." {
			break
		}
	}
	return "", errNotFound
}

func autoDetect() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	names := []string{
		"go.mod",
		filepath.Join("conf"),
	}
	dir, err1 := findDirMatch(wd, names)
	if err1 == nil {
		return dir
	}
	return wd
}
