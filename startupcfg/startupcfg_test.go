package startupcfg_test

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tianlin0/go-plat-startupcfg/startupcfg"
	"github.com/tianlin0/go-plat-utils/conv"
	"github.com/tianlin0/go-plat-utils/crypto"
	"github.com/tianlin0/go-plat-utils/utils"
	"gopkg.in/yaml.v3"
	"testing"
)

const (
	encKey = "____tianlin0____"
)

func init1() {
	startupcfg.SetDecryptHandler(func(m startupcfg.Encrypted) (string, error) {
		mysqlPwd, _ := crypto.CbcDecrypt(string(m), encKey, new(crypto.HexCoder))
		return mysqlPwd, nil
	})

}

func TestStartConfig(t *testing.T) {
	init1()

	conf, err := startupcfg.NewByYamlFile("config_test.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	mysqlPwd, _ := crypto.CbcEncrypt("root", encKey, new(crypto.HexCoder))
	tCRPullCommConn, _ := crypto.CbcEncrypt("datamore@2019", encKey, new(crypto.HexCoder))

	fmt.Println(mysqlPwd, tCRPullCommConn)

	testCases := []*utils.TestStruct{
		{"api service-discovery domain", []any{"service-discovery"}, []any{"http://www.baidu.com"}, func(paasName string) string {
			return conf.ServiceAPI(paasName).DomainName()
		}},
		{"api service-discovery url TemplateCopyInClusterService", []any{"service-discovery", "TemplateCopyInClusterService"}, []any{"/v1/project/{{.projectName}}/get"}, func(paasName string, urlKey string) string {
			return conf.ServiceAPI(paasName).Url(urlKey)
		}},
		{"mysql MysqlConnect address", []any{"MysqlConnect"}, []any{"127.0.0.1:3306"}, func(mysqlName string) string {
			return conf.Mysql(mysqlName).ServerAddress()
		}},
		{"mysql MysqlConnect username", []any{"MysqlConnect"}, []any{"root"}, func(mysqlName string) string {
			return conf.Mysql(mysqlName).User()
		}},
		{"mysql MysqlConnect database", []any{"MysqlConnect"}, []any{"db_gdp_server"}, func(mysqlName string) string {
			return conv.String(conf.Mysql(mysqlName).DatabaseName())
		}},
		{"custom normal AppId", []any{"AppId"}, []any{"gdp-appserver-go"}, func(mysqlName string) string {
			str := conf.CustomNormal(mysqlName)
			return conv.String(str)
		}},
		{"custom normal CloseHttpMemCache", []any{"CloseHttpMemCache"}, []any{false}, func(mysqlName string) bool {
			str := conf.CustomNormal(mysqlName)
			b, _ := conv.Bool(str)
			return b
		}},
		{"mysql MysqlConnect pwd", []any{"MysqlConnect"}, []any{"root"}, func(mysqlName string) string {
			return conf.Mysql(mysqlName).Password()
		}},
		{"custom sensitive tCRPullCommConn", []any{"tCRPullCommConn"}, []any{"datamore@2019"}, func(mysqlName string) string {
			str, _ := conf.CustomSensitive(mysqlName)
			return str
		}},
	}
	utils.TestFunction(t, testCases, nil)

	t.Log(conv.String(conf.All()))
}

func TestEncrypt(t *testing.T) {
	mysqlPwd, _ := crypto.CbcEncrypt("root", encKey)
	mysqlCode, _ := crypto.CbcDecrypt(mysqlPwd, encKey)

	t.Log(mysqlPwd)
	t.Log(mysqlCode)
}
func TestDecrypt(t *testing.T) {
	mysqlCode, _ := crypto.CbcDecrypt("a792a5b35ac43cc0132a093ff06b395b5412eed5f9ee71f48e1b62ef048052d5", encKey)

	t.Log(mysqlCode)
}

func TestEncryptedMarshal(t *testing.T) {
	str := startupcfg.Encrypted("d278cee958a4a2521245a1e80224c3b31deedbce24e7eb3c33137a6ab4ca99e5")
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	js, _ := json.Marshal(&str)
	t.Log(string(js)) // "root"

	yml, err := yaml.Marshal(&str)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(yml)) // root

	mysqlTemp := &startupcfg.MysqlConfig{
		PasswordEncoded: "d278cee958a4a2521245a1e80224c3b31deedbce24e7eb3c33137a6ab4ca99e5",
	}
	js, _ = json.Marshal(mysqlTemp)
	t.Log(string(js))
}

func TestEncryptedUnMarshal(t *testing.T) {
	init1()

	jsonStr := `{"username":"","pwEncoded":"root","address":"","database":"","charset":""}`

	mysqlTemp := &startupcfg.MysqlConfig{}
	_ = json.Unmarshal([]byte(jsonStr), mysqlTemp)
	t.Log(string(mysqlTemp.PasswordEncoded))
}
