package main

import "github.com/gin-gonic/gin"

func main() {
    str := returnStr()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pongss",
			"str": str,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func returnStr() string {
    str := "テストです"
    return str
}