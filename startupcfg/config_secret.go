package startupcfg

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/tianlin0/go-plat-utils/crypto"
	"text/template"
)

// FormatWithSecret 组合解密和模板替换操作，还原密码配置
func FormatWithSecret(cfgTemplate string, encryptedMap map[string]string, encryptionKey string) (string, error) {
	if encryptedMap == nil || len(encryptedMap) == 0 {
		return cfgTemplate, nil
	}

	decodeSecretMap := make(map[string]string)
	for key, encryptedValue := range encryptedMap {
		decodedValue, err := crypto.ConfigDecryptSecret(encryptedValue, encryptionKey)
		if err != nil {
			return "", fmt.Errorf("failed to decode hex string for key %s: %w", key, err)
		}
		decodeSecretMap[key] = decodedValue
	}

	// 计算配置模板的 MD5 哈希值
	hash := md5.New()
	hash.Write([]byte(cfgTemplate))
	sum := hash.Sum(nil)

	var err error
	templateInstance := template.New(fmt.Sprintf("%x", sum))
	templateInstance = templateInstance.Option("missingkey=zero")
	templateInstance, err = templateInstance.Parse(cfgTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var outputBuffer bytes.Buffer
	err = templateInstance.Execute(&outputBuffer, decodeSecretMap)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return outputBuffer.String(), nil
}

// EncryptSecret 对密钥映射进行加密
func EncryptSecret(secretMap map[string]string, encryptionKey string) (map[string]string, error) {
	encryptedMap := make(map[string]string)
	var retErr error
	for key, value := range secretMap {
		encryptedStr, err := crypto.ConfigEncryptSecret(value, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt value for key %s: %w", key, err)
		}
		encryptedMap[key] = encryptedStr
	}
	return encryptedMap, retErr
}
