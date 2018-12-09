package user

import (
	"books-backend/app/book"
	"books-backend/app/common"
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

var image_url = "https://golang.org/doc/gopher/frontpage.png"
var test_db *gorm.DB

func newUserModel() UserModel {
	return UserModel{
		ID:           2,
		Username:     "asd123!@#ASD",
		Email:        "wzt@g.cn",
		Bio:          "heheda",
		Image:        &image_url,
		PasswordHash: "",
	}
}

type request struct {
	init           func(*http.Request)
	url            string
	method         string
	bodyData       string
	expectedCode   int
	responseRegexg string
	msg            string
}

func userModelMocker(n int) []UserModel {
	var offset int
	test_db.Model(&UserModel{}).Count(&offset)
	var ret []UserModel
	for i := offset + 1; i <= offset+n; i++ {
		image := fmt.Sprintf("http://image/%v.jpg", i)
		userModel := UserModel{
			Username: fmt.Sprintf("user%v", i),
			Email:    fmt.Sprintf("user%v@linkedin.com", i),
			Bio:      fmt.Sprintf("bio%v", i),
			Image:    &image,
		}
		userModel.setPassword("password123")
		test_db.Create(&userModel)
		ret = append(ret, userModel)
	}
	return ret
}

func bookModelMocker(n int) []book.BookModel {
	var offset int
	test_db.Model(&book.BookModel{}).Count(&offset)
	var ret []book.BookModel
	for i := offset + 1; i <= offset+n; i++ {
		bookModel := book.BookModel{
			Authors: pq.StringArray{fmt.Sprintf("author%v", i)},
			Title:   fmt.Sprintf("title%v", i),
			UrlId: uint64(i),
			ImageFile: fmt.Sprintf("image_file%v", i),
		}
		err := test_db.Create(&bookModel)
		fmt.Println(err)
		ret = append(ret, bookModel)
	}
	return ret
}

func HeaderTokenMock(req *http.Request, u uint) {
	token := common.GenToken(u)
	fmt.Println(token)
	req.Header.Set("Authorization", fmt.Sprintf("JWT %v", common.GenToken(u)))
}

func TestUserSignUp(t *testing.T) {
	asserts := assert.New(t)
	common.TestDBFree()
	common.TestDBInit()
	//userModelMocker(10)

	r := gin.New()
	UsersRegister(r.Group("/user"))
	var unauthRequestTests = []request{
		{
			func(req *http.Request) {},
			"/user/",
			"POST",
			`{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`,
			http.StatusCreated,
			`{"data":.*,"success":true}`,
			"valid data and should return StatusCreated",
		},
		{
			func(req *http.Request) {},
			"/user/",
			"POST",
			`{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"duplicated data and should return StatusUnprocessableEntity",
		},
		{
			func(req *http.Request) {
				//common.TestDBFree(test_db)
				//test_db = common.TestDBInit()
				//test_db.AutoMigrate(&UserModel{})
			},
			"/user/",
			"POST",
			`{"username": "u","email": "wzt@gg.cn","password": "jakejxke"}`,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"too short username should return error",
		},
		{
			func(req *http.Request) {
				//common.TestDBFree(test_db)
				//test_db = common.TestDBInit()
				//test_db.AutoMigrate(&UserModel{})
			},
			"/user/",
			"POST",
			`{"username": "wangzitian0","email": "wzt@gg.cn","password": "j"}`,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"too short password should return error",
		},
		{
			func(req *http.Request) {
				//common.TestDBFree(test_db)
				//test_db = common.TestDBInit()
				//test_db.AutoMigrate(&UserModel{})
			},
			"/user/",
			"POST",
			`{"username": "wangzitian0","email": "wztgg.cn","password": "jakejxke"}`,
			http.StatusUnprocessableEntity,
			`{"Email":"{key: email}","success":false}`,
			"email invalid should return error",
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

func TestUserLogin(t *testing.T) {
	asserts := assert.New(t)
	test_db = common.TestDBInit()

	r := gin.New()
	UsersRegister(r.Group("/user"))
	var unauthRequestTests = []request{
		{
			func(req *http.Request) {},
			"/user/",
			"POST",
			`{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`,
			http.StatusCreated,
			`{"data":.*,"success":true}`,
			"create user and should return StatusCreated",
		},
		{
			func(req *http.Request) {},
			"/user/login",
			"POST",
			`{"username": "wangzitian0", "password": "jakejxke"}`,
			http.StatusOK,
			`{.*"success":true.*}`,
			"successfully log in  and should return StatusOK",
		},
		{
			func(req *http.Request) {},
			"/user/login",
			"POST",
			`{"username": "wrong_username", "password": "jakejxke"}`,
			http.StatusForbidden,
			`{.*"success":false.*}`,
			"wrong username in  and should return StatusForbidden",
		},
		{
			func(req *http.Request) {},
			"/user/login",
			"POST",
			`{"username": "wangzitian0", "password": "wrongpassword"}`,
			http.StatusForbidden,
			`{.*"success":false.*}`,
			"wrong password  and should return StatusForbidden",
		},
		{
			func(req *http.Request) {},
			"/user/login",
			"POST",
			`{"usernakl;kme": "wangzitian0", "passwokkkkrd": "wrongpassword"}`,
			http.StatusUnprocessableEntity,
			`{.*"success":false.*}`,
			"wrong password  and should return StatusUnprocessableEntity",
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

func TestUserAddBooks(t *testing.T) {
	asserts := assert.New(t)
	test_db = common.TestDBInit()
	ret := bookModelMocker(10)
	fmt.Println(ret)

	r := gin.New()
	UsersRegister(r.Group("/user"))
	var unauthRequestTests = []request{
		{
			func(req *http.Request) {},
			"/user/",
			"POST",
			`{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`,
			http.StatusCreated,
			`{"data":.*,"success":true}`,
			"create user and should return StatusCreated",
		},
		{
			func(req *http.Request) {},
			"/user/login",
			"POST",
			`{"username": "wangzitian0", "password": "jakejxke"}`,
			http.StatusOK,
			`{.*"success":true.*}`,
			"successfully log in  and should return StatusOK",
		},
		{
			func(req *http.Request) {
				user, _ := FindOneUser(UserModel{Username: "user1"})
				fmt.Println(user)
				HeaderTokenMock(req, user.ID)
			},
			"/user/books", "POST",
			`{"book_id": [1,2,3,4]}`,
			http.StatusOK,
			`{.*"success":true.*}`,
			"add books and should return StatusOK",
		},
		//{
		//	func(req *http.Request) {},
		//	"/user/login",
		//	"POST",
		//	`{"username": "wangzitian0", "password": "wrongpassword"}`,
		//	http.StatusForbidden,
		//	`{.*"success":false.*}`,
		//	"wrong password  and should return StatusForbidden",
		//},
		//{
		//	func(req *http.Request) {},
		//	"/user/login",
		//	"POST",
		//	`{"usernakl;kme": "wangzitian0", "passwokkkkrd": "wrongpassword"}`,
		//	http.StatusUnprocessableEntity,
		//	`{.*"success":false.*}`,
		//	"wrong password  and should return StatusUnprocessableEntity",
		//},
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
	//common.TestDBFree(test_db)
}

func TestUserModel(t *testing.T) {
	asserts := assert.New(t)

	userModel := newUserModel()
	err := userModel.checkPassword("")
	asserts.Error(err, "empty password should return err")

	userModel = newUserModel()
	err = userModel.setPassword("")
	asserts.Error(err, "empty password can not be set null")

	userModel = newUserModel()
	err = userModel.setPassword("asd123!@#ASD")
	asserts.NoError(err, "password should be set successful")
	asserts.Len(userModel.PasswordHash, 60, "password hash length should be 60")

	err = userModel.checkPassword("sd123!@#ASD")
	asserts.Error(err, "password should be checked and not validated")

	err = userModel.checkPassword("asd123!@#ASD")
	asserts.NoError(err, "password should be checked and validated")
}
