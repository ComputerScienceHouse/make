package areas

import (
	"database/sql"
	"errors"
	"makedotcsh/database"
	"makedotcsh/models"

	"github.com/gin-gonic/gin"
)

func getAllAreas(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description FROM areas")
	if err != nil {
		c.Error(err)
		return
	}

	defer rows.Close()

	var areas []models.Area

	for rows.Next() {
		var area models.Area

		err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.Description,
		)

		if err != nil {
			c.Error(err)
			return
		}

		areas = append(areas, area)
	}

	c.JSON(200, areas)
}

func getArea(c *gin.Context) {
	row := database.DB.QueryRow("SELECT id, name, description FROM areas WHERE id = ?", c.Param("id"))

	var area models.Area

	err := row.Scan(
		&area.ID,
		&area.Name,
		&area.Description,
	)

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

	_, err = database.DB.Exec(
		"INSERT INTO areas (name, description) VALUES (?, ?)",
		req.Name,
		req.Description,
	)
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
