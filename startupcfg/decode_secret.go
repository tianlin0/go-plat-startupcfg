package startupcfg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/forgoer/openssl"
	"text/template"
)

// FormatWithSecretMap 组合解密和模板替换操作，还原密码配置
func FormatWithSecretMap(cfgTemplate string, encryptedMap map[string]string, encryptionKey string) (string, error) {
	if encryptedMap == nil || len(encryptedMap) == 0 {
		return cfgTemplate, nil
	}

	decodeSecretMap := make(map[string]string)
	for key, encryptedValue := range encryptedMap {
		decodedValue, err := hex.DecodeString(encryptedValue)
		if err != nil {
			return "", fmt.Errorf("failed to decode hex string for key %s: %w", key, err)
		}
		dst, err := openssl.AesECBDecrypt(decodedValue, []byte(encryptionKey), openssl.PKCS7_PADDING)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt value for key %s: %w", key, err)
		}
		decodeSecretMap[key] = string(dst)
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

// EncryptSecretMap 对密钥映射进行加密
func EncryptSecretMap(secretMap map[string]string, encryptionKey string) (map[string]string, error) {
	encryptedMap := make(map[string]string)
	var retErr error
	for key, value := range secretMap {
		encryptedBytes, err := openssl.AesECBEncrypt([]byte(value), []byte(encryptionKey), openssl.PKCS7_PADDING)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt value for key %s: %w", key, err)
		}
		encryptedMap[key] = hex.EncodeToString(encryptedBytes)
	}
	return encryptedMap, retErr
}
