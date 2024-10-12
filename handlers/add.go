package handlers

import (
	"log"
	"medods_test_task/models"
	"medods_test_task/service"
	"medods_test_task/utils"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Add user
// @Description Adds user's email and hashed password to the database
// @Tags User
// @Produce json
// @Param body body handlers.UserRequest true "Request body"
// @Success 200 {object} models.User
// @Failure 400 {string} Bad request
// @Failure 500 {string} Internal server error
// @Security AccessToken
// @Router /user/add [post]
func AddUserHandler(db service.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserRequest

		if err := c.BindJSON(&user); err != nil {
			c.String(400, "Bad request")
			return
		}

		if user.Email != "" && user.Password != "" {
			hashedPassword, err := utils.HashBCrypt([]byte(user.Password))

			if err != nil {
				log.Printf(" [Error] %s\n", err)
				c.String(500, "Failed to hash a password")
				return
			}

			user := models.User{
				Email:          user.Email,
				HashedPassword: string(hashedPassword),
			}

			err = db.AddUser(c, user)

			if err != nil {
				log.Printf(" [Error] %s\n", err)
				c.String(500, "Failed to add a user")
				return
			}

			c.IndentedJSON(200, user)
			return
		}

		c.IndentedJSON(400, gin.H{
			"message": "Bad request",
		})
	}
}
