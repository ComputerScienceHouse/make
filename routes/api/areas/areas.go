package areas

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/middleware"
	"makedotcsh/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllAreas godoc
//
// @Summary Get all areas
// @Description Returns all areas
// @Tags areas
// @Produce json
// @Success 200 {array} models.Area
// @Failure 500 {object} models.ErrorResponse
// @Router /areas [get]
func getAllAreas(c *gin.Context) {
	areas, err := database.Helper.GetAllAreas()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, areas)
}

// GetArea godoc
//
// @Summary Get area
// @Description Returns an area by ID
// @Tags areas
// @Produce json
// @Param id path int true "Area ID"
// @Success 200 {object} models.Area
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /areas/{id} [get]
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

// CreateArea godoc
//
// @Summary Create area
// @Description Creates a new area
// @Tags areas
// @Accept json
// @Param area body models.CreateAreaRequest true "Area"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /areas [post]
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

// UpdateArea godoc
//
// @Summary Update area
// @Description Updates an existing area
// @Tags areas
// @Accept json
// @Param id path int true "Area ID"
// @Param area body models.CreateAreaRequest true "Updated area"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /areas/{id} [put]
func updateArea(c *gin.Context) {
	var req models.CreateAreaRequest

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = database.Helper.UpdateArea(req, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

// DeleteArea godoc
//
// @Summary Deletes an area
// @Description Deletes an area
// @Tags areas
// @Accept json
// @Param id path int true "Area ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /areas/{id} [delete]
func deleteArea(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	err = database.Helper.DeleteArea(id)
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

	areas.PUT("/:id", middleware.RequireGroup("eboard"), updateArea)
	areas.POST("/", middleware.RequireGroup("eboard"), createArea)

	areas.DELETE("/:id", middleware.RequireGroup("eboard"), deleteArea)
}
