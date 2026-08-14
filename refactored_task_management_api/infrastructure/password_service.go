package infrastructure

import (
	"refactored_task_management_api/domain"

	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

func NewPasswordService() domain.PasswordService {
	return &passwordService{}
}

func (p *passwordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *passwordService) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
