package members

import (
	"makedotcsh/ldap"

	"github.com/gin-gonic/gin"
)

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
