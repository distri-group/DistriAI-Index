package routes

import (
	"distriai-backend-solana/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	mailbox := engine.Group("/mailbox")
	{
		mailbox.POST("/subscribe", handlers.Subscribe)
		mailbox.POST("/unsubscribe", handlers.Unsubscribe)
	}
	machine := engine.Group("/machine")
	{
		machine.POST("/filter", handlers.MachineFilter)
		machine.POST("/market", handlers.MachineMarket)
		machine.POST("/mine", handlers.MachineMine)
	}
	order := engine.Group("/order")
	{
		order.POST("/mine", handlers.OrderMine)
		order.POST("/all", handlers.OrderAll)
	}
	log := engine.Group("/log")
	{
		log.POST("/add", handlers.LogAdd)
		log.POST("/list", handlers.LogList)
	}
}
