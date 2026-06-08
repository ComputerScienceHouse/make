package me

import (
	"errors"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	authAny, exists := c.Get("cshauth")
	if !exists {
		err := errors.New("cshauth does not exist in context")
		c.Error(err)
		return
	}

	auth, ok := authAny.(*csh_auth.Claims)
	if !ok {
		err := errors.New("authentication data in gin context does not match structure")
		c.Error(err)
		return
	}

	user := csh_auth.UserInfo{
		Uuid:     auth.Uuid,
		Email:    auth.Email,
		Username: auth.Username,
		FullName: auth.FullName,
		Groups:   auth.Groups,
	}

	c.JSON(200, user)
}

func Routes(route *gin.RouterGroup) {
	me := route.Group("/me")
	me.GET("/", test)
}
