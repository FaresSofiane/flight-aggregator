package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("SERVER_PORT", "3001")
	viper.SetDefault("JSERVER1_NAME", "j-server1")
	viper.SetDefault("JSERVER1_PORT", "4001")
	viper.SetDefault("JSERVER2_NAME", "j-server2")
	viper.SetDefault("JSERVER2_PORT", "4002")

	port := viper.GetString("SERVER_PORT")
	url1 := "http://" + viper.GetString("JSERVER1_NAME") + ":" + viper.GetString("JSERVER1_PORT")
	url2 := "http://" + viper.GetString("JSERVER2_NAME") + ":" + viper.GetString("JSERVER2_PORT")

	repo1 := NewServer1Repository(url1)
	repo2 := NewServer2Repository(url2)

	service := NewFlightService([]FlightRepository{repo1, repo2})

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/flight", func(c *gin.Context) {
		sortBy := c.DefaultQuery("sort", "price")

		flights, err := service.GetFlightsSortedBy(c, sortBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, flights)
	})

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
