package chain

import (
	"context"
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"fmt"
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func fetchAllOrder(out rpc.GetProgramAccountsResult) {
	var orders []model.Order
	for _, keyedAcct := range out {
		acct := keyedAcct.Account
		o := new(distri_ai.Order)
		if err := o.UnmarshalWithDecoder(bin.NewBorshDecoder(acct.Data.GetBinary())); err != nil {
			continue
		}
		order := buildOrderModel(*o)
		orders = append(orders, order)
	}

	if len(orders) > 0 {
		if dbResult := common.Db.Create(&orders); dbResult.Error != nil {
			logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		}
	}
}

func addOrder(orderId [16]uint8, buyer solana.PublicKey) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("order"),
			buyer[:],
			orderId[:],
		},
		distriProgramID,
	)
	if err != nil {
		return
	}
	var o distri_ai.Order
	if err := rpcClient.GetAccountDataBorshInto(context.TODO(), address, &o); err != nil {
		return
	}

	order := buildOrderModel(o)
	if dbResult := common.Db.Create(&order); dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}

func removeOrder(orderId [16]uint8) {
	dbResult := common.Db.
		Where("uuid = ?", fmt.Sprintf("%#x", orderId)).
		Delete(&model.Order{})
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}

func updateOrder(orderId [16]uint8, buyer solana.PublicKey) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("order"),
			buyer[:],
			orderId[:],
		},
		distriProgramID,
	)
	if err != nil {
		logs.Warn(fmt.Sprintf("Can not find Program Address: %s \n", err))
		return
	}
	var o distri_ai.Order
	if err := rpcClient.GetAccountDataBorshInto(context.TODO(), address, &o); err != nil {
		logs.Warn(fmt.Sprintf("GetAccountDataBorshInto error: %s \n", err))
		return
	}

	uuidStr := fmt.Sprintf("%#x", orderId)
	var order model.Order
	dbResult := common.Db.
		Where("uuid = ?", uuidStr).
		Take(&order)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		return
	}

	updateOrder := buildOrderModel(o)
	updateOrder.Id = order.Id
	dbResult = common.Db.Save(&updateOrder)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}
