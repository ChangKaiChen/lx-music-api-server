package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	EAPIKEY = []byte("e82ckenh8dichen8")
)

func EApiEncrypt(url, text string) map[string]string {
	digest := MD5(fmt.Sprintf("nobody%suse%smd5forencrypt", url, text))
	data := fmt.Sprintf("%s-36cd479b6b5-%s-36cd479b6b5-%s", url, text, digest)
	return map[string]string{"params": AesEncrypt([]byte(data), EAPIKEY)}
}
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func AesEncrypt(text, key []byte) string {
	blockSize := 16
	padding := blockSize - len(text)%blockSize
	padText := append(text, bytes.Repeat([]byte{byte(padding)}, padding)...)
	block, _ := aes.NewCipher(key)
	var ciphertext []byte
	ciphertext = make([]byte, len(padText))
	for i := 0; i < len(padText); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], padText[i:i+blockSize])
	}
	return strings.ToUpper(hex.EncodeToString(ciphertext))
}
