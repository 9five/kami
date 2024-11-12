package main

import (
	// "github.com/gin-gonic/gin"
	// "kami/config"
	"kami/router"
)

func main() {
	r := router.NewRouter()
	r.Run(":8000")
}
