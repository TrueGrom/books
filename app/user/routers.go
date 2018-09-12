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

func UsersModify(router *gin.RouterGroup) {
	router.POST("/reset_password/:username", LoginReset)
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
	userLoginRequestValidator := NewLoginRequestValidator()
	if err := userLoginRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	user, err := FindOneUser(&UserModel{Username: userLoginRequestValidator.Username})
	if err != nil {
		common.RenderResponse(c,
			http.StatusForbidden,
			common.CommonError{Errors: gin.H{"errors": "Invalid username or password"}},
			nil)
		return
	}
	if err = user.checkPassword(userLoginRequestValidator.Password); err != nil {
		common.RenderResponse(c,
			http.StatusForbidden,
			common.CommonError{Errors: gin.H{"errors": "Invalid username or password"}},
			nil)
		return
	}
	token := common.GenToken(user.ID)
	common.RenderResponse(c, http.StatusOK, nil, gin.H{"token": token})
}

func LoginReset(c *gin.Context) {
	username := c.Param("username")
	loginResetRequestValidator := NewLoginResetRequestValidator()
	if err := loginResetRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	user, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		common.RenderResponse(c,
			http.StatusBadRequest,
			common.CommonError{Errors: gin.H{"errors": "Invalid username or password"}},
			nil)
		return
	}
	if err = user.checkPassword(loginResetRequestValidator.Password); err != nil {
		common.RenderResponse(c,
			http.StatusForbidden,
			common.CommonError{Errors: gin.H{"errors": "Invalid username or password"}},
			nil)
		return
	}
	user.setPassword(loginResetRequestValidator.NewPassword)
	user.Update(user)
	common.RenderResponse(c, http.StatusOK, nil, nil)
}
