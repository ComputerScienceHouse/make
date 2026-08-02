package resources

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/middleware"
	"makedotcsh/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RemoveResourceRequest struct {
	ID int `json:"id"`
}

// GetAllTrainings godoc
//
// @Summary Get all trainings
// @Description Returns all trainings
// @Tags trainings
// @Produce json
// @Success 200 {array} models.Training
// @Failure 500 {object} models.ErrorResponse
// @Router /trainings [get]
func getAllResources(c *gin.Context) {
	rs, err := database.Helper.GetAllResources()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, rs)
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
func getResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	r, err := database.Helper.GetResource(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, gin.H{"error": "resource not found"})
			return
		}

		c.Error(err)
		return
	}

	c.JSON(200, r)
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
func createResource(c *gin.Context) {
	var req models.Resource

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.CreateResource(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func updateResource(c *gin.Context) {
	var req models.Resource

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.UpdateResource(req, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func deleteResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.DeleteResource(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func getAreaResources(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	resources, err := database.Helper.GetAreaResources(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resources)
}

func addResourceToArea(c *gin.Context) {
	var req []int
	id, err := strconv.Atoi(c.Param("id"))
	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.AddResourceToArea(id, req)
	if err != nil {
		c.Error(err)
		return
	}
}

func removeResourceFromArea(c *gin.Context) {
	var req RemoveResourceRequest
	id, err := strconv.Atoi(c.Param("id"))
	err = c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.RemoveResourceFromArea(id, req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func Routes(route *gin.RouterGroup) {
	resources := route.Group("/resources")
	resources.GET("/", getAllResources)
	resources.GET("/:id", getResource)
	resources.GET("/area/:id", getAreaResources)

	resources.PUT("/:id", middleware.RequireGroup("eboard"), updateResource)

	resources.POST("/", middleware.RequireGroup("eboard"), createResource)
	resources.POST("/area/:id", middleware.RequireGroup("eboard"), addResourceToArea)

	resources.DELETE("/area/:id", middleware.RequireGroup("eboard"), removeResourceFromArea)
	resources.DELETE("/:id", middleware.RequireGroup("eboard"), deleteResource)
}
