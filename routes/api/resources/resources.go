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

// GetAllResources godoc
//
// @Summary Get all resources
// @Description Returns all resources
// @Tags resources
// @Produce json
// @Success 200 {array} models.Resource
// @Failure 500 {object} models.ErrorResponse
// @Router /resources [get]
func getAllResources(c *gin.Context) {
	rs, err := database.Helper.GetAllResources()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, rs)
}

// GetResource godoc
//
// @Summary Get resource
// @Description Returns a resource by ID
// @Tags resources
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} models.Resource
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources/{id} [get]
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

// CreateResource godoc
//
// @Summary Create resource
// @Description Creates a new resource
// @Tags resources
// @Accept json
// @Param resource body models.Resource true "Resource"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources [post]
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

// UpdateResource godoc
//
// @Summary Update resource
// @Description Updates an existing resource
// @Tags resources
// @Accept json
// @Param id path int true "Resource ID"
// @Param resource body models.Resource true "Resource"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources/{id} [put]
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

// DeleteResource godoc
//
// @Summary Delete resource
// @Description Deletes a resource by ID
// @Tags resources
// @Param id path int true "Resource ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources/{id} [delete]
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

// GetAreaResources godoc
//
// @Summary Get area resources
// @Description Returns all resources associated with an area
// @Tags resources
// @Produce json
// @Param id path int true "Area ID"
// @Success 200 {array} models.Resource
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources/area/{id} [get]
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

// AddResourceToArea godoc
//
// @Summary Add resources to area
// @Description Adds resources to an area
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Area ID"
// @Param request body []int true "Resource IDs"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources/area/{id} [post]
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

// RemoveResourceFromArea godoc
//
// @Summary Remove resource from area
// @Description Removes a resource from an area
// @Tags resources
// @Accept json
// @Param id path int true "Area ID"
// @Param request body resources.RemoveResourceRequest true "Resource removal request"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /resources/area/{id} [delete]
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
