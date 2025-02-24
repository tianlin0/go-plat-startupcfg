package startupcfg_test

import (
	"fmt"
	"github.com/tianlin0/go-plat-utils/conf/startupcfg"
	"github.com/tianlin0/go-plat-utils/conn"
	"github.com/tianlin0/go-plat-utils/conv"
	"github.com/tianlin0/go-plat-utils/crypto"
	"github.com/tianlin0/go-plat-utils/startconf"
	"testing"
)

// TemplateURL 模板访问地址
type TemplateURL struct {
	TemplateGetCdBiz            string //获取CD业务信息
	TemplateBindCdBizTag        string //绑定部署标签
	TemplateAppserverSubmitById string //gdp_appserver_go 提交template执行的方法
	GetPodListByCdId            string
	InsertCI                    string
}

// GdpConfig gdp全局配置
type GdpConfig struct {
	HostAndPort              *conn.Connect
	MysqlConnect             *startupcfg.MysqlConfig
	MysqlConnectODP          *startupcfg.MysqlConfig
	RedisConnect             *startupcfg.RedisConfig
	GdpExternalOrigin        string
	ClientSecret             string
	TemplateIdBatchDeleteCd  string //批量删除部署的模板ID
	TemplateIdCopyCdWithCdId string
	TemplateIdCopyCd         string

	DefaultRTXLoginToken string //rtxLoginToken

	DefaultSystemRoleNameMap map[string][]string
}

func TestGetAllApiUrlMap(t *testing.T) {

	keyStr := "tianlin020250214"

	startupcfg.SetDecryptHandler(func(e startupcfg.Encrypted) (string, error) {
		str, err := crypto.CbcDecrypt(string(e), keyStr, new(crypto.HexCoder))
		if err != nil {
			return "", err
		}
		return str, nil
	})

	one, _ := startconf.NewStartupForYamlFile("all_start_cfg_test.yaml")
	mapTemp := one.AllApiUrlMap()
	tempUrl := new(TemplateURL)
	conv.Unmarshal(mapTemp, tempUrl)

	fmt.Println(conv.String(tempUrl))

	cMap, _ := one.AllCustomMap()

	tempCMap := new(GdpConfig)
	conv.Unmarshal(cMap, tempCMap)

	ccTemp, _ := one.AllMysqlMap()
	conv.Unmarshal(ccTemp, tempCMap)

	ccTemp2, _ := one.AllRedisMap()

	conv.Unmarshal(ccTemp2, tempCMap)

	fmt.Println(conv.String(tempCMap))

}
