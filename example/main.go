package main

import (
	"context"
	"errors"
	"example/resource"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"time"
)
import "github.com/gin-gonic/gin"

func init() {
	fmt.Printf("init example main")
}

func Sum(a, b int) int {
	time.Sleep(10 * time.Millisecond)
	return a + b
}

func Callback(fn func()) {
	fn()
}

func Handle() {
	_, _, _ = WithUnusualParam(1, 2, "ccc", true, 0)
	NoName(1, "xx")
}

func Trade() {
	a := 0
	for i := 0; i < 10; i++ {
		a += i
	}
	switch a {
	case 10:
		a -= 1
	case 20:
		a -= 2
	default:
		a -= 3
	}

	defer func() {
		_ = recover()
	}()

	defer Sum(1, 2)

	time.Sleep(5 * time.Millisecond)

	if a > 10 {
		a += 10
	} else if a < 5 {
		a += 5
	} else {
		a += 1
	}

	for {
		a--
		if a <= 10 {
			break
		}
	}

	Callback(func() {
		k := 0
		for i := 0; i < 10; i++ {
			k++
		}
		if k > 10 {
			k = 0
		}
	})

	go func() {
		Handle()
	}()

	go Handle()

	c := make(chan string, 100)
	c <- "hello"
	select {
	case msg := <-c:
		_ = fmt.Sprintf("say: %v", msg)
	}
	if 1 == 1 {
	}
	func() {
	}()
	// empty
	go func() {
		select {} // 阻塞
	}()
}

func N(n int) int {
	if n == 0 {
		return 1
	}
	time.Sleep(10 * time.Millisecond)
	return n * N(n-1)
}

func WithUnusualParam(a, b int,
	c string, _ bool, bool int) (x, y int,
	_ error) {
	time.Sleep(10 * time.Millisecond)
	if bool == 0 {
		return 0, 0, errors.New("test")
	}
	return 0, 0, nil
}

func NoName(int,
	string) (byte, bool) {
	_ = fmt.Sprintf("x%d", runtime.NumCPU())
	return 0, false
}

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(101 * time.Millisecond)
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			Callback(func() {
				Sum(1, 2)
			})
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	r.GET("/profile/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		homeResource := resource.HomeResource{}
		uid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errMsg": err.Error()})
		} else {
			user, posts := homeResource.MyProfile(context.WithValue(context.Background(), "id", uid))
			c.JSON(http.StatusOK, gin.H{"user": user, "posts": posts})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
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
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
