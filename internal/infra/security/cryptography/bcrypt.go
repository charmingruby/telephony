package cryptography

import "golang.org/x/crypto/bcrypt"

func NewCryptography() *Cryptography {
	return &Cryptography{}
}

type Cryptography struct{}

func (h *Cryptography) GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *Cryptography) ValidateHash(password, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}
