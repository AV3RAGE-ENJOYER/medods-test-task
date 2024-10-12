package middlewares

import (
	"errors"
	"log"
	"medods_test_task/service"
	"medods_test_task/tokens"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(db service.UserRepository, tc tokens.TokenController) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenString := c.Request.Header.Get("Authorization")

		if accessTokenString == "" {
			c.String(401, "No authorization header provided")
			c.Abort()
			return
		}

		_, err := tc.Parse(accessTokenString)

		defer func() {
			if err != nil {
				log.Printf(" [Error] %s\n", err)
			}
		}()

		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			c.String(400, "Token is malformed")
			c.Abort()
			return
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			c.String(401, "Invalid token signature")
			c.Abort()
			return
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			c.String(401, "Token has expired")
			c.Abort()
			return
		}

		c.Next()
	}
}
