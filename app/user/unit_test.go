package user

import (
	"books-backend/app/common"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

func TestUserSignUp(t *testing.T) {
	asserts := assert.New(t)

	r := gin.New()
	UsersRegister(r.Group("/user"))
	req := request{
		func(req *http.Request) {
			common.TestDBFree(test_db)
			test_db = common.TestDBInit()
			test_db.AutoMigrate(&UserModel{})
		},
		"/user/",
		"POST",
		`{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`,
		http.StatusCreated,
		`{"data":.*,"success":true}`,
		"valid data and should return StatusCreated",
	}
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

func TestUserModel(t *testing.T) {
	asserts := assert.New(t)

	//Testing UserModel's password feature
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
