package iconf

import (
	"github.com/BurntSushi/toml"
)

func MustParseToml(filePath string, v any) {
	_, err := toml.DecodeFile(filePath, v)
	if err != nil {
		panic(err)
	}
}
