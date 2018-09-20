package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

const NBSecretPassword = "A String Very Very Very Strong!!@##$!@#$"

func GenToken(id uint) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	jwt_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token, _ := jwt_token.SignedString([]byte(NBSecretPassword))
	return token
}

func KeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(NBSecretPassword), nil
}

func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		if v.Param != "" {
			res.Errors[v.Field] = fmt.Sprintf("{%v: %v}", v.Tag, v.Param)
		} else {
			res.Errors[v.Field] = fmt.Sprintf("{key: %v}", v.Tag)
		}

	}
	return res
}

func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

func RenderResponse(c *gin.Context, code int, errors interface{}, data interface{}) {
	body := gin.H{}

	if errors != nil {
		err := errors.(CommonError)
		for k, v := range err.Errors {
			body[k] = v
		}
	}
	if code >= 200 && code <= 299 {
		body["success"] = true
		body["data"] = data
	} else {
		body["success"] = false
	}
	c.JSON(code, body)
}
