package main

import (
	"fmt"
	"log"
	"makedotcsh/database"
	"makedotcsh/routes"
	"os"

	cshauth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func serveIndex(c *gin.Context) {
	c.File("./web/dist/index.html")
}

func main() {
	godotenv.Load()
	router := gin.New()

	host := os.Getenv("MAKE_HOST")

	fmt.Println(os.Getenv("MAKE_OIDC_ID"))

	// init db
	database.Init()

	// init auth
	auth, err := cshauth.Init(
		os.Getenv("MAKE_OIDC_ID"),
		os.Getenv("MAKE_OIDC_SECRET"),
		host,
		host+"/auth/login",
		host+"/auth/callback",
		[]string{"profile", "email", "groups"},
	)

	if err != nil {
		log.Panicf("Error initializing CSH auth %v", err)
	}

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(cors.Default())
	router.Use(logger.SetLogger())
	router.Use(errorHandler)

	// auth
	router.GET("/auth/login", auth.HandleLogin)       // This endpoint should match the path for loginURL
	router.GET("/auth/callback", auth.HandleCallback) // This endpoint should match the path for callbackURL
	router.GET("/auth/logout", auth.HandleLogout)

	// api
	routes.SetRoutes(router, auth)

	// frontend
	frontend := router.Group("/")
	frontend.Use(auth.CookieMiddleware())

	if os.Getenv("DEV") == "true" {
		router.NoRoute(auth.CookieMiddleware(), createViteProxy())
	} else {
		frontend.Static("/assets", "./web/dist/assets")
		frontend.GET("/", serveIndex)
		frontend.GET("/:path", serveIndex)
		frontend.GET("/:path/*rest", serveIndex)
	}

	log.Println("running")
	router.Run()
}
