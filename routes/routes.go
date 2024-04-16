package routes

import (
	"distriai-index-solana/handlers"
	"distriai-index-solana/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes unified management routes.
func RegisterRoutes(engine *gin.Engine) {
	mailbox := engine.Group("/mailbox")
	{
		mailbox.POST("/subscribe", handlers.Subscribe)
		mailbox.POST("/unsubscribe", handlers.Unsubscribe)
	}
	user := engine.Group("/user")
	{
		user.POST("/login", handlers.Login)
	}
	machine := engine.Group("/machine")
	{
		machine.POST("/filter", handlers.MachineFilter)
		machine.POST("/market", handlers.MachineMarket)
		machine.POST("/mine", handlers.MachineMine)
		machine.GET("/:Owner/:Uuid", handlers.MachineGet)
	}
	order := engine.Group("/order")
	{
		order.POST("/mine", handlers.OrderMine)
		order.POST("/all", handlers.OrderAll)
		order.GET("/:Uuid", handlers.OrderGet)
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
		model.POST("/list", handlers.ModelList)
		model.GET("/:Owner/:Name", handlers.ModelGet)
	}
	modelAuth := engine.Group("/model", middleware.Jwt())
	{
		modelAuth.POST("/create", handlers.ModelCreate)
		modelAuth.POST("/presign", handlers.ModelPresign)
	}
	engine.POST("/faucet", handlers.Faucet)
	dataset := engine.Group("/dataset", middleware.Jwt())
	{
		dataset.POST("/create", handlers.DataSetCreate)
		dataset.POST("/presign", handlers.DatasetPresign)
	}
	datasetInfo := engine.Group("/datasetinfo")
	{
		datasetInfo.POST("/list", handlers.DatasetList)
		datasetInfo.GET("/:Owner/:Name", handlers.DatasetGet)
	}
}
