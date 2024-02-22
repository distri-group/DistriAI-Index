package chain

import (
	"context"
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"encoding/binary"
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

// fetch all account data on the Solana blockchain and storage
func fetchAllReward(out rpc.GetProgramAccountsResult) {
	var rewards []model.Reward
	for _, keyedAcct := range out {
		acct := keyedAcct.Account
		r := new(distri_ai.Reward)
		if err := r.UnmarshalWithDecoder(bin.NewBorshDecoder(acct.Data.GetBinary())); err != nil {
			continue
		}
		reward := buildRewardModel(*r)
		rewards = append(rewards, reward)
	}

	if len(rewards) > 0 {
		if dbResult := common.Db.Create(&rewards); dbResult.Error != nil {
			log.Printf("Database error: %s \n", dbResult.Error)
		}
	}
}

// Create or update account
func saveReward(period uint32) {
	var periodBytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(periodBytes, period)

	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("reward"),
			periodBytes,
		},
		distriProgramID,
	)
	if err != nil {
		return
	}
	var r distri_ai.Reward
	if err := rpcClient.GetAccountDataBorshInto(context.TODO(), address, &r); err != nil {
		return
	}

	saveReward := buildRewardModel(r)
	var reward model.Reward
	dbResult := common.Db.
		Where("period = ?", period).
		Take(&reward)
	if dbResult.Error == nil {
		saveReward.Id = reward.Id
	}

	dbResult = common.Db.Save(&saveReward)
	if dbResult.Error != nil {
		log.Printf("Database error: %s \n", dbResult.Error)
	}
}
