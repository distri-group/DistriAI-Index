package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
)

// Convert blockchain format to application format
func buildRewardModel(o distri_ai.Reward) model.Reward {
	return model.Reward{
		Period:     o.Period,
		Pool:       o.Pool,
		MachineNum: o.MachineNum,
	}
}
