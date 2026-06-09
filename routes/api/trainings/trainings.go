package trainings

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/middleware"
	"makedotcsh/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getAllTrainings(c *gin.Context) {
	trainings, err := database.Helper.GetAllTrainings()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, trainings)
}

func getTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	training, err := database.Helper.GetTraining(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "training not found"})
			return
		}

		c.Error(err)
		return
	}

	c.JSON(200, training)
}

func createTraining(c *gin.Context) {
	var req models.CreateTrainingRequest

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.CreateTraining(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func getAreaTrainings(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	trainings, err := database.Helper.GetAreaTrainings(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, trainings)
}

func Routes(route *gin.RouterGroup) {
	areas := route.Group("/trainings")
	areas.GET("/", getAllTrainings)
	areas.GET("/:id", getTraining)

	areas.GET("/area/:id", getAreaTrainings)

	areas.POST("/create/", middleware.RequireGroup("eboard"), createTraining)
}
