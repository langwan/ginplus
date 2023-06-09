package ginplus

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"reflect"
)

var Trans ut.Translator

func UseValidator() gin.HandlerFunc {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhHans := zh_Hans.New()
		uni := ut.New(zhHans)

		Trans, _ = uni.GetTranslator("")
		zh.RegisterDefaultTranslations(v, Trans)

		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("label")
		})

		v.RegisterTranslation("required", Trans, func(ut ut.Translator) error {
			return ut.Add("required", "请输入{0}", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})
	}
	return func(c *gin.Context) {

	}
}

type JsonError struct {
	Response gin.H
}

func (j *JsonError) Error() string {
	marshal, _ := json.Marshal(j.Response)
	return string(marshal)
}

func NewJsonErrors() *JsonError {
	err := JsonError{}
	err.Response = gin.H{}
	return &err
}

func NewJsonError(name, message string) *JsonError {
	err := JsonError{}
	err.Response = gin.H{name: message}
	return &err
}

func (j *JsonError) Set(name, message string) {
	j.Response[name] = message
}

func Any(r any, relativePath string, handler any) {
	if n, ok := r.(*gin.Engine); ok {
		n.Any(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	} else if n, ok := r.(*gin.RouterGroup); ok {
		n.Any(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	}
}

func GET(r any, relativePath string, handler any) {
	if n, ok := r.(*gin.Engine); ok {
		n.GET(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	} else if n, ok := r.(*gin.RouterGroup); ok {
		n.Any(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	}
}

func POST(r any, relativePath string, handler any) {
	if n, ok := r.(*gin.Engine); ok {
		n.POST(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	} else if n, ok := r.(*gin.RouterGroup); ok {
		n.Any(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	}
}

func PUT(r any, relativePath string, handler any) {
	if n, ok := r.(*gin.Engine); ok {
		n.PUT(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	} else if n, ok := r.(*gin.RouterGroup); ok {
		n.Any(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	}
}

func DELETE(r any, relativePath string, handler any) {
	if n, ok := r.(*gin.Engine); ok {
		n.DELETE(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	} else if n, ok := r.(*gin.RouterGroup); ok {
		n.Any(relativePath, func(c *gin.Context) {
			execute(c, handler)
		})
	}
}

func execute(c *gin.Context, handler any) {
	if n, ok := handler.(func(*gin.Context)); ok {
		n(c)
	} else if n, ok := handler.(func(*gin.Context) error); ok {
		err := n(c)
		if err != nil {
			handlerError(c, err)
		} else {
			c.JSON(http.StatusOK, gin.H{"body": "ok"})
		}

	} else if n, ok := handler.(func(*gin.Context) (any, error)); ok {
		resp, err := n(c)
		handlerResponse(c, resp, err)
	} else {
		method := reflect.ValueOf(handler)
		in := make([]reflect.Value, 0)
		in = append(in, reflect.ValueOf(c))
		if method.Type().NumIn() == 2 {
			parameter := method.Type().In(1)
			req := reflect.New(parameter.Elem()).Interface()
			err := c.ShouldBindJSON(req)
			if err != nil {
				handlerValidator(c, req, err)
				return
			}
			in = append(in, reflect.ValueOf(req))
		}

		call := method.Call(in)
		if call != nil {
			if len(call) == 1 {
				if call[0].IsNil() {
					handlerError(c, nil)
				} else {
					callErr := call[0].Interface().(error)
					handlerError(c, callErr)
				}

			} else if len(call) == 2 {
				if call[1].IsNil() {
					handlerResponse(c, call[0].Interface(), nil)
				} else {
					callErr := call[1].Interface().(error)
					handlerResponse(c, call[0].Interface(), callErr)
				}

			}
		}
	}
}

func handlerError(c *gin.Context, err error) {
	if err != nil {
		var jsonError *JsonError
		if errors.As(err, &jsonError) {
			resp := gin.H{}
			resp["errors"] = jsonError.Response
			c.JSON(http.StatusBadRequest, resp)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"body": "ok"})
	}
}

func handlerValidator(c *gin.Context, req any, err error) {
	resp := gin.H{}
	errs := gin.H{}
	elem := reflect.TypeOf(req).Elem()
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {

			if field, ok := elem.FieldByName(e.StructField()); ok {
				jsonName, _ := field.Tag.Lookup("json")
				errs[jsonName] = e.Translate(Trans)
			}
		}
	}
	resp["errors"] = errs
	c.JSON(http.StatusBadRequest, resp)
}

func handlerResponse(c *gin.Context, response any, err error) {
	resp := gin.H{}
	if err != nil {
		handlerError(c, err)
	} else {
		resp["body"] = response
		c.JSON(http.StatusOK, resp)
	}
}

func ResponseError(name string, err error) error {
	resp := gin.H{}
	h := gin.H{name: err.Error()}
	resp["errors"] = h
	marshal, _ := json.Marshal(resp)
	return errors.New(string(marshal))
}

func HandleRecovery(c *gin.Context, err any) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": gin.H{".exception": err}})
}

func UseRecovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, HandleRecovery)
}
