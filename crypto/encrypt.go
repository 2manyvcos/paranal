package crypto

func Encrypt(key, value string) (string, error) {
	return encryptAES(key, value)
}

func Decrypt(key, value string) (string, error) {
	return decryptAES(key, value)
}
