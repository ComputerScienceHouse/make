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

func addTrainingToArea(c *gin.Context) {
	var req []int
	id, err := strconv.Atoi(c.Param("id"))
	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.AddTrainingsToArea(id, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

type RemoveTrainingRequest struct {
	ID int `json:"id"`
}

func removeTrainingFromArea(c *gin.Context) {
	var req RemoveTrainingRequest
	id, err := strconv.Atoi(c.Param("id"))
	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.RemoveTrainingFromArea(id, req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func deleteTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.DeleteTraining(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func Routes(route *gin.RouterGroup) {
	trainings := route.Group("/trainings")
	trainings.GET("/", getAllTrainings)
	trainings.GET("/:id", getTraining)

	trainings.GET("/area/:id", getAreaTrainings)

	trainings.POST("/area/:id", middleware.RequireGroup("eboard"), addTrainingToArea)
	trainings.POST("/create/", middleware.RequireGroup("eboard"), createTraining)

	trainings.DELETE("/area/:id", middleware.RequireGroup("eboard"), removeTrainingFromArea)
	trainings.DELETE("/:id", middleware.RequireGroup("eboard"), deleteTraining)
}
