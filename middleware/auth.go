package middleware

import (
	"errors"
	"slices"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

func RequireGroup(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ignore exists because other middleware should've handled it
		authAny, _ := c.Get("cshauth")

		auth, ok := authAny.(*csh_auth.Claims)
		if !ok {
			c.AbortWithError(500, errors.New(
				"authentication data in gin context does not match expected structure",
			))
			return
		}

		if !slices.Contains(auth.Groups, group) {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "route require elevated permissions",
			})
			return
		}

		c.Next()
	}
}
