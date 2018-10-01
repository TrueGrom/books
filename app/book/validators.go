package book

import (
	"books-backend/app/common"
	"github.com/gin-gonic/gin"
)

type SearchQueryRequestValidator struct {
	Q string `form:"q" binding:"min=3,max=5000,required"`
}

func (self *SearchQueryRequestValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	return nil
}

func NewSearchQueryRequestValidator() SearchQueryRequestValidator {
	searchQueryRequestValidator := SearchQueryRequestValidator{}
	return searchQueryRequestValidator
}
