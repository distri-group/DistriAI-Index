package chain

import (
	"context"
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"encoding/binary"
	"fmt"
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

// fetch all account data on the Solana blockchain and storage
func fetchAllRewardMachine(out rpc.GetProgramAccountsResult) {
	var rewardMachines []model.RewardMachine
	for _, keyedAcct := range out {
		acct := keyedAcct.Account
		r := new(distri_ai.RewardMachine)
		if err := r.UnmarshalWithDecoder(bin.NewBorshDecoder(acct.Data.GetBinary())); err != nil {
			continue
		}
		rewardMachine := buildRewardMachineModel(*r)
		rewardMachines = append(rewardMachines, rewardMachine)
	}

	if len(rewardMachines) > 0 {
		if dbResult := common.Db.Create(&rewardMachines); dbResult.Error != nil {
			log.Printf("Database error: %s \n", dbResult.Error)
		}
	}
}

// Create or update account
func saveRewardMachine(period uint32, owner solana.PublicKey, machineId [16]uint8) {
	var periodBytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(periodBytes, period)

	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("reward-machine"),
			periodBytes,
			owner[:],
			machineId[:],
		},
		distriProgramID,
	)
	if err != nil {
		return
	}
	var r distri_ai.RewardMachine
	if err := rpcClient.GetAccountDataBorshInto(context.TODO(), address, &r); err != nil {
		return
	}

	saveRewardMachine := buildRewardMachineModel(r)
	machineIdStr := fmt.Sprintf("%#x", machineId)
	var rewardMachine model.RewardMachine
	dbResult := common.Db.
		Where("period = ?", period).
		Where("owner = ?", owner.String()).
		Where("machine_id = ?", machineIdStr).
		Take(&rewardMachine)
	if dbResult.Error == nil {
		saveRewardMachine.Id = rewardMachine.Id
	}

	dbResult = common.Db.Save(&saveRewardMachine)
	if dbResult.Error != nil {
		log.Printf("Database error: %s \n", dbResult.Error)
	}
}
