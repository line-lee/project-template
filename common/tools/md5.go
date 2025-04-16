package tools

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

func MD5X(source string) (string, error) {
	hash := md5.New()
	_, err := io.WriteString(hash, source)
	if err != nil {
		return "", err
	}
	b := hash.Sum(nil)
	// 16进制
	return fmt.Sprintf("%x", b), nil
}

// SHA加密[1, 256, 512]
func SHAX(data []byte, shaType crypto.Hash, isHex bool) (sh []byte) {

	var hashs hash.Hash = nil
	switch shaType {
	case crypto.SHA1:
		hashs = sha1.New()
	case crypto.SHA256:
		hashs = sha256.New()
	case crypto.SHA3_512:
		hashs = sha512.New()
	default:
		return nil
	}

	hashs.Write(data)
	hashed := hashs.Sum(nil)

	if isHex {
		return []byte(hex.EncodeToString(hashed))
	}
	return hashed
}
