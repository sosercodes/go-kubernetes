package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIpAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	return ipAddress.IP.String()
}

func getMessage(context *gin.Context) {
	// Add CORS headers
	//context.Header("Access-Control-Allow-Origin", "http://127.0.0.1:8181")
	//context.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	context.JSON(http.StatusOK, gin.H{
		"title": "Hello from Go!",
		"body":  "Welcome to Kubernetes pod@'" + getIpAddress() + "'.",
	})
}

func main() {
	router := gin.Default()
	router.GET("/message", getMessage)
	router.Run(":80")
}
