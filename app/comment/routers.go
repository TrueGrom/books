package comment

import (
	"books-backend/app/common"
	"books-backend/app/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CommentsRegister(router *gin.RouterGroup) {
	router.POST("/books", user.JWTAuthorization(), AddCommentToBook)
	router.DELETE("books", user.JWTAuthorization(), DeleteCommentFromBook)
	router.PATCH("/books")
}

func AddCommentToBook(c *gin.Context) {
	commentModelValidator := NewCommentModelValidator()
	if err := commentModelValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	userAuth, _ := c.Get("user")
	user, _ := userAuth.(user.UserModel)
	commentModelValidator.CommentModel.UserId = user.ID
	err := SaveOne(&commentModelValidator.CommentModel)
	if err != nil {
		common.RenderResponse(c,
			http.StatusBadRequest,
			common.CommonError{Errors: gin.H{"errors": "Invalid data"}},
			nil)
		return
	}
	common.RenderResponse(c, http.StatusOK, nil, nil)
}

func DeleteCommentFromBook(c *gin.Context) {
	deleteCommentRequrestValidator := NewDeleteCommentRequestValidator()
	if err := deleteCommentRequrestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	userAuth, _ := c.Get("user")
	user, _ := userAuth.(user.UserModel)

	num, err := DeleteOneComment(CommentModel{ID: deleteCommentRequrestValidator.CommentId, UserId: user.ID},
		time.Now().Add(time.Duration(-3)*time.Hour))
	if err != nil || num == 0 {
		common.RenderResponse(c,
			http.StatusBadRequest,
			common.CommonError{Errors: gin.H{"errors": "You cannot delete this comment or there is no comment"}},
			nil)
		return
	}
	common.RenderResponse(c, http.StatusOK, nil, nil)
}
