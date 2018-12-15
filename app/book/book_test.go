package book

import (
	"books/app/common"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var test_db *gorm.DB

type request struct {
	init           func(*http.Request)
	url            string
	method         string
	bodyData       string
	expectedCode   int
	responseRegexg string
	msg            string
}

func BookModelMocker(n int) []BookModel {
	var offset int
	test_db.Model(&BookModel{}).Count(&offset)
	var ret []BookModel
	for i := offset + 1; i <= offset+n; i++ {
		bookModel := BookModel{
			Authors:   pq.StringArray{fmt.Sprintf("author%v", i)},
			Title:     fmt.Sprintf("название%v", i),
			UrlId:     uint64(i),
			ImageFile: fmt.Sprintf("image_file%v", i),
		}
		test_db.Create(&bookModel)
		ret = append(ret, bookModel)
	}
	return ret
}

func HeaderTokenMock(req *http.Request, u uint) {
	req.Header.Set("Authorization", fmt.Sprintf("JWT %v", common.GenToken(u)))
}

func HeaderSetTokenMock(req *http.Request, token string) {
	req.Header.Set("Authorization", token)
}

func TestBook(t *testing.T) {
	asserts := assert.New(t)
	common.TestDBFree()
	test_db = common.TestDBInit()
	BookModelMocker(10)

	r := gin.New()
	BooksRegister(r.Group("books/"))
	var unauthRequestTests = []request{
		{
			func(req *http.Request) {},
			"/books/search?q=наз",
			"GET",
			``,
			http.StatusOK,
			`{"data":.*,"success":true}`,
			"valid data and should return StatusOK",
		},
		{
			func(req *http.Request) {},
			"/books/search?k=tit",
			"GET",
			``,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"valid data and should return StatusUnprocessableEntity",
		},
	}
	for _, req := range unauthRequestTests {
		bodyData := req.bodyData
		req_serv, err := http.NewRequest(req.method, req.url, bytes.NewBufferString(bodyData))
		req_serv.Header.Set("Content-Type", "application/json")
		asserts.NoError(err)

		req.init(req_serv)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req_serv)

		asserts.Equal(req.expectedCode, w.Code, "Response Status - "+req.msg)
		asserts.Regexp(req.responseRegexg, w.Body.String(), "Response Content - "+req.msg)
	}
	common.TestDBFree()
}
