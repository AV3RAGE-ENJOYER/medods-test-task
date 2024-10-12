package handlers

import (
	"log"
	"medods_test_task/service"
	"medods_test_task/tokens"
	"medods_test_task/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// @BasePath /api/v1

// @Summary Login user
// @Description Login user via email and password
// @Tags Authentication
// @Produce json
// @Accept json
// @Param body body handlers.UserRequest true "User email and password"
// @Success 200 {object} tokens.TokenPair
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/login [post]
func LoginHandler(db service.UserRepository, tc tokens.TokenController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userRequest UserRequest

		if err := c.BindJSON(&userRequest); err != nil {
			log.Printf(" [Error] %s\n", err)
			c.IndentedJSON(500, "Failed to unmarshall json")
			return
		}

		user, err := db.GetUser(c, userRequest.Email)

		if err == pgx.ErrNoRows {
			log.Printf(" [Error] %s\n", err)
			c.IndentedJSON(404, "User not found")
			return
		}

		err = user.CheckPassword(userRequest.Password)

		if err != nil {
			log.Printf(" [Error] %s\n", err)
			c.IndentedJSON(401, "Unauthorized")
			return
		}

		tokenPair, err := tc.NewJWT(user.Email, c.ClientIP())

		if err != nil {
			log.Printf(" [Error] %s\n", err)
			c.IndentedJSON(500, "Failed to generate token pair")
			return
		}

		refreshTokenHash, err := utils.HashBCrypt([]byte(tokenPair.RefreshToken))

		if err != nil {
			log.Printf(" [Error] %s\n", err)
			c.IndentedJSON(500, "Failed to hash refresh token")
			return
		}

		err = db.AddRefreshToken(
			c,
			user.Email,
			string(refreshTokenHash),
			tc.RefreshTokenTTL,
			c.ClientIP(),
		)

		if err != nil {
			log.Printf(" [Error] %s\n", err)
			c.IndentedJSON(500, "Failed to add refresh token")
			return
		}

		c.IndentedJSON(200, tokenPair)
	}
}
