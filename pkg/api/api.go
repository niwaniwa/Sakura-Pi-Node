package api

import (
	"Sakura-Hardware/pkg/device"
	"github.com/gin-gonic/gin"
	"net/http"
)

func initializeApiServer() {
	router := gin.Default()
	router.GET("/Key", getKeyState)
	router.POST("/Key", postKeyState)
	router.GET("/pasori", getID)
	go router.Run("localhost:5001")
}

func getKeyState(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, device.GetKeyState())
}

func postKeyState(c *gin.Context) {
	state := c.PostForm("state")
	c.IndentedJSON(http.StatusOK, state)
}

func getID(c *gin.Context) {

}
