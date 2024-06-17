package adapter

type CryptographyContract interface {
	GenerateHash(password string) (string, error)
	ValidateHash(password, passwordHash string) error
}
