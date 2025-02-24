package startupcfg

import "fmt"

// ServiceAPI 服务接口抽象
type ServiceAPI interface {
	// DomainName 接口域名
	DomainName() string
	// Url 接口Url
	Url(apiName string) string
	// AuthData 接口其他数据（鉴权数据等）
	AuthData(key string) (string, error)
}

// ServiceApiConfig 服务接口
type ServiceApiConfig struct {
	Domain string               `json:"domain" yaml:"domain"`
	Auth   map[string]Encrypted `json:"auth" yaml:"auth"`
	Urls   map[string]string    `json:"urls" yaml:"urls"`
}

// DomainName 接口域名
func (c *ServiceApiConfig) DomainName() string {
	return c.Domain
}

// Url 接口Url
func (c *ServiceApiConfig) Url(name string) string {
	if c.Urls != nil {
		if url, ok := c.Urls[name]; ok {
			return url
		}
	}
	return ""
}

// AuthData 接口其他数据（鉴权数据等）
func (c *ServiceApiConfig) AuthData(key string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("auth data %s empty", key)
	}
	if c.Auth != nil {
		if valueEncrypt, ok := c.Auth[key]; ok {
			return valueEncrypt.Get()
		}
	}
	return "", nil
}
