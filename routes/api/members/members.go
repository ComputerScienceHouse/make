package members

import (
	"makedotcsh/ldap"

	"github.com/gin-gonic/gin"
)

// GetActiveMembers godoc
//
// @Summary Get active CSH members
// @Description Returns all members in the LDAP active group
// @Tags members
// @Produce json
// @Success 200 {array} ldap.UserWUUID
// @Failure 500 {object} models.ErrorResponse
// @Router /members/active [get]
func getActiveMembers(c *gin.Context) {
	members, err := ldap.GetGroupMembers("active")
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, members)
}
func Routes(route *gin.RouterGroup) {
	areas := route.Group("/members")
	areas.GET("/active", getActiveMembers)

}
