package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type OrderMetadataJson struct {
	FormData `json:"formData"`
}

type FormData struct {
	OrderTime time.Time `json:"orderTime"`
}

func buildOrderModel(o distri_ai.Order) model.Order {
	var mj OrderMetadataJson
	if err := json.Unmarshal([]byte(o.Metadata), &mj); err != nil {
		log.Printf("Unmarshal 'Metadata' error: %s \n", err)
	}

	return model.Order{
		Uuid:        fmt.Sprintf("%#x", o.OrderId),
		Buyer:       o.Buyer.String(),
		Seller:      o.Seller.String(),
		MachineUuid: fmt.Sprintf("%#x", o.MachineId),
		Price:       o.Price,
		Duration:    o.Duration,
		Total:       o.Total,
		Metadata:    o.Metadata,
		Status:      uint8(o.Status),
		OrderTime:   mj.FormData.OrderTime,
	}
}
