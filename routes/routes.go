package routes

import (
	"distriai-index-solana/handlers"
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
	reward := engine.Group("/reward")
	{
		reward.POST("/total", handlers.RewardTotal)
		reward.POST("/claimable/list", handlers.RewardClaimableList)
		reward.POST("/period/list", handlers.RewardPeriodList)
		reward.POST("/machine/list", handlers.RewardMachineList)
	}
	log := engine.Group("/log")
	{
		log.POST("/add", handlers.LogAdd)
		log.POST("/list", handlers.LogList)
	}
	model := engine.Group("/model")
	{
		model.POST("/create", handlers.ModelCreate)
		model.POST("/list", handlers.ModelList)
	}
	engine.POST("/faucet", handlers.Faucet)
}
