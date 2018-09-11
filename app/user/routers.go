package user

import (
	"books-backend/app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UserSighup)
	router.POST("/login", LoginUser)
}

func UserSighup(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	if err := SaveOne(&userModelValidator.userModel); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewError("database", err), nil)
		return
	}
	common.RenderResponse(c, http.StatusCreated, nil, nil)
}

func LoginUser(c *gin.Context) {
	userModelValidator := NewLoginRequestValidator()
	if err := userModelValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	user, err := FindOneUser(&UserModel{Username: userModelValidator.Username})
	if err != nil {
		common.RenderResponse(c, http.StatusBadRequest, common.CommonError{gin.H{"errors": "user not found"}}, nil)
		return
	}
	token := common.GenToken(user.ID)
	common.RenderResponse(c, http.StatusOK, nil, gin.H{"token": token})
}
