package startupcfg_test

import (
	"fmt"
	"github.com/tianlin0/go-plat-startupcfg/startupcfg"
	"testing"
)

func TestEncodeSecretMap(t *testing.T) {
	newMap, err := startupcfg.EncryptSecretMap(map[string]string{
		"aaa": "kkk",
	}, "aa")
	fmt.Println(newMap, err)
}
func TestFormatWithSecretMap(t *testing.T) {
	newMap, err := startupcfg.FormatWithSecretMap("aaaaa{{.aaa}}", map[string]string{
		"aaa": "635d662e90ba51615e811efed2f8be98",
	}, "aa")
	fmt.Println(newMap, err)
}
