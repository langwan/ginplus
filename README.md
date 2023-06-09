# ginplus

`ginplus`是gin的一个微型扩展模块，`gin`原版只支持一种handle形式，`ginplus`支持多种API形式，支持表单验证、异常处理、自定义错误等。

## 调用方式

```go
func main() {
	r := gin.New()
	r.Use(ginplus.UseRecovery(), ginplus.UseValidator())
	ginplus.POST(r, "api/login", Login)
	acc := r.Group("account")
    ginplus.Any(acc, "create", Create)
	r.Run(":8010")
}
```

1. 需要返回错误

```go
func Hello(c *gin.Context) error {
	return NewJsonError(".message", "call function error")
}
```

2. 需要输入返回错误
```go
type helloRequest struct {
    Message string `json:"message"`
}
func(c *gin.Context, request *helloRequest) error {
    fmt.Println("request", request.Message)
    return nil
}
```

3. 需要输入返回输出和错误

```go

type helloRequest struct {
    Message string `json:"message"`
}

type helloResponse struct {
    Message string `json:"message"`
}

func(c *gin.Context, request *helloRequest) (*helloResponse, error) {
    fmt.Println("request", request.Message)
    return &helloResponse{Message: request.Message}, nil
}
```
4. 只需要输出和错误

```go

type helloResponse struct {
    Message string `json:"message"`
}


func(c *gin.Context) (*helloResponse, error) {
    return &helloResponse{Message: "ok"}, nil
}
```

5. 只需要输入

```go

type helloRequest struct {
    Message string `json:"message"`
}

func(c *gin.Context, request *helloRequest) (*helloResponse, error) {
    fmt.Println("request", request.Message)
}
```

6. 原始gin函数

```
func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}
```

## 异常处理

```go
func main() {
	r := gin.New()
	r.Use(ginplus.UseRecovery())
	ginplus.POST(r, "hello", func(c *gin.Context) {
        panic("系统异常")	
    })
	r.Run(":8010")
}
```

## 自定义错误

```go
func Login(c *gin.Context, req *request) error {
    return ginplus.NewJsonError(".message", "密码无效登录失败")
}
```

## 表单验证

```go
func main() {
	r := gin.New()
	r.Use(ginplus.UseValidator())
	ginplus.POST(r, "login", Login)
	r.Run(":8010")
}

type request struct {
    Email    string `json:"email" binding:"required" label:"邮箱地址"`
    Password string `json:"password" binding:"required" label:"密码"`
}

func Login(c *gin.Context, req *request) error {
    fmt.Println(req)
}

```
