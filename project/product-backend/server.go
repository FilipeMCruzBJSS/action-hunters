package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitServer(p Producer) *gin.Engine {
	r := gin.Default()

	api := r.Group("api/v1")

	api.POST("products", submitProduct(p))

	return r
}

func submitProduct(p Producer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var dto InputProductDto

		err := c.BindJSON(&dto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		out, err := Verify(dto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = p.Send(out)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, out)
	}
}
