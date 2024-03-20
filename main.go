package main

import (
	"distriai-index-solana/common"
	"distriai-index-solana/middleware"
	"distriai-index-solana/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	common.Init()
	//chain.Sync()

	gin.SetMode(common.Conf.Server.Mode)
	router := gin.Default()
	router.Use(middleware.Cors())
	routes.RegisterRoutes(router)
	if err := router.Run("0.0.0.0:" + common.Conf.Server.Port); err != nil {
		panic(err)
	}
}
