package comment

import (
	"books/app/book"
	"books/app/common"
	"books/app/user"
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

func UserModelMocker(n int) []user.UserModel {
	var offset int
	test_db.Model(&user.UserModel{}).Count(&offset)
	var ret []user.UserModel
	for i := offset + 1; i <= offset+n; i++ {
		image := fmt.Sprintf("http://image/%v.jpg", i)
		userModel := user.UserModel{
			Username: fmt.Sprintf("user%v", i),
			Email:    fmt.Sprintf("user%v@linkedin.com", i),
			Bio:      fmt.Sprintf("bio%v", i),
			Image:    &image,
		}
		test_db.Create(&userModel)
		ret = append(ret, userModel)
	}
	return ret
}

func BookModelMocker(n int) []book.BookModel {
	var offset int
	test_db.Model(&book.BookModel{}).Count(&offset)
	var ret []book.BookModel
	for i := offset + 1; i <= offset+n; i++ {
		bookModel := book.BookModel{
			Authors:   pq.StringArray{fmt.Sprintf("author%v", i)},
			Title:     fmt.Sprintf("title%v", i),
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

func TestComment(t *testing.T) {
	asserts := assert.New(t)
	common.TestDBFree()
	test_db = common.TestDBInit()
	UserModelMocker(10)
	BookModelMocker(10)

	r := gin.New()
	user.UsersRegister(r.Group("/user"))
	CommentsRegister(r.Group("comments/"))
	var unauthRequestTests = []request{
		{
			func(req *http.Request) {
				user, _ := user.FindOneUser(user.UserModel{Username: "user1"})
				HeaderTokenMock(req, user.ID)
			},
			"/comments/books",
			"POST",
			`{"text": "Book comment","book_id": 1}`,
			http.StatusCreated,
			`{"data":.*,"success":true}`,
			"valid data and should return StatusCreated",
		},
		{
			func(req *http.Request) {
				user, _ := user.FindOneUser(user.UserModel{Username: "user2"})
				HeaderTokenMock(req, user.ID)
			},
			"/comments/books",
			"POST",
			`{"text": "Book comment","book_id": 1}`,
			http.StatusCreated,
			`{"data":.*,"success":true}`,
			"valid data and should return StatusCreated",
		},
		{
			func(req *http.Request) {},
			"/comments/books/1",
			"GET",
			`{}`,
			http.StatusOK,
			`{"data":.*,"success":true}`,
			"valid data and should return StatusOK",
		},
		{
			func(req *http.Request) {
				user, _ := user.FindOneUser(user.UserModel{Username: "user1"})
				HeaderTokenMock(req, user.ID)
			},
			"/comments/books",
			"DELETE",
			`{"comment_id": 1}`,
			http.StatusOK,
			`{"data":.*,"success":true}`,
			"valid data and should return StatusOK",
		},
		{
			func(req *http.Request) {
				user, _ := user.FindOneUser(user.UserModel{Username: "user1"})
				HeaderTokenMock(req, user.ID)
			},
			"/comments/books",
			"DELETE",
			`{"comment_i": 1}`,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"valid data and should return StatusUnprocessableEntity",
		},
		{
			func(req *http.Request) {
				user, _ := user.FindOneUser(user.UserModel{Username: "user1"})
				HeaderTokenMock(req, user.ID)
			},
			"/comments/books",
			"POST",
			`{"texlt": "Book comment","book_id": 1}`,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"not valid field and should return StatusUnprocessableEntity",
		},
		{
			func(req *http.Request) {
				user, _ := user.FindOneUser(user.UserModel{Username: "user1"})
				HeaderTokenMock(req, user.ID)
			},
			"/comments/books",
			"POST",
			`{"text": "Book comment","book_id": 100000}`,
			http.StatusBadRequest,
			`{.*"success":false.*}`,
			"not valid book_id and should return StatusBadRequest",
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
