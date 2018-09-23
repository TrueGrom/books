package comment

import (
	"books-backend/app/common"
	"github.com/gin-gonic/gin"
	"time"
)

type CommentModelValidator struct {
	Text         string       `json:"text" binding:"required,min=3,max=8192"`
	BookId       uint         `json:"book_id" binding:"required"`
	CommentModel CommentModel `json:"-"`
}

func (self *CommentModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.CommentModel.BookId = self.BookId
	self.CommentModel.Text = self.Text
	self.CommentModel.CreatedAt = time.Now()
	return nil
}

func NewCommentModelValidator() CommentModelValidator {
	commentModelValidator := CommentModelValidator{}
	return commentModelValidator
}

type DeleteCommentRequestValidator struct {
	CommentId uint `json:"comment_id" binding:"required"`
}

func (self *DeleteCommentRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewDeleteCommentRequestValidator() DeleteCommentRequestValidator {
	deleteCommentRequestValidator := DeleteCommentRequestValidator{}
	return deleteCommentRequestValidator
}
