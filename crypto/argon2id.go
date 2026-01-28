package crypto

import "github.com/alexedwards/argon2id"

func Argon2IDHash(pw string) (string, error) {
	return argon2id.CreateHash(pw, argon2id.DefaultParams)
}

func Argon2IDValidateHash(hash string) error {
	_, _, _, err := argon2id.DecodeHash(hash)
	return err
}

func Argon2IDCompare(pw, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pw, hash)
}
