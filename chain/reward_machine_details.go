package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"fmt"
)

// Convert blockchain format to application format
func buildRewardMachineModel(o distri_ai.RewardMachine) model.RewardMachine {
	return model.RewardMachine{
		Period:         o.Period,
		Owner:          o.Owner.String(),
		MachineId:      fmt.Sprintf("%#x", o.MachineId),
		TaskNum:        o.TaskNum,
		Claimed:        o.Claimed,
		PeriodicReward: o.PeriodicReward,
	}
}
