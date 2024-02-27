package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"time"
)

// Convert blockchain format to application format
func buildRewardModel(r distri_ai.Reward) model.Reward {
	return model.Reward{
		Period:             r.Period,
		StartTime:          time.Unix(r.StartTime, 0),
		Pool:               r.Pool,
		MachineNum:         r.MachineNum,
		UnitPeriodicReward: r.UnitPeriodicReward,
		TaskNum:            r.TaskNum,
		UnitTaskReward:     r.UnitTaskReward,
	}
}
