package main

import "github.com/lvyahui8/ellyn/ellyn_agent"

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)
import "github.com/gin-gonic/gin"

func Sum(a, b int) int {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 0)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	time.Sleep(10 * time.Millisecond)
	return a + b
}

func Callback(fn func()) {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 1)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	fn()
}

func Handle() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 2)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	_, _, _ = WithUnusualParam(1, 2, "ccc", true, 0)
}

func Trade() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 3)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	a := 0
	for i := 0; i < 10; i++ {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
		a += i
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
	switch a {
	case 10:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 3)
		a -= 1
	case 20:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 4)
		a -= 2
	default:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 5)
		a -= 3
	}

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 6)
	defer func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 4)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		_ = recover()
	}()

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 7)
	defer Sum(1, 2)

	time.Sleep(5 * time.Millisecond)

	if a > 10 {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 8)
		a += 10
	} else {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 9)
		if a < 5 {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 10)
			a += 5
		} else {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 11)
			{
				a += 1
			}
		}
	}

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 12)
	for {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 13)
		a--
		if a <= 10 {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 14)
			break
		}
	}

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 15)
	Callback(func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 5)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		k := 0
		for i := 0; i < 10; i++ {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
			k++
		}
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
		if k > 10 {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 3)
			k = 0
		}
	})

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 16)
	go func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 6)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		Handle()
	}()

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 17)
	go func() { Handle() }()

	c := make(chan string, 100)
	c <- "hello"
	select {
	case msg := <-c:
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 18)
		_ = fmt.Sprintf("say: %v", msg)
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 19)
	if 1 == 1 {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 20)
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 21)
	func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 7)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	}()
	// empty
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 22)
	go func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 8)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		select {} // 阻塞
	}()
}

func N(n int) int {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 9)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	if n == 0 {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
		return 1
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
	time.Sleep(10 * time.Millisecond)
	return n * N(n-1)
}

func WithUnusualParam(a, b int, c string, _ bool, bool int) (x, y int, _ error) {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 10)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	time.Sleep(10 * time.Millisecond)
	if bool == 0 {
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
		return 0, 0, errors.New("test")
	}
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
	return 0, 0, nil
}

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 11)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 12)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		time.Sleep(101 * time.Millisecond)
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
	r.GET("/user/:name", func(c *gin.Context) {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 13)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
			{
				Callback(func() {
					_ellynCtx := ellyn_agent.Agent.GetCtx()
					ellyn_agent.Agent.Push(_ellynCtx, 14)
					defer ellyn_agent.Agent.Pop(_ellynCtx)
					ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
					Sum(1, 2)
				})
				ellyn_agent.Agent.VisitBlock(_ellynCtx, 3)
				c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
			}
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 2)
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 15)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	ellyn_agent.Agent.VisitBlock(_ellynCtx, 3)
	return r
}

func main() {
	_ellynCtx := ellyn_agent.Agent.GetCtx()
	ellyn_agent.Agent.Push(_ellynCtx, 16)
	defer ellyn_agent.Agent.Pop(_ellynCtx)
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
	go func() {
		_ellynCtx := ellyn_agent.Agent.GetCtx()
		ellyn_agent.Agent.Push(_ellynCtx, 17)
		defer ellyn_agent.Agent.Pop(_ellynCtx)
		ellyn_agent.Agent.VisitBlock(_ellynCtx, 0)
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	ellyn_agent.Agent.VisitBlock(_ellynCtx, 1)
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
