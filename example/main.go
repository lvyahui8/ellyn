package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)
import "github.com/gin-gonic/gin"

func Sum(a, b int) int {
	return a + b
}

func Callback(fn func()) {
	fn()
}

func Handle() {

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
	return n * N(n-1)
}

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
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
