package initial

import (
	"net/http"
	"regexp"

	"gitee.com/don178/m3u8/global"
	"github.com/forgoer/openssl"
	"go.uber.org/zap"
)

// Decryption handler
type Decryption interface {
	// isNeed 判断是否需要该解密算法
	// 获取解密需要使用的密钥等
	// head 为m3u8文件
	// 如果返回 true 会使用 Decrypt 进行解密
	IsNeed(head []byte) bool
	// decrypt 进行解密, 返回解密后的数据
	Decrypt(dst []byte) ([]byte, error)
}

// DecryptionList 解密算法集合
var DecryptionList = []Decryption{
	&Aes128{},
}

// ----------------------------------------------------------------
// Aes128 解密
type Aes128 struct {
	key []byte
	iv  []byte
}

func (aes *Aes128) IsNeed(head []byte) bool {
	re := regexp.MustCompile(`#EXT-X-KEY:METHOD=(.*?),URI="(.*?)"`)
	submatch := re.FindSubmatch(head)
	if len(submatch) != 3 {
		return false
	}
	if string(submatch[2]) != "AES-128" {
		return false
	}
	// 获取key
	rep, err := http.Get(string(submatch[1]))
	if err != nil {
		global.Log.Fatal("获取key uri 失败", zap.Error(err))
	}
	aes.key = make([]byte, rep.ContentLength)
	n, err := rep.Body.Read(aes.key)
	if int64(n) != rep.ContentLength {
		global.Log.Fatal("key 写入 m.key 失败", zap.Error(err))
	}
	// 获取 iv
	re = regexp.MustCompile(`IV="(.*?)"`)
	submatch = re.FindSubmatch(head)
	if len(submatch) != 2 {
		aes.iv = []byte("0000000000000000")
	} else {
		aes.iv = submatch[1]
	}
	return true
}

func (aes *Aes128) Decrypt(dst []byte) ([]byte, error) {
	dst, err := openssl.AesCBCDecrypt(dst, aes.key, aes.iv, openssl.PKCS7_PADDING)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// ----------------------------------------------------------------
