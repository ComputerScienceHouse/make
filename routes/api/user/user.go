package user

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/middleware"
	"makedotcsh/models"

	"github.com/gin-gonic/gin"
)

func getUserTrainings(c *gin.Context) {
	uuid := c.Param("uuid")
	trainings, err := database.Helper.GetUserTrainings(uuid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "training not found"})
			return
		}

		c.Error(err)
		return
	}

	c.JSON(200, trainings)
}

func createUserTraining(c *gin.Context) {
	var req models.UserTraining

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.CreateUserTraining(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

// TODO: can this be made more efficient?
func getUserAreaAccess(c *gin.Context) {
	uuid := c.Param("uuid")
	areas, err := database.Helper.GetAllAreasWithUserAccess(uuid)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, areas)
}

func Routes(route *gin.RouterGroup) {
	areas := route.Group("/user")
	areas.GET("/:uuid/trainings", getUserTrainings)
	areas.GET("/:uuid/areaAccess", getUserAreaAccess)

	areas.POST("/:uuid/trainings/", middleware.RequireGroup("eboard"), createUserTraining)
}
