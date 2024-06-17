package fake

func NewFakeCryptography() *FakeCryptography {
	return &FakeCryptography{}
}

type FakeCryptography struct{}

func (h *FakeCryptography) GenerateHash(value string) (string, error) {
	hash := value + "-hashed"
	return hash, nil
}

func (h *FakeCryptography) ValidateHash(password, passwordHash string) bool {
	hashedValue := password + "-hashed"
	return hashedValue == passwordHash
}
