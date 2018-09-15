package book

import (
	"books-backend/app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BooksRegister(router *gin.RouterGroup) {
	router.GET("/search", SearchBooks)
}

func SearchBooks(c *gin.Context) {
	searchQueryRequestValidator := NewSearchQueryRequestValidator()
	if err := searchQueryRequestValidator.Bind(c); err != nil {
		common.RenderResponse(c, http.StatusUnprocessableEntity, common.NewValidatorError(err), nil)
		return
	}
	books, err := FindBooksByTitle(searchQueryRequestValidator.Q, 20)
	if err != nil {
		common.RenderResponse(c,
			http.StatusBadRequest,
			common.CommonError{Errors: gin.H{"errors": "Invalid search"}},
			nil)
		return
	}
	common.RenderResponse(c, http.StatusOK, nil, books)
}
