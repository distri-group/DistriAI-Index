package main

import (
	"distriai-backend-solana/chain"
	"distriai-backend-solana/common"
	"distriai-backend-solana/middleware"
	"distriai-backend-solana/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	common.InitConfig()
	common.InitDatabase()
	chain.Sync()

	gin.SetMode(common.Conf.Server.Mode)
	router := gin.Default()
	router.Use(middleware.Cors())
	routes.RegisterRoutes(router)
	if err := router.Run("0.0.0.0:" + common.Conf.Server.Port); err != nil {
		panic(err)
	}
}
