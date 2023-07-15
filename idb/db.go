package idb

// Support MySQL, PostgreSQL, SQlite, SQL Server

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	// 一个不基于cgo的库
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbMySQL    = "mysql"
	dbPostgres = "postgres"
	dbSQLite   = "sqlite"
)

type Config struct {
	DatabaseType          string `desc:"mysql|postgres|sqlite"`
	DatabaseFile          string `desc:"sqlite专用字段"`
	Host                  string
	Port                  int
	User                  string
	Password              string
	DBName                string
	SSLMode               string
	TimeZone              string `desc:"时区，如Local、Asia/Shanghai"`
	Charset               string `desc:"mysql专用字段,如utf8、utf8mb4"`
	ParseTime             bool   `desc:"mysql专用字段"`
	MaxIdleConns          int
	MaxOpenConns          int
	ConnMaxLifetimeMinute uint8
}

func parseConfig(filename string) (config Config) {
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		panic(err)
	}
	return
}

func MustInit(cfgFilePath string) (db *gorm.DB) {
	cfg := parseConfig(cfgFilePath)
	var skip bool
	switch cfg.DatabaseType {
	case dbMySQL:
		db = initMySQL(cfg)
	case dbPostgres:
		db = initPostgres(cfg)
	case dbSQLite:
		db = initSQLite(cfg)
		skip = true
	default:
		panic(fmt.Errorf("invalid database type: %s", cfg.DatabaseType))
	}
	if skip {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(cfg.ConnMaxLifetimeMinute))
	return

}

func initSQLite(cfg Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseFile), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}
func initMySQL(cfg Config) *gorm.DB {
	var parseTime string
	if cfg.ParseTime {
		parseTime = "True"
	} else {
		parseTime = "False"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset, parseTime, cfg.TimeZone)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initPostgres(cfg Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode, cfg.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db

}
