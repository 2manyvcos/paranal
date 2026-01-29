package crypto

func Hash(value string) (string, error) {
	return Argon2IDHash(value)
}

func ValidateHash(hash string) error {
	return ValidateArgon2IDHash(hash)
}

func CompareToHash(value, hash string) (bool, error) {
	return CompareToArgon2IDHash(value, hash)
}
