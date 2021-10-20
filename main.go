package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/isaqueveras/scrape-google-results/crawler"
)

func main() {
	r := gin.Default()

	crawler.Router(&r.RouterGroup)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Erro ao inicializar aplicação: " + err.Error())
	}
}
