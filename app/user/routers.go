package user

import (
	"books/app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UserSighup)
	router.POST("/login", LoginUser)
	router.POST("/forget_password", ForgetPassword)
	router.POST("/reset_forget_password", GetNewPassword)
	router.POST("/books", JWTAuthorization(), AddBookToUser)
	router.DELETE("/books", JWTAuthorization(), DeleteBookFromUser)
	router.GET("/books", JWTAuthorization(), GetBooks)
	router.POST("/books/rating", JWTAuthorization(), AddRating)
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

// View for reset User password by username
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

func ForgetPassword(c *gin.Context) {
	forgetPasswordRequestValidator := NewForgetPasswordRequestValidator()
	if err := forgetPasswordRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	user, err := FindOneUser(&UserModel{Username: forgetPasswordRequestValidator.Username})
	if err != nil {
		common.RenderResponse(c,
			http.StatusBadRequest,
			common.CommonError{Errors: gin.H{"errors": "Invalid username"}},
			nil)
		return
	}
	SendEmailWithResetLink(user)
	common.RenderResponse(c, http.StatusOK, nil, nil)
}

func GetNewPassword(c *gin.Context) {
	getNewPasswordRequestValidator := NewGetNewPasswordRequestValidator()
	if err := getNewPasswordRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	user, err := TokenParse(getNewPasswordRequestValidator.Token)
	if err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, "Token is invalid", nil)
		return
	}
	user.setPassword(getNewPasswordRequestValidator.Password)
	user.Update(user)
	common.RenderResponse(c, http.StatusOK, nil, nil)
}

func AddBookToUser(c *gin.Context) {
	addBookToUserRequestValidator := NewAddBookToUserRequestValidator()
	if err := addBookToUserRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	userInt, _ := c.Get("user")
	user, _ := userInt.(UserModel)
	user.AddBooksToUser(addBookToUserRequestValidator.BookId)
	common.RenderResponse(c, http.StatusOK, nil, nil)
}

func DeleteBookFromUser(c *gin.Context) {
	addBookToUserRequestValidator := NewAddBookToUserRequestValidator()
	if err := addBookToUserRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	userInt, _ := c.Get("user")
	user, _ := userInt.(UserModel)
	user.DeleteBooksToUser(addBookToUserRequestValidator.BookId)
	common.RenderResponse(c, http.StatusOK, nil, nil)
}

func GetBooks(c *gin.Context) {
	userAuth, _ := c.Get("user")
	user, _ := userAuth.(UserModel)
	books, err := user.GetAllBooksFromUser()
	if err != nil {
		return
	}
	common.RenderResponse(c, http.StatusOK, nil, books)
}

func AddRating(c *gin.Context) {
	addRatingToBooksRequestValidator := NewAddRatingToBooksRequestValidator()
	if err := addRatingToBooksRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	userAuth, _ := c.Get("user")
	user, _ := userAuth.(UserModel)
	err := AddRatingToBook(&user, addRatingToBooksRequestValidator.Books)
	if err != nil {
		return
	}
	common.RenderResponse(c, http.StatusOK, nil, nil)
}
