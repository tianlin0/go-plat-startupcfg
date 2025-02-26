package startupcfg

import (
	"fmt"
	"github.com/tianlin0/go-plat-utils/crypto"
)

type (
	Encrypted string // Encrypted 加密串
)

var (
	encryptFunc = func(e string) (Encrypted, error) {
		return Encrypted(e), nil
	} // 默认将字符串转化为Encrypted类型
	hasSetDecryptHandler = false // 是否设置过解密函数，一个程序里只能设置一次，不然所有以前的就都解不开了
	decryptFunc          = func(e Encrypted) (string, error) {
		return string(e), nil
	} // 默认的密码直接返回
)

// setEncryptHandler 设置加密函数,一般不需要设置
func setEncryptHandler(encryptF func(m string) (Encrypted, error)) {
	if encryptF != nil {
		encryptFunc = encryptF
	}
}

// EncryptedPlainStringByCbcDecrypt 将字符串转化为Encrypted类型，可以使用默认的
// EncryptedPlainStringByCbcDecrypt 和 SetStringForCbcDecrypt 是一一对应的
func EncryptedPlainStringByCbcDecrypt(str string, key string) Encrypted {
	str, err := crypto.CbcEncrypt(str, key, new(crypto.HexCoder))
	if err != nil {
		return ""
	}
	return Encrypted(str)
}

// SetStringForCbcDecrypt 给默认的解密方法设置加密key
func SetStringForCbcDecrypt(key string) error {
	if hasSetDecryptHandler {
		return fmt.Errorf("decryptFunc has seted")
	}
	if key == "" {
		return fmt.Errorf("key is empty")
	}
	//可以直接使用默认的加密解密函数
	decryptFunc = func(e Encrypted) (string, error) {
		str, err := crypto.CbcDecrypt(string(e), key, new(crypto.HexCoder))
		if err != nil {
			return "", err
		}
		return str, nil
	}
	hasSetDecryptHandler = true
	return nil
}

// SetDecryptHandler 设置解密函数
func SetDecryptHandler(decryptF func(m Encrypted) (string, error)) error {
	if hasSetDecryptHandler {
		return fmt.Errorf("decryptFunc has seted")
	}
	if decryptF != nil {
		decryptFunc = decryptF
		hasSetDecryptHandler = true
		return nil
	}
	return fmt.Errorf("decryptF is nil")
}

// Get 获取解密串儿
func (e Encrypted) Get() (string, error) {
	if string(e) == "" {
		return "", nil
	}
	if decryptFunc != nil {
		return decryptFunc(e)
	}
	return "", fmt.Errorf("no set defaultDecrypt")
}

//// MarshalJSON 实现json Marshaler接口 自定义json 编码
//func (e *Encrypted) MarshalJSON() ([]byte, error) {
//	decrypted, err := e.Get()
//	if err != nil {
//		return nil, err
//	}
//	return ([]byte)(fmt.Sprintf("\"%s\"", decrypted)), nil
//}
//
//// UnmarshalJSON 实现Unmarshaler接口 自定义json解码
//func (e *Encrypted) UnmarshalJSON(data []byte) error {
//	var raw json.RawMessage = data
//	var kindType string
//	if err := json.Unmarshal(raw, &kindType); err != nil {
//		return err
//	}
//
//	var str string
//	if err := json.Unmarshal(raw, &str); err != nil {
//		return fmt.Errorf("UnmarshalJSON Encrypted(%s) to string failed: %w ", data, err)
//	}
//	if str == "" {
//		return nil
//	}
//	//str 有可能是已加密过的字符串，或者是解密后的字符串，这里需要区分开
//
//	dd, err := encryptFunc(str)
//	if err != nil {
//		return fmt.Errorf("UnmarshalJSON Encrypted(%s) failed: %w ", str, err)
//	}
//	*e = dd
//	return nil
//}
//
//// MarshalYAML 实现yaml Marshaler接口
//func (e *Encrypted) MarshalYAML() (interface{}, error) {
//	decrypted, err := e.Get()
//	if err != nil {
//		return nil, err
//	}
//	return decrypted, nil
//}
//
//// UnmarshalYAML 实现Unmarshaler接口 自定义yaml解码
//func (e *Encrypted) UnmarshalYAML(value *yaml.Node) error {
//	if value.Value == "" {
//		return nil
//	}
//	dd, err := encryptFunc(value.Value)
//	if err != nil {
//		return fmt.Errorf("UnmarshalYAML Encrypted(%s) failed: %w ", value.Value, err)
//	}
//	*e = dd
//	return nil
//}
