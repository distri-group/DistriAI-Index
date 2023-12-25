package chain

import (
	"distriai-backend-solana/chain/distri_ai"
	"distriai-backend-solana/model"
	"encoding/json"
	"fmt"
	"log"
)

type MetadataJson struct {
	GPUInfo      GPUInfo
	LocationInfo LocationInfo
	InfoFlop     InfoFlop
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

func buildMachineModel(m distri_ai.Machine) model.Machine {
	var mj MetadataJson
	if err := json.Unmarshal([]byte(m.Metadata), &mj); err != nil {
		log.Printf("Unmarshal 'Metadata' error: %s \n", err)
	}

	return model.Machine{
		Owner:          m.Owner.String(),
		Uuid:           fmt.Sprintf("%#x", m.Uuid),
		Metadata:       m.Metadata,
		Status:         uint8(m.Status),
		Price:          m.Price,
		MaxDuration:    m.MaxDuration,
		Disk:           m.Disk,
		CompletedCount: m.CompletedCount,
		FailedCount:    m.FailedCount,
		Score:          uint32(m.Score),
		Gpu:            mj.GPUInfo.Model,
		GpuCount:       mj.GPUInfo.Number,
		Region:         mj.LocationInfo.Country,
		Tflops:         mj.InfoFlop.Flops,
	}
}
