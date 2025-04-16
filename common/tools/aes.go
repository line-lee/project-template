package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const (
	key = "nT9jE1dH9aV9oI3x"
	iv  = "kI5tO2fJ6zR2mS9l"
)

// AESEncrypt 使用AES/CBC/PKCS5Padding进行加密
func AESEncrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plainText = pkcs5Padding(plainText, block.BlockSize())
	thisIv := []byte(iv)[:aes.BlockSize] // 初始化向量IV
	mode := cipher.NewCBCEncrypter(block, thisIv)
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt 使用AES/CBC/PKCS5Padding进行解密
func AESDecrypt(cipherText string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	thisIv := []byte(iv)[:aes.BlockSize] // 初始化向量IV，需要与加密时使用的IV相同
	mode := cipher.NewCBCDecrypter(block, thisIv)
	plainText := make([]byte, len(cipherBytes))
	mode.CryptBlocks(plainText, cipherBytes)
	plainText = pkcs5UnPadding(plainText)
	return plainText, nil
}

// pkcs5Padding 填充明文数据，使其长度是AES块大小的整数倍
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pkcs5UnPadding 移除填充的数据
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
