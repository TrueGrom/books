package user

import (
	"books-backend/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UserSighup)
	router.POST("/login", UserLogin)
}

func UserSighup(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	if err := SaveOne(&userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)
	c.JSON(http.StatusCreated, gin.H{"user": "suc"})
}

func UserLogin(c *gin.Context) {

}
