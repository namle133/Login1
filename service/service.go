package service

import (
	"context"

	"github.com/namle133/Login1.git/LOGIN/domain"
)

type IUser interface {
	CreateUser(c context.Context, u *domain.UserInit) error
	SignIn(c context.Context, u *domain.UserInit) (*domain.Claims, error)
	UserAdmin() error
	CheckUserAdmin(c context.Context, token string) error
	LogOut(c context.Context, token string) error
}
