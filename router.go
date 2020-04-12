package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/github_flavored_markdown"
)

type Page struct {
	Title  string
	Source string
	Body   template.HTML
}

func Router() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	router.Static("/assets", "./assets")

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
	authorized.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "root")
	})

	// /:page display page
	authorized.GET("/docs/:page", func(c *gin.Context) {
		page := loadPage(c)
		c.HTML(http.StatusOK, "show.html", gin.H{
			"page": page,
		})
	})

	// /:page/edit edit page
	authorized.GET("/docs/:page/edit", func(c *gin.Context) {
		page := loadPage(c)
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"page": page,
		})
	})

	// /:page/edit edit page action
	authorized.POST("/docs/:page/edit", func(c *gin.Context) {
		title := c.Param("page")
		source := c.PostForm("source")
		ioutil.WriteFile("pages/"+title+".mw.html.md", []byte(source), 0644)

		c.Redirect(http.StatusMovedPermanently, "/docs/"+title)
	})

	// /:page/edit edit page action
	authorized.POST("/docs/:page/preview", func(c *gin.Context) {
		source := c.PostForm("source")
		page := Page{
			c.Param("page"),
			source,
			template.HTML(github_flavored_markdown.Markdown([]byte(source))),
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"page": page,
		})
	})

	return router
}

func loadPage(c *gin.Context) Page {
	title := c.Param("page")

	source, err := ioutil.ReadFile("pages/" + title + ".mw.html.md")
	if err != nil {
		source = []byte("not found")
	}

	page := Page{
		title,
		string(source),
		template.HTML(github_flavored_markdown.Markdown(source)),
	}
	return page
}
