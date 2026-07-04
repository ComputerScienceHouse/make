package me

import (
	"makedotcsh/utils"

	"github.com/gin-gonic/gin"
)

// GetAuthUser godoc
//
// @Summary Get authenticated user
// @Description Returns information about the currently authenticated user
// @Tags me
// @Produce json
// @Success 200 {object} csh_auth.UserInfo
// @Failure 500 {object} models.ErrorResponse
// @Router /me [get]
func getAuthUser(c *gin.Context) {
	user, err := utils.GetCSHAuth(c)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user)
}

func Routes(route *gin.RouterGroup) {
	me := route.Group("/me")
	me.GET("/", getAuthUser)
}
