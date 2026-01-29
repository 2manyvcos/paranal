package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(value string) string {
	hasher := md5.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func compareToMD5Hash(value, hash string) bool {
	return MD5Hash(value) == hash
}
