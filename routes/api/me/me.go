package me

import (
	"makedotcsh/utils"

	"github.com/gin-gonic/gin"
)

// bullshit type that matches csh.auth.UserInfo because swag can't resolve dependencies
type MeResponse struct {
	Uuid     string   `json:"uuid"`
	Email    string   `json:"email"`
	Username string   `json:"preferred_username"`
	FullName string   `json:"name"`
	Groups   []string `json:"groups"`
}

// GetAuthUser godoc
//
// @Summary Get authenticated user
// @Description Returns information about the currently authenticated user
// @Tags me
// @Produce json
// @Success 200 {object} MeResponse
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
