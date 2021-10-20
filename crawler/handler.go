package crawler

import (
	"log"

	"github.com/gin-gonic/gin"
)

func crawler(c *gin.Context) {
	searchTerm := c.Query("search")
	log.Println(searchTerm)

	res, err := GoogleScrape(searchTerm, "br", "pt-BR", nil, 1, 5)
	if err != nil {
		log.Print(err.Error())
		return
	}

	c.JSON(200, res)
}
