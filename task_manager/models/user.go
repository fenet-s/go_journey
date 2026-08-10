package models

import "strings"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password,omitempty" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

func NormalizeRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleAdmin:
		return RoleAdmin
	case RoleUser:
		return RoleUser
	default:
		return RoleUser
	}
}
