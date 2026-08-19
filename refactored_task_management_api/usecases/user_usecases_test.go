package usecases

import (
	"context"
	"refactored_task_management_api/domain"
	"refactored_task_management_api/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo        *mocks.UserRepositoryMock
	passwordService *mocks.PasswordServiceMock
	userUsecase     domain.UserUsecase
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.userRepo = new(mocks.UserRepositoryMock)
	s.passwordService = new(mocks.PasswordServiceMock)
	s.userUsecase = NewUserUsecase(
		s.userRepo,
		s.passwordService,
		nil,
	)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (s *UserUsecaseTestSuite) TestRegister_Success() {
	user := &domain.User{
		Email:    "test@gmail.com",
		Password: "password",
	}

	s.userRepo.
		On("FetchByEmail", mock.Anything, user.Email).
		Return(nil, domain.ErrUserNotFound)

	s.passwordService.
		On("Hash", "password").
		Return("hashed-password", nil)

	s.userRepo.
		On("Count", mock.Anything).
		Return(int64(0), nil)

	s.userRepo.
		On("Create", mock.Anything, user).
		Return(nil)

	err := s.userUsecase.Register(
		context.Background(),
		user,
	)

	s.NoError(err)

	s.Equal("hashed-password", user.Password)
	s.Equal(domain.RoleAdmin, user.Role)

	s.userRepo.AssertExpectations(s.T())
	s.passwordService.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestRegister_UserAlreadyExists() {

	user := &domain.User{
		Email:    "test@gmail.com",
		Password: "password",
	}

	s.userRepo.
		On("FetchByEmail", mock.Anything, user.Email).
		Return(user, nil)

	err := s.userUsecase.Register(
		context.Background(),
		user,
	)
	s.ErrorIs(err, domain.ErrUserExists)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestRegister_HashError() {
	user := &domain.User{
		Email:    "fenet",
		Password: "password",
	}
	s.userRepo.
		On("FetchByEmail", mock.Anything, user.Email).
		Return(nil, domain.ErrUserNotFound)

	s.passwordService.
		On("Hash", "password").
		Return("", domain.ErrHashingFailed)
	err := s.userUsecase.Register(
		context.Background(),
		user,
	)
	s.ErrorIs(err, domain.ErrHashingFailed)

	s.userRepo.AssertExpectations(s.T())
	s.passwordService.AssertCalled(s.T(), "Hash", "password")

}

func (s *UserUsecaseTestSuite) TestRegister_CountError() {
	user := &domain.User{
		Email:    "",
		Password: "password",
	}
	s.userRepo.
		On("FetchByEmail", mock.Anything, user.Email).
		Return(nil, domain.ErrUserNotFound)
	s.passwordService.
		On("Hash", "password").
		Return("hashed-password", nil)
	s.userRepo.
		On("Count", mock.Anything).
		Return(int64(0), domain.ErrUserNotFound)

	err := s.userUsecase.Register(
		context.Background(),
		user,
	)
	s.ErrorIs(err, domain.ErrUserNotFound)

	s.userRepo.AssertExpectations(s.T())
	s.passwordService.AssertExpectations(s.T())
}

func (s *UserUsecaseTestSuite) TestRegister_ExistingUser() {
	user := &domain.User{
		Email:    "fenet@gmail.com",
		Password: "password",
	}
	s.userRepo.
		On("FetchByEmail", mock.Anything, user.Email).
		Return(user, nil)

	err := s.userUsecase.Register(
		context.Background(),
		user,
	)
	s.ErrorIs(err, domain.ErrUserExists)

	s.userRepo.AssertExpectations(s.T())
	s.passwordService.AssertNotCalled(s.T(), "Hash", mock.Anything)
}
