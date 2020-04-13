package main

import (
	"html/template"
	"io/ioutil"

	"github.com/buaazp/fasthttprouter"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/valyala/fasthttp"
)

type Page struct {
	Title  string
	Source string
	Body   template.HTML
}

func Router() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.ServeFiles("/assets/*filepath", "assets")

	// / root path
	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		ctx.Redirect("/docs/"+title, fasthttp.StatusMovedPermanently)
	})

	// /:page display page
	router.GET("/docs/:page", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		page := loadPage(title)
		t := template.Must(template.ParseFiles("views/layout.html", "views/show.html"))
		t.Execute(ctx, map[string]interface{}{
			"page": page,
		})
		ctx.SetContentType("text/html")
	})

	// /:page/edit edit page
	router.GET("/docs/:page/edit", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		page := loadPage(title)
		t := template.Must(template.ParseFiles("views/layout.html", "views/edit.html"))
		t.Execute(ctx, map[string]interface{}{
			"page": page,
		})
		ctx.SetContentType("text/html")
	})

	// /:page/edit edit page action
	router.POST("/docs/:page/edit", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		source := ctx.FormValue("source")
		ioutil.WriteFile("pages/"+title+".mw.html.md", []byte(source), 0644)

		ctx.Redirect("/docs/"+title, fasthttp.StatusMovedPermanently)
	})

	// /:page/edit edit page action
	router.POST("/docs/:page/preview", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		source := ctx.FormValue("source")
		page := Page{
			title,
			string(source),
			template.HTML(github_flavored_markdown.Markdown([]byte(source))),
		}

		t := template.Must(template.ParseFiles("views/layout.html", "views/edit.html"))
		t.Execute(ctx, map[string]interface{}{
			"page": page,
		})
		ctx.SetContentType("text/html")
	})

	return router
}

func loadPage(title string) Page {
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

func pageTitle(ctx *fasthttp.RequestCtx) string {
	title, _ := ctx.UserValue("page").(string)
	if title == "" {
		title = "Main_page"
	}
	return title
}
