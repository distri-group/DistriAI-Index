package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"encoding/json"
	"fmt"
)

type MetadataJson struct {
	GPUInfo      GPUInfo
	LocationInfo LocationInfo
	InfoFlop     InfoFlop
	InfoMemory   InfoMemory
}

type GPUInfo struct {
	Model  string
	Number uint32
}

type LocationInfo struct {
	Country string
}

type InfoFlop struct {
	Flops float64
}

type InfoMemory struct {
	Ram float64 `json:"RAM"`
}

func buildMachineModel(m distri_ai.Machine) model.Machine {
	var mj MetadataJson
	if err := json.Unmarshal([]byte(m.Metadata), &mj); err != nil {
		logs.Error(fmt.Sprintf("Unmarshal 'Metadata' error: %s \n", err))
	}

	return model.Machine{
		Owner:                  m.Owner.String(),
		Uuid:                   fmt.Sprintf("%#x", m.Uuid),
		Metadata:               m.Metadata,
		Status:                 uint8(m.Status),
		Price:                  m.Price,
		MaxDuration:            m.MaxDuration,
		Disk:                   m.Disk,
		CompletedCount:         m.CompletedCount,
		FailedCount:            m.FailedCount,
		Score:                  uint32(m.Score),
		ClaimedPeriodicRewards: m.ClaimedPeriodicRewards,
		ClaimedTaskRewards:     m.ClaimedTaskRewards,
		Gpu:                    mj.GPUInfo.Model,
		GpuCount:               mj.GPUInfo.Number,
		Region:                 mj.LocationInfo.Country,
		Tflops:                 mj.InfoFlop.Flops,
		Ram:                    mj.InfoMemory.Ram,
	}
}
