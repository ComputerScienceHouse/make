package user

import (
	"makedotcsh/database"
	"makedotcsh/middleware"
	"makedotcsh/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getUserTrainings(c *gin.Context) {
	uuid := c.Param("uuid")

	userTrainings, err := database.Helper.GetUserTrainings(uuid)
	if err != nil {
		c.Error(err)
		return
	}

	trainings := []models.Training{}

	for _, id := range userTrainings {
		training, err := database.Helper.GetTraining(id)
		if err != nil {
			c.Error(err)
			return
		}

		trainings = append(trainings, training)
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

func deleteUserTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	uuid := c.Param("uuid")
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.DeleteTrainingFromUser(id, uuid)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func Routes(route *gin.RouterGroup) {
	areas := route.Group("/user")
	areas.GET("/:uuid/trainings", getUserTrainings)
	areas.GET("/:uuid/areaAccess", getUserAreaAccess)

	areas.DELETE("/:uuid/trainings/:id", middleware.RequireGroup("eboard"), deleteUserTraining)

	areas.POST("/:uuid/trainings/", middleware.RequireGroup("eboard"), createUserTraining)
}
