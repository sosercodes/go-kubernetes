package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getMessage(context *gin.Context) {
	// Add CORS headers
	//context.Header("Access-Control-Allow-Origin", "http://127.0.0.1:8181")
	//context.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	context.JSON(http.StatusOK, gin.H{
		"title": "Hello from Go!",
		"body":  "Welcome to Kubernetes.",
	})
}

func main() {
	router := gin.Default()
	router.GET("/message", getMessage)
	router.Run(":81")
}
