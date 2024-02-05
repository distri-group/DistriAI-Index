package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"fmt"
	"time"
)

// Convert blockchain format to application format
func buildOrderModel(o distri_ai.Order) model.Order {
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
		OrderTime:   time.Unix(o.OrderTime, 0),
		RefundTime:  time.Unix(o.RefundTime, 0),
	}
}
