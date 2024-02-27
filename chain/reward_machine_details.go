package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"fmt"
)

// Convert blockchain format to application format
func buildRewardMachineModel(r distri_ai.RewardMachine) model.RewardMachine {
	return model.RewardMachine{
		Period:    r.Period,
		Owner:     r.Owner.String(),
		MachineId: fmt.Sprintf("%#x", r.MachineId),
		TaskNum:   r.TaskNum,
		Claimed:   r.Claimed,
	}
}
