package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("SERVER_PORT", "3001")
	
	port := viper.GetString("SERVER_PORT")
	
	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Serveur en bonne santé",
		})
	})
	
	log.Printf("Serveur démarré sur le port %s", port)
	r.Run(":" + port)
}
