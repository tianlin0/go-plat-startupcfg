package startupcfg_test

import (
	"fmt"
	"github.com/tianlin0/go-plat-startupcfg/startupcfg"
	"testing"
)

type AAA struct {
	Name string
}

func TestAesCbc(t *testing.T) {

	_, err := startupcfg.NewI18nFile("language.yaml", "zh")
	if err != nil {
		return
	}
	mm := startupcfg.I18nTranslate("aaa.bbbb.ccccc", &AAA{
		Name: "mmmm",
	})

	fmt.Println(mm)
}
