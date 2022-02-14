package service

import (
	"github.com/lemon997/lemonMall/common/authJWT"
)

type LogoutRequest struct {
	authJWT.Claims
}
