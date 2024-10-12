package tests

import (
	"encoding/json"
	"medods_test_task/database"
	"medods_test_task/email"
	"medods_test_task/handlers"
	"medods_test_task/middlewares"
	"medods_test_task/models"
	"medods_test_task/service"
	"medods_test_task/tokens"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter(userService *service.UserService, tc tokens.TokenController, es *service.EmailService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		user.Use(middlewares.AuthMiddleware(userService.Repo, tc))
		user.GET("/ping", handlers.PingHandler())
		user.POST("/add", handlers.AddUserHandler(userService.Repo))

		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.LoginHandler(userService.Repo, tc))
			auth.POST("/refresh", handlers.RefreshTokensHandler(userService.Repo, tc, es.Repo))
		}
	}
	return r
}

func sendRequest(user handlers.UserRequest, path string) *http.Request {
	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest(
		"POST",
		path,
		strings.NewReader(string(userJSON)),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

var MockDB database.MockDBUserRepository = database.MockDBUserRepository{
	Users: map[string]models.User{
		"test@gmail.com": {
			Email:          "test@gmail.com",
			HashedPassword: "$2a$12$sJWDAVM8PDIO62Xiz9nSF.kSWVs/ZikVElCqdvmSUiPJL3Pb/ZlVW",
		},
		"admin@gmail.com": {
			Email:          "admin@gmail.com",
			HashedPassword: "$2a$12$kSmFPc0kKZBjmBUJyMRiGu1uHC6QFR43RBXehFYKk4tcsOMeo7PfK",
		},
	},
	RefreshTokens: make(map[string]models.RefreshToken),
}

var MockUserService *service.UserService = service.NewUserService(&MockDB)

var MockTokenController tokens.TokenController = tokens.TokenController{
	SigningKey:      []byte("test-secret-key"),
	AccessTokenTTL:  15 * time.Minute,
	RefreshTokenTTL: 72 * time.Hour,
}

var MockEmailRepository email.MockEmailRepository = email.MockEmailRepository{}
var MockEmailService *service.EmailService = service.NewEmailService(&MockEmailRepository)

var router *gin.Engine = setupRouter(MockUserService, MockTokenController, MockEmailService)
