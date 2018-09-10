package user

import (
	"books-backend/common"
	"github.com/gin-gonic/gin"
)

type UserModelValidator struct {
	Username  string    `json:"username" binding:"min=4,max=255"`
	Email     string    `json:"email" binding:"exists,email"`
	Password  string    `json:"password" binding:"exists,min=8,max=255"`
	Bio       string    `json:"bio" binding:"max=1024"`
	Image     string    `json:"image" binding:"omitempty,url"`
	userModel UserModel `json:"-"`
}

func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.userModel.Username = self.Username
	self.userModel.Email = self.Email
	self.userModel.Bio = self.Bio

	if self.Password != common.NBRandomPassword {
		self.userModel.setPassword(self.Password)
	}
	if self.Image != "" {
		self.userModel.Image = &self.Image
	}
	return nil
}

func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}

type LoginRequestValidator struct {
	Username string `json:"username" binding:"min=4,max=255"`
	Password string `json:"password" binding:"exists,min=8,max=255"`
}

func (self *LoginRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewLoginRequestValidator() LoginRequestValidator {
	loginRequestValidator := LoginRequestValidator{}
	return loginRequestValidator
}
