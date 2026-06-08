package routes

import (
	"makedotcsh/routes/api/areas"
	"makedotcsh/routes/api/me"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine, auth csh_auth.Auth) {
	api := router.Group("/api")
	api.Use(auth.CookieMiddleware())
	areas.Routes(api)
	me.Routes(api)
}
