package idb

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Title string
	Code  string
	Price uint
}

func Test_MustInit(t *testing.T) {
	db, close := MustInit("./testdata/sqlite.toml")
	// db := MustInit("./testdata/postgres.toml")
	db.AutoMigrate(&Product{})
	// 插入内容
	db.Create(&Product{Title: "新款手机", Code: "D42", Price: 1000})
	db.Create(&Product{Title: "新款电脑", Code: "D43", Price: 3500})

	// 读取内容
	var product Product
	db.First(&product, 1) // find product with integer primary key
	fmt.Println(product)
	db.First(&product, "code = ?", "D42") // find product with code D42
	fmt.Println(product)
	err := close.Close()
	fmt.Println(err)
}
