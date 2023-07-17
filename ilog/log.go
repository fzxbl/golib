package ilog

import (
	iconf "github.com/fzxbl/golib/iconf"
	ienv "github.com/fzxbl/golib/ienv"

	"github.com/natefinch/lumberjack"
	"golang.org/x/exp/slog"
)

type Config struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`

	AddSource bool
}

// MustInitFromFile slog.Logger未实现Close方法
func MustInitFromFile(confPath string) (logger *slog.Logger) {
	cfg := mustLoadConfig(confPath)
	// 设置日志切割选项
	lj := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,    // 单个日志文件的最大尺寸（MB）
		MaxBackups: cfg.MaxBackups, // 保留的日志文件数量
		MaxAge:     cfg.MaxAge,     // 日志文件（天）的最大生命周期
		LocalTime:  cfg.LocalTime,
	}
	logger = slog.New(slog.NewJSONHandler(lj, &slog.HandlerOptions{AddSource: cfg.AddSource}))
	return
}

// InitFromConfig slog.Logger未实现Close方法
func InitFromConfig(cfg Config) (logger *slog.Logger) {
	// 设置日志切割选项
	lj := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,    // 单个日志文件的最大尺寸（MB）
		MaxBackups: cfg.MaxBackups, // 保留的日志文件数量
		MaxAge:     cfg.MaxAge,     // 日志文件（天）的最大生命周期
	}
	logger = slog.New(slog.NewJSONHandler(lj, &slog.HandlerOptions{AddSource: cfg.AddSource}))
	return
}

// 加载配置文件
func mustLoadConfig(filePath string) (v Config) {
	iconf.MustParseToml(filePath, &v)
	v.Filename = ienv.EnvExpand(v.Filename)
	return
}
