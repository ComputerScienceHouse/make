package routes

import (
	"makedotcsh/routes/api/areas"
	"makedotcsh/routes/api/me"
	"makedotcsh/routes/api/trainings"
	"makedotcsh/routes/api/user"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

// for browser -> header auth api endpoint,
// moves the JWT from cookie to request header
func cookieToAuthHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("Auth")
		if err == nil && cookie != "" {
			c.Request.Header.Set("Authorization", "Bearer "+cookie)
		}

		c.Next()
	}
}

func SetRoutes(router *gin.Engine, auth csh_auth.Auth) {
	api := router.Group("/api")
	api.Use(cookieToAuthHeader())
	api.Use(auth.HeaderMiddleware())

	areas.Routes(api)
	me.Routes(api)
	trainings.Routes(api)
	user.Routes(api)

}
