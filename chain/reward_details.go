package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
)

// Convert blockchain format to application format
func buildRewardModel(r distri_ai.Reward) model.Reward {
	return model.Reward{
		Period:     r.Period,
		Pool:       r.Pool,
		MachineNum: r.MachineNum,
	}
}
