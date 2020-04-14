package main

import (
	"html/template"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/valyala/fasthttp"
)

const EmptyPageText = "Page does not exist. Click on the Go button above to create it."

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
		ctx.Redirect("/docs/Main_page", fasthttp.StatusMovedPermanently)
	})

	// /edit launch edit mode from top bar
	router.GET("/edit", func(ctx *fasthttp.RequestCtx) {
		term := string(ctx.FormValue("term"))
		redirection := "/docs/"+term+"/edit"
		if term == "" {
			redirection = "/docs"
		}
		ctx.Redirect(redirection, fasthttp.StatusMovedPermanently)
	})

	// /search find pages
	router.GET("/search", func(ctx *fasthttp.RequestCtx) {
		term := string(ctx.FormValue("term"))
		if term == "" {
			ctx.Redirect("/docs", fasthttp.StatusMovedPermanently)
			return
		}
		files, _ := ioutil.ReadDir("pages/")
		var pages []string
		for _, f := range files {
			title := f.Name()
			if strings.Contains(strings.ToLower(title), strings.ToLower(term)) {
				pages = append(pages, title)
			}
		}

		var content = map[string]string{}
		results, _ := exec.Command("bash", "-c", "cd pages; grep '"+term+"' *").Output() // TODO case insensitive search
		for _, result := range strings.Split(string(results), "\n") {
			matches := strings.Split(result, ":")
			page := matches[0]
			if page == "" {
				continue
			}
			line := strings.Join(matches[1:], "\n")
			content[page] += strings.TrimSpace(line) + "\n"
		}

		page := Page{"Main page", "", ""}
		t := template.Must(template.ParseFiles("views/layout.html", "views/search.html"))
		t.Execute(ctx, H{
			"page":    page,
			"pages":   pages,   // pages with search term in name
			"content": content, // pages with search term in content
			"term":    term,
			"title":   "Main page",
		})
		ctx.SetContentType("text/html")
	})

	// /docs display all pages
	router.GET("/docs", func(ctx *fasthttp.RequestCtx) {
		files, _ := ioutil.ReadDir("pages/")
		var pages []string
		for _, f := range files {
			title := f.Name()
			if title[0] != '.' && title != "Main_page" {
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
			"title": "All",
		})
		ctx.SetContentType("text/html")
	})

	// /docs/:title display page
	router.GET("/docs/:title", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		if title == "All" {
			ctx.Redirect("/docs", fasthttp.StatusMovedPermanently)
			return
		}
		page := loadPage(title)
		term := ""
		if page.Source == EmptyPageText {
			term = title
		}
		t := template.Must(template.ParseFiles("views/layout.html", "views/show.html"))
		t.Execute(ctx, H{
			"page":  page,
			"title": strings.ReplaceAll(title, "_", " "),
			"term":  term,
		})
		ctx.SetContentType("text/html")
	})

	// /docs/:title/edit edit page
	router.GET("/docs/:title/edit", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		page := loadPage(title)
		t := template.Must(template.ParseFiles("views/layout.html", "views/edit.html"))
		t.Execute(ctx, H{
			"page":  page,
			"title": strings.ReplaceAll(title, "_", " "),
		})
		ctx.SetContentType("text/html")
	})

	// /docs/:title/edit edit page action
	router.POST("/docs/:title/edit", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		source := ctx.FormValue("source")
		ioutil.WriteFile("pages/"+title+".mw.html.md", []byte(source), 0644)

		ctx.Redirect("/docs/"+title, fasthttp.StatusMovedPermanently)
	})

	// /docs/:title/preview preview page action
	router.POST("/docs/:title/preview", func(ctx *fasthttp.RequestCtx) {
		title := pageTitle(ctx)
		source := ctx.FormValue("source")
		page := Page{
			title,
			string(source),
			template.HTML(github_flavored_markdown.Markdown(source)),
		}

		t := template.Must(template.ParseFiles("views/layout.html", "views/edit.html"))
		t.Execute(ctx, H{
			"page":  page,
			"title": strings.ReplaceAll(title, "_", " "),
		})
		ctx.SetContentType("text/html")
	})

	return router
}

func loadPage(title string) Page {
	source, err := ioutil.ReadFile("pages/" + title)
	if err != nil {
		source = []byte(EmptyPageText)
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
