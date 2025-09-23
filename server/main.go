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
	server1URL := "http://" + viper.GetString("JSERVER1_NAME") + ":" + viper.GetString("JSERVER1_PORT")
	server2URL := "http://" + viper.GetString("JSERVER2_NAME") + ":" + viper.GetString("JSERVER2_PORT")
	
	server1Repo := NewServer1Repository(server1URL)
	server2Repo := NewServer2Repository(server2URL)
	
	flightService := NewFlightService([]FlightRepository{server1Repo, server2Repo})
	
	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Serveur en bonne santé",
		})
	})
	
	r.GET("/flight", func(c *gin.Context) {
		sortBy := c.DefaultQuery("sort", "price")
		
		flights, err := flightService.GetFlightsSortedBy(c, sortBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"flights": flights,
			"sort_by": sortBy,
		})
	})
	
	log.Printf("Serveur démarré sur le port %s", port)
	r.Run(":" + port)
}
