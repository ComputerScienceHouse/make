package routes

import (
	"makedotcsh/routes/api/areas"

	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine) {
	api := router.Group("/api")
	areas.Routes(api)
}
