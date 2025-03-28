package i18n_test

import (
	"fmt"
	"github.com/tianlin0/go-plat-startupcfg/i18n"
	"testing"
	"time"
)

type AAA struct {
	Name string
}

func TestAesCbc(t *testing.T) {

	// 初始化
	_, err := i18n.NewYamlFile("language.yaml", "zh")
	if err != nil {
		return
	}

	mm := i18n.Translate("aaa.bbbb.ccccc", &AAA{
		Name: "mmmm",
	})
	fmt.Println(mm)

	go func() {
		for i := 0; i < 100; i++ {
			mm := i18n.Translate("aaa.bbbb.ccccc", &AAA{
				Name: "mmmm",
			})
			fmt.Println(mm)
			time.Sleep(2 * time.Second)
		}
	}()

	time.Sleep(5 * time.Second)

	mm = i18n.Translate("aaa.bbbb.ccccc", &AAA{
		Name: "mmmm",
	})

	fmt.Println(mm)

	time.Sleep(time.Minute)
}

func TestAesCbcv(t *testing.T) {

	// 初始化
	_, err := i18n.NewYamlFile("language.yaml", "zh")
	if err != nil {
		return
	}

	mm := i18n.Translate("aaa.bbbb.cccccddd", "nihao")

	fmt.Println(mm)
}
