package mocks

import (
	"context"
	"refactored_task_management_api/domain"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(
	ctx context.Context,
	user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FetchByEmail(
	ctx context.Context,
	email string) (*domain.User, error) {
	args := m.Called(ctx, email)

	var user *domain.User

	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

func (m *UserRepositoryMock) Count(
	ctx context.Context,
) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
