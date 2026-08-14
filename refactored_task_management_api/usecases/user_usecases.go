package usecases

import (
	"context"
	"errors"

	"refactored_task_management_api/domain"
)

type userUsecase struct {
	userRepo        domain.UserRepository
	passwordService domain.PasswordService
	jwtService      domain.JWTService
}

func NewUserUsecase(userRepo domain.UserRepository, passwordService domain.PasswordService, jwtService domain.JWTService) domain.UserUsecase {
	return &userUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uu *userUsecase) Register(ctx context.Context, user *domain.User) error {
	existing, err := uu.userRepo.FetchByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return domain.ErrUserExists
	} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	hashed, err := uu.passwordService.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	count, err := uu.userRepo.Count(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		user.Role = domain.RoleAdmin
	} else {
		user.Role = domain.RoleUser
	}

	return uu.userRepo.Create(ctx, user)
}

func (uu *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uu.userRepo.FetchByEmail(ctx, email)
	if err != nil || user == nil {
		return "", domain.ErrInvalidCreds
	}

	if err := uu.passwordService.Compare(user.Password, password); err != nil {
		return "", domain.ErrInvalidCreds
	}

	return uu.jwtService.GenerateToken(user.ID.Hex(), user.Role)
}
