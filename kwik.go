package main

import (
	"github.com/AubSs/fasthttplogger"
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2005"
	}

	s := &fasthttp.Server{
		Handler: fasthttplogger.CombinedColored(Router().Handler),
		Name:    "FastHttpLogger",
	}

	s.ListenAndServe(":" + port)
}
