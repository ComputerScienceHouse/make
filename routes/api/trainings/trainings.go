package trainings

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/middleware"
	"makedotcsh/models"
	"makedotcsh/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllTrainings godoc
//
// @Summary Get all trainings
// @Description Returns all trainings
// @Tags trainings
// @Produce json
// @Success 200 {array} models.Training
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings [get]
func getAllTrainings(c *gin.Context) {
	trainings, err := database.Helper.GetAllTrainings()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, trainings)
}

// GetTraining godoc
//
// @Summary Get training
// @Description Returns a training by ID
// @Tags trainings
// @Produce json
// @Param id path int true "Training ID"
// @Success 200 {object} models.Training
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/{id} [get]
func getTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	training, err := database.Helper.GetTraining(id, false)

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

// GetTraining godoc
//
// @Summary Get training
// @Description Returns a training by ID
// @Tags trainings
// @Produce json
// @Param id path int true "Training ID"
// @Success 200 {object} models.Training
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/{id}/full [get]
func getTrainingFull(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	training, err := database.Helper.GetTraining(id, true)

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

// CreateTraining godoc
//
// @Summary Create training
// @Description Create a new training
// @Tags trainings
// @Accept json
// @Param training body models.CreateTrainingRequest true "Training"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings [post]
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

// GetAreaTrainings godoc
//
// @Summary Get area trainings
// @Description Returns all trainings required for an area
// @Tags trainings
// @Produce json
// @Param id path int true "Area ID"
// @Success 200 {array} models.Training
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/area/{id} [get]
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

// AddTrainingToArea godoc
//
// @Summary Add trainings to area
// @Description Associates one or more trainings with an area
// @Tags trainings
// @Accept json
// @Param id path int true "Area ID"
// @Param trainingIds body []int true "Training IDs"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/area/{id} [post]
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

// RemoveTrainingFromArea godoc
//
// @Summary Remove training from area
// @Description Removes a training requirement from an area
// @Tags trainings
// @Accept json
// @Param id path int true "Area ID"
// @Param request body trainings.RemoveTrainingRequest true "Training removal request"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/area/{id} [delete]
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

// DeleteTraining godoc
//
// @Summary Delete training
// @Description Deletes a training
// @Tags trainings
// @Param id path int true "Training ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/{id} [delete]
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

// UpdateTraining godoc
//
// @Summary Update training
// @Description Updates an existing training
// @Tags trainings
// @Accept json
// @Param id path int true "Training ID"
// @Param training body models.CreateTrainingRequest true "Updated training"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/{id} [put]
func updateTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var req models.CreateTrainingRequest

	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.UpdateTraining(req, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

// submitCompletedTraining godoc
//
// @Summary Submits and grades a completed training
// @Description Submits and grades a completed training
// @Tags trainings
// @Accept json
// @Param training body models.Submission true "Training"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings/{id}/submissions [post]
func submitCompletedTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var req models.Submission

	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := utils.GetCSHAuth(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := utils.GradeTraining(id, user.Uuid, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, res)
}

func Routes(route *gin.RouterGroup) {
	trainings := route.Group("/trainings")
	trainings.GET("/", getAllTrainings)
	trainings.GET("/:id", getTraining)
	trainings.GET("/:id/full", middleware.RequireGroup("eboard"), getTrainingFull)
	trainings.GET("/area/:id", getAreaTrainings)

	trainings.PUT("/:id", middleware.RequireGroup("eboard"), updateTraining)

	trainings.POST("/area/:id", middleware.RequireGroup("eboard"), addTrainingToArea)
	trainings.POST("/", middleware.RequireGroup("eboard"), createTraining)
	trainings.POST("/:id/submissions", submitCompletedTraining)

	trainings.DELETE("/area/:id", middleware.RequireGroup("eboard"), removeTrainingFromArea)
	trainings.DELETE("/:id", middleware.RequireGroup("eboard"), deleteTraining)
}
