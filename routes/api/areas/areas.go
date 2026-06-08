package areas

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAllAreas(c *gin.Context) {
	areas, err := database.Helper.GetAllAreas()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, areas)
}

func getArea(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	area, err := database.Helper.GetArea(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "area not found"})
			return
		}

		c.Error(err)
		return
	}

	c.JSON(200, area)
}

func createArea(c *gin.Context) {
	var req models.CreateAreaRequest

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.CreateArea(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func Routes(route *gin.RouterGroup) {
	areas := route.Group("/areas")
	areas.GET("/", getAllAreas)
	areas.GET("/:id", getArea)

	areas.POST("/create/", createArea)
}
