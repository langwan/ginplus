package ginplus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}
func Test_A(t *testing.T) {
	r := gin.New()
	r.Use(UseValidator())
	Any(r, "/hello", hello)

	Any(r, "/hello2", func(c *gin.Context) error {
		return nil
	})
	Any(r, "/hello3", func(c *gin.Context) error {
		return NewJsonError("errmsg", "hello3 error")
	})

	type helloResponse struct {
		Message string `json:"message"`
	}
	Any(r, "/hello4", func(c *gin.Context) (*helloResponse, error) {
		return &helloResponse{Message: "hello4 message"}, nil
	})
	Any(r, "/hello5", func(c *gin.Context) (*helloResponse, error) {
		return nil, NewJsonError("errmsg", "hello5 error")
	})

	type helloRequest struct {
		Message string `json:"message" binding:"required" label:"消息"`
	}

	Any(r, "/hello6", func(c *gin.Context, request *helloRequest) error {
		fmt.Println("request", request.Message)
		return nil
	})
	Any(r, "/hello7", func(c *gin.Context, request *helloRequest) error {
		fmt.Println("request", request.Message)
		return NewJsonError("errmsg", "hello7 error")
	})

	Any(r, "/hello8", func(c *gin.Context, request *helloRequest) (*helloResponse, error) {
		fmt.Println("request", request.Message)
		return &helloResponse{Message: request.Message}, nil
	})
	Any(r, "/hello9", func(c *gin.Context, request *helloRequest) (*helloResponse, error) {
		fmt.Println("request", request.Message)
		return nil, NewJsonError("errmsg", "hello9 error")
	})

	Any(r, "/hello10", func(c *gin.Context, request *helloRequest) (*helloResponse, error) {
		fmt.Println("request", request.Message)
		return nil, NewJsonError("errmsg", "hello10 error")
	})

	r.Run(":4000")
}

func Test_Recovery(t *testing.T) {
	r := gin.Default()
	r.Use(UseRecovery())
	POST(r, "/api/login", func(g *gin.Context) {
		panic("err")
	})
	r.Run(":8810")
}

func UseRecovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, HandleRecovery)
}
