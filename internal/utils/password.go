package utils

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintext
	p.Hash = hash

	return nil
}

func (p *Password) Compare(plaintext string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintext))
}
