package main

import (
	"github.com/dncrht/kwik/controllers"
	"github.com/buaazp/fasthttprouter"
)

func Router() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.ServeFiles("/assets/*filepath", "assets")

	// / root path
	router.GET("/", controllers.Root)

	// /edit launch edit mode from top bar
	router.GET("/edit", controllers.EditRedirection)

	// /search find pages
	router.GET("/search", controllers.Search)

	// /docs display all pages
	router.GET("/docs", controllers.ShowAll)

	// /docs/:title display page
	router.GET("/docs/:title", controllers.Show)

	// /docs/:title/edit edit page
	router.GET("/docs/:title/edit", controllers.Edit)

	// /docs/:title/edit edit page action
	router.POST("/docs/:title/edit", controllers.EditAction)

	// /docs/:title/preview preview page action
	router.POST("/docs/:title/preview", controllers.PreviewAction)

	return router
}
