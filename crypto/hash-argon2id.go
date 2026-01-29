package crypto

import "github.com/alexedwards/argon2id"

func Argon2IDHash(value string) (string, error) {
	return argon2id.CreateHash(value, argon2id.DefaultParams)
}

func ValidateArgon2IDHash(hash string) error {
	_, _, _, err := argon2id.DecodeHash(hash)
	return err
}

func CompareToArgon2IDHash(value, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(value, hash)
}
