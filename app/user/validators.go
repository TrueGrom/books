package user

import (
	"books-backend/app/common"
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

	self.userModel.setPassword(self.Password)

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

type LoginResetRequestValidator struct {
	Password    string `json:"password" binding:"exists,min=8,max=255"`
	NewPassword string `json:"new_password" binding:"exists,min=8,max=255"`
}

func (self *LoginResetRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewLoginResetRequestValidator() LoginResetRequestValidator {
	loginResetRequestValidator := LoginResetRequestValidator{}
	return loginResetRequestValidator
}

type ForgetPasswordRequestValidator struct {
	Username string `json:"username" binding:"min=4,max=255"`
}

func (self *ForgetPasswordRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewForgetPasswordRequestValidator() ForgetPasswordRequestValidator {
	forgetPasswordRequestValidator := ForgetPasswordRequestValidator{}
	return forgetPasswordRequestValidator
}

type GetNewPasswordRequestValidator struct {
	Password string `json:"password" binding:"min=4,max=255"`
	Token    string `json:"token" binding:"required"`
}

func (self *GetNewPasswordRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewGetNewPasswordRequestValidator() GetNewPasswordRequestValidator {
	getNewPasswordRequestValidator := GetNewPasswordRequestValidator{}
	return getNewPasswordRequestValidator
}

type AddBookToUserRequestValidator struct {
	BookId []uint `json:"book_id" binding:"required,dive,min=1"`
}

func (self *AddBookToUserRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewAddBookToUserRequestValidator() AddBookToUserRequestValidator {
	addBookToUserRequestValidator := AddBookToUserRequestValidator{}
	return addBookToUserRequestValidator
}

type BookRating struct {
	Book_id uint `json:"book_id"`
	Rating  int8 `json:"rating"`
}

type AddRatingToBooksRequestValidator struct {
	Books []BookRating `json:"books" binding:"required"`
}

func (self *AddRatingToBooksRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewAddRatingToBooksRequestValidator() AddRatingToBooksRequestValidator {
	addRatingToBooksRequestValidator := AddRatingToBooksRequestValidator{}
	return addRatingToBooksRequestValidator
}
