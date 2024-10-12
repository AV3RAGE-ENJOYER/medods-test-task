package main

import (
	"context"
	"log"
	"medods_test_task/database"
	"medods_test_task/handlers"
	"medods_test_task/middlewares"
	"medods_test_task/service"
	"medods_test_task/tokens"
	"os"
	"time"

	docs "medods_test_task/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           MEDODS Golang test task
// @version         1.0
// @description     This is a test task for Juniour Go Developer in MEDODS.

// @contact.name   Andrei Dombrovskii
// @contact.email  andrushathegames@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey AccessToken
// @in header
// @name Authorization

// @host      127.0.0.1:8080
// @BasePath  /api/v1
func main() {
	godotenv.Load("config.env")

	GIN_MODE := os.Getenv("GIN_MODE")
	GIN_ADDR := os.Getenv("GIN_ADDR")

	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

	tokenController, err := tokens.NewTokenController([]byte(JWT_SECRET_KEY), 15*time.Minute, 72*time.Hour)

	if err != nil {
		log.Fatalf(" [Error] Failed to initialize jwt token controller. %s", err)
	}

	POSTGRES_URL := os.Getenv("POSTGRES_URL")

	db, err := database.NewDB(context.Background(), POSTGRES_URL)

	if err != nil {
		log.Fatalf(" [Error] Failed to establish a connection to PostgresDB. %s", err)
	}

	defer db.Pool.Close()

	userService := service.NewUserService(&db)
	// Setup migrations

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf(" [Error] %s\n", err)
	}

	driver := stdlib.OpenDBFromPool(db.Pool)

	log.Println(" [Info] Migrating database")

	if err := goose.Up(driver, "migrations"); err != nil {
		log.Fatalf(" [Error] %s\n", err)
	}

	if err := driver.Close(); err != nil {
		log.Fatalf(" [Error] %s", err)
	}

	// Setup Gin

	gin.SetMode(GIN_MODE)
	r := setupRouter(userService, tokenController)
	r.Run(GIN_ADDR)
}

func setupRouter(userService *service.UserService, tc tokens.TokenController) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		user.Use(middlewares.AuthMiddleware(userService.Repo, tc))
		user.GET("/ping", handlers.PingHandler())
		user.POST("/add", handlers.AddUserHandler(userService.Repo))

		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.LoginHandler(userService.Repo, tc))
			auth.POST("/refresh", handlers.RefreshTokensHandler(userService.Repo, tc))
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
