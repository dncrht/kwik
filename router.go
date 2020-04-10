package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/github_flavored_markdown"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	var authorized *gin.RouterGroup
	if user != "" && password != "" {
		fmt.Println("* PROTECTED BY USER AND PASSWORD *")

		authorized = router.Group("/", gin.BasicAuth(gin.Accounts{
			user: password,
		}))
	} else {
		fmt.Println("* OPEN ACCESS *")

		authorized = router.Group("/")
	}

	// / root path
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "root")
	})

	// /:page display page
	authorized.GET("/:page", func(c *gin.Context) {
		page := c.Param("page")

		filename := "pages/" + page + ".mw.html.md"
		body, _ := ioutil.ReadFile(filename)
		body = github_flavored_markdown.Markdown(body)

		c.HTML(http.StatusOK, "show.html", gin.H{
			"body": template.HTML(body),
		})
	})

	return router
}
