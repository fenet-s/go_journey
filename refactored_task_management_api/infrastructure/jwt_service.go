package infrastructure

import (
	"errors"
	"refactored_task_management_api/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret     string
	expiration time.Duration
}

func NewJWTService(secret string, expiration time.Duration) domain.JWTService {
	if expiration == 0 {
		expiration = 24 * time.Hour
	}
	return &jwtService{
		secret:     secret,
		expiration: expiration,
	}
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *jwtService) GenerateToken(userID, role string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) ValidateToken(tokenString string) (*domain.Claims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(j.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &domain.Claims{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
