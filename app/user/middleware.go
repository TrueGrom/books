package user

import (
	"books-backend/app/common"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_token := c.Request.Header.Get("Authorization")
		if auth_token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Authorization header not provided"})
			c.Abort()
			return
		}
		token_header := strings.Split(auth_token, " ")
		if len(token_header) != 2 {
			c.JSON(http.StatusForbidden, gin.H{"message": "Usage: JWT <token_header>"})
			c.Abort()
			return
		}
		if token_header[0] != "JWT" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Authorization header need to start with 'JWT'"})
			c.Abort()
			return
		}
		user, err := TokenParse(token_header[1])
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Token is invalid"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func TokenParse(token string) (UserModel, error) {
	tokenParse, err := jwt.Parse(token, common.KeyFunc)
	if err != nil {
		return UserModel{}, errors.New("Token is invalid")
	}
	claims, ok := tokenParse.Claims.(jwt.MapClaims)
	if !ok {
		return UserModel{}, errors.New("Token is invalid")
	}
	user, err := FindOneUser(&UserModel{ID: uint(claims["id"].(float64))})
	return user, err
}
