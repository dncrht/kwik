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

type H map[string]interface{} // a standard hash map

func Router() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.ServeFiles("/assets/*filepath", "assets")

	// / root path
	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		ctx.Redirect("/docs/"+title, fasthttp.StatusMovedPermanently)
	})

	// /docs display all pages
	router.GET("/docs", func(ctx *fasthttp.RequestCtx) {
		files, _ := ioutil.ReadDir("pages/")
		var pages []string
		for _, f := range files {
			title := f.Name()
			if title[0] != '.' || title != "Main_page" {
				pages = append(pages, title)
			}
		}

		page := Page{
			"All",
			"",
			"",
		}
		t := template.Must(template.ParseFiles("views/layout.html", "views/show_all.html"))
		t.Execute(ctx, H{
			"page":  page,
			"pages": pages,
		})
		ctx.SetContentType("text/html")
	})

	// /:title display page
	router.GET("/docs/:title", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		if title == "All" {
			ctx.Redirect("/docs", fasthttp.StatusMovedPermanently)
			return
		}
		page := loadPage(title)
		t := template.Must(template.ParseFiles("views/layout.html", "views/show.html"))
		t.Execute(ctx, H{
			"page": page,
		})
		ctx.SetContentType("text/html")
	})

	// /:title/edit edit page
	router.GET("/docs/:title/edit", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		page := loadPage(title)
		t := template.Must(template.ParseFiles("views/layout.html", "views/edit.html"))
		t.Execute(ctx, H{
			"page": page,
		})
		ctx.SetContentType("text/html")
	})

	// /:title/edit edit page action
	router.POST("/docs/:title/edit", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		source := ctx.FormValue("source")
		ioutil.WriteFile("pages/"+title+".mw.html.md", []byte(source), 0644)

		ctx.Redirect("/docs/"+title, fasthttp.StatusMovedPermanently)
	})

	// /:title/edit edit page action
	router.POST("/docs/:title/preview", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		source := ctx.FormValue("source")
		page := Page{
			title,
			string(source),
			template.HTML(github_flavored_markdown.Markdown([]byte(source))),
		}

		t := template.Must(template.ParseFiles("views/layout.html", "views/edit.html"))
		t.Execute(ctx, H{
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
	title, _ := ctx.UserValue("title").(string)
	if title == "" {
		title = "Main_page"
	}
	return title
}
