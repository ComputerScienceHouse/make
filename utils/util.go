package utils

import (
	"errors"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

func GetCSHAuth(c *gin.Context) (csh_auth.UserInfo, error) {
	authAny, exists := c.Get("cshauth")
	if !exists {
		err := errors.New("cshauth does not exist in context")
		return csh_auth.UserInfo{}, err
	}

	auth, ok := authAny.(*csh_auth.Claims)
	if !ok {
		err := errors.New("authentication data in gin context does not match structure")
		return csh_auth.UserInfo{}, err
	}

	user := csh_auth.UserInfo{
		Uuid:     auth.Uuid,
		Email:    auth.Email,
		Username: auth.Username,
		FullName: auth.FullName,
		Groups:   auth.Groups,
	}

	return user, nil
}
