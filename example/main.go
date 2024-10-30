package main

import (
	"context"
	"example/resource"
	"github.com/lvyahui8/ellyn/api"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)
import "github.com/gin-gonic/gin"

var db = make(map[string]string)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers",
				"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,X-Ellyn-Gid, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Wrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Ellyn-Gid", strconv.FormatUint(api.Agent.GetGraphId(), 10))
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.Use(Cors())
	r.Use(Wrapper())

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

	r.GET("/trade", func(c *gin.Context) {
		Trade()
		c.JSON(http.StatusOK, gin.H{"code": 1})
	})

	r.GET("/profile/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		homeResource := resource.HomeResource{}
		uid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errMsg": err.Error()})
		} else {
			user, posts := homeResource.MyProfile(context.WithValue(context.Background(), "uid", uint(uid)))
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
		// pprof
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
