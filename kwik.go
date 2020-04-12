package main

import (
	"github.com/valyala/fasthttp"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fasthttp.ListenAndServe(":"+port, Router().Handler)
}
