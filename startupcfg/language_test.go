package startupcfg_test

import (
	"fmt"
	"github.com/tianlin0/go-plat-startupcfg/startupcfg"
	"testing"
	"time"
)

type AAA struct {
	Name string
}

func TestAesCbc(t *testing.T) {

	// 初始化
	_, err := startupcfg.NewI18nFile("language.yaml", "zh")
	if err != nil {
		return
	}

	mm := startupcfg.I18nTranslate("zh", "aaa.bbbb.ccccc", &AAA{
		Name: "mmmm",
	})
	fmt.Println(mm)

	go func() {
		for i := 0; i < 100; i++ {
			mm := startupcfg.I18nTranslate("en", "aaa.bbbb.ccccc", &AAA{
				Name: "mmmm",
			})
			fmt.Println(mm)
			time.Sleep(2 * time.Second)
		}
	}()

	time.Sleep(5 * time.Second)

	mm = startupcfg.I18nTranslate("zh", "aaa.bbbb.ccccc", &AAA{
		Name: "mmmm",
	})

	fmt.Println(mm)

	time.Sleep(time.Minute)
}
