package mocks

import "github.com/stretchr/testify/mock"

type PasswordServiceMock struct {
	mock.Mock
}

func (p *PasswordServiceMock) Hash(password string) (string, error) {
	args := p.Called(password)

	return args.String(0), args.Error(1)

}

func (p *PasswordServiceMock) Compare(hashedPassword string, password string) error {
	args := p.Called(hashedPassword, password)
	return args.Error(0)
}
