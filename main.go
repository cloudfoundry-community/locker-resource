package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func main() {
	lockChan := make(chan LockRequest)
	lockConfig := os.Getenv("LOCKER_CONFIG")
	if lockConfig == "" {
		fmt.Fprintf(os.Stderr, "No LOCKER_CONFIG file specified. Bailing out.\n")
		os.Exit(1)
	}
	_, err := os.Stat(lockConfig)
	if os.IsNotExist(err) {
		err := ioutil.WriteFile(lockConfig, []byte("{}"), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create LOCKER_CONFIG file (%s): %s\n", lockConfig, err)
			fmt.Fprintf(os.Stderr, "Bailing out.\n")
			os.Exit(1)
		}
	}

	authUser := os.Getenv("AUTH_USER")
	authPass := os.Getenv("AUTH_PASS")
	if authUser != "" && authPass == "" {
		fmt.Fprintf(os.Stderr, "AUTH_USER specified, but no AUTH_PASS was provided. Bailing out.\n")
		os.Exit(1)
	}

	if authUser == "" && authPass != "" {
		fmt.Fprintf(os.Stderr, "AUTH_PASS specified, but no AUTH_USER was provided. Bailing out.\n")
		os.Exit(1)
	}

	go lockServer(lockChan, lockConfig)

	router := gin.Default()
	var routes *gin.RouterGroup
	if authUser != "" && authPass != "" {
		routes = router.Group("/", gin.BasicAuth(gin.Accounts{authUser: authPass}))
	} else {
		routes = router.Group("/", func(c *gin.Context) {})
	}
	routes.GET("/locks", func(c *gin.Context) {
		rc := make(chan LockResponse)

		lr := LockRequest{
			Command:  ListOp,
			Response: rc,
		}

		lockChan <- lr

		response := <-rc
		if response.Error != nil {
			fmt.Fprintf(os.Stderr, "Error listing locks: %s\n", response.Error)
			c.Render(500, render.JSON{Data: map[string]string{"error": response.Error.Error()}})
			return
		}

		c.Render(200, render.JSON{Data: response.Message})
	})
	routes.PUT("/lock/:pool", func(c *gin.Context) {
		rc := make(chan LockResponse)

		pool := c.Param("pool")
		input := LockInput{}
		c.BindJSON(&input)
		if input.Lock == "" {
			c.Render(400, render.JSON{Data: map[string]string{"error": "No lock specified in request body"}})
			return
		}

		lr := LockRequest{
			Command:  LockOp,
			Pool:     pool,
			Lock:     input.Lock,
			Response: rc,
		}
		lockChan <- lr
		res := <-rc
		if res.Error != nil {
			fmt.Fprintf(os.Stderr, "Error locking %s with %s: %s\n", pool, input.Lock, res.Error)
			c.Render(423, render.JSON{Data: map[string]string{"error": res.Error.Error()}})
			return
		}

		if res.Status == Locked {
			c.Render(200, render.JSON{Data: res.Message})
		} else {
			fmt.Fprintf(os.Stderr, "Unexpected response to a lock request: %#v\n", res)
			c.Render(500, render.JSON{Data: map[string]string{"error": "Unknown response from lock request"}})
		}
	})
	routes.DELETE("/lock/:pool", func(c *gin.Context) {
		rc := make(chan LockResponse)
		pool := c.Param("pool")
		input := LockInput{}
		c.BindJSON(&input)
		if input.Lock == "" {
			c.Render(400, render.JSON{Data: map[string]string{"error": "No lock specified in request body"}})
			return
		}

		lr := LockRequest{
			Command:  UnlockOp,
			Pool:     pool,
			Lock:     input.Lock,
			Response: rc,
		}

		lockChan <- lr
		res := <-rc
		if res.Error != nil {
			fmt.Fprintf(os.Stderr, "Error unlocking %s: %s\n", pool, res.Error)
			c.Render(423, render.JSON{Data: map[string]string{"error": res.Error.Error()}})
			return
		}

		if res.Status != Unlocked {
			fmt.Fprintf(os.Stderr, "Unexpected response to unlock request: %#v\n", res)
			c.Render(500, render.JSON{Data: map[string]string{"error": "Unknown response from unlock request"}})
			return
		}

		c.Render(200, render.JSON{Data: res.Message})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)
}
