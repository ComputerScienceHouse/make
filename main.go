package main

import (
	"fmt"
	"log"
	"makedotcsh/database"
	"makedotcsh/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	}
}

func notFoundHandler(context *gin.Context) {
	context.JSON(404, gin.H{
		"status":  "error",
		"message": fmt.Sprintf("path %s was not found", context.Request.URL.Path),
	})
}

func main() {
	router := gin.New()

	// init db
	database.Init()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cors.Default())
	router.Use(logger.SetLogger())
	router.Use(errorHandler)
	router.NoRoute(notFoundHandler)

	routes.SetRoutes(router)

	log.Println("running")
	router.Run()
}
