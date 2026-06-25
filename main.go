package main

import (
	"embed"
	"io/fs"
	"log"
	"makedotcsh/database"
	"makedotcsh/routes"
	"net/http"
	"os"
	"strings"

	_ "makedotcsh/docs"

	cshauth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed web/dist/*
var staticFS embed.FS
var distFS fs.FS
var assetsFS fs.FS

func errorHandler(c *gin.Context) {
	c.Next()

	if c.Writer.Written() {
		return
	}

	err := c.Errors.Last()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Internal server error",
		})
	}
}

func serveIndex(c *gin.Context) {
	data, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(200, "text/html; charset=utf-8", data)
}

// @title		makedotcsh API
// @host		localhost:8080
// @BasePath	/api/
func main() {
	godotenv.Load()
	router := gin.New()

	host := os.Getenv("MAKE_HOST")

	var err error
	// init embed fs
	distFS, err = fs.Sub(staticFS, "web/dist")
	if err != nil {
		panic(err)
	}

	assetsFS, err = fs.Sub(distFS, "assets")
	if err != nil {
		panic(err)
	}

	// init db
	database.Init()
	defer database.DB.Close()

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
	if os.Getenv("DEV") == "true" {
		router.NoRoute(auth.CookieMiddleware(), createViteProxy())
	} else {
		gin.SetMode(gin.ReleaseMode)

		router.StaticFS("/assets", http.FS(assetsFS))

		router.NoRoute(auth.CookieMiddleware(), func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}
			serveIndex(c)
		})
	}

	// swag
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Started!")
	log.Fatal(router.Run())
}
