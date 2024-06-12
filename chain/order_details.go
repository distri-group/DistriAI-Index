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
		StartTime:   time.Unix(o.StartTime, 0),
		RefundTime:  time.Unix(o.RefundTime, 0),
		Model1Owner: o.Model1Owner.String(),
		Model1Name:  o.Model1Name,
		Model2Owner: o.Model2Owner.String(),
		Model2Name:  o.Model2Name,
		Model3Owner: o.Model3Owner.String(),
		Model3Name:  o.Model3Name,
		Model4Owner: o.Model4Owner.String(),
		Model4Name:  o.Model4Name,
		Model5Owner: o.Model5Owner.String(),
		Model5Name:  o.Model5Name,
	}
}
