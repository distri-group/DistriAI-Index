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

// fetch all account data on the Solana blockchain and storage
func fetchAllMachine(out rpc.GetProgramAccountsResult) {
	var machines []model.Machine
	for _, keyedAcct := range out {
		acct := keyedAcct.Account
		m := new(distri_ai.Machine)
		if err := m.UnmarshalWithDecoder(bin.NewBorshDecoder(acct.Data.GetBinary())); err != nil {
			continue
		}
		machine := buildMachineModel(*m)
		machines = append(machines, machine)
	}

	if len(machines) > 0 {
		if dbResult := common.Db.Create(&machines); dbResult.Error != nil {
			logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		}
	}
}

// Create a new  account
func addMachine(owner solana.PublicKey, uuid [16]uint8) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("machine"),
			owner[:],
			uuid[:],
		},
		distriProgramID,
	)
	if err != nil {
		logs.Error(fmt.Sprintf("FindProgramAddress error: %s \n", err))
		return
	}

	resp, err := rpcClient.GetAccountInfoWithOpts(
		context.TODO(),
		address,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		logs.Warn(fmt.Sprintf("GetAccountInfoWithOpts error: %s \n", err))
		return
	}
	var m distri_ai.Machine
	if err := m.UnmarshalWithDecoder(bin.NewBorshDecoder(resp.Value.Data.GetBinary())); err != nil {
		return
	}

	machine := buildMachineModel(m)
	if dbResult := common.Db.Create(&machine); dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}

// remove account
func removeMachine(owner solana.PublicKey, uuid [16]uint8) {
	dbResult := common.Db.
		Where("owner = ?", owner.String()).
		Where("uuid = ?", fmt.Sprintf("%#x", uuid)).
		Delete(&model.Machine{})
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}

// update account
func updateMachine(owner solana.PublicKey, uuid [16]uint8) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("machine"),
			owner[:],
			uuid[:],
		},
		distriProgramID,
	)
	if err != nil {
		logs.Error(fmt.Sprintf("FindProgramAddress error: %s \n", err))
		return
	}

	resp, err := rpcClient.GetAccountInfoWithOpts(
		context.TODO(),
		address,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		logs.Warn(fmt.Sprintf("GetAccountInfoWithOpts error: %s \n", err))
		return
	}
	var m distri_ai.Machine
	if err := m.UnmarshalWithDecoder(bin.NewBorshDecoder(resp.Value.Data.GetBinary())); err != nil {
		return
	}

	uuidStr := fmt.Sprintf("%#x", uuid)
	var machine model.Machine
	dbResult := common.Db.
		Where("owner = ?", owner.String()).
		Where("uuid = ?", uuidStr).
		Take(&machine)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		return
	}

	updateMachine := buildMachineModel(m)
	updateMachine.Id = machine.Id
	dbResult = common.Db.Save(&updateMachine)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}
