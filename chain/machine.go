package chain

import (
	"context"
	"distriai-backend-solana/chain/distri_ai"
	"distriai-backend-solana/common"
	"distriai-backend-solana/model"
	"fmt"
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

func fetchAllMachine() {
	resp, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		DistriAIProgramID,
		&rpc.GetProgramAccountsOpts{
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		log.Printf("GetProgramAccounts error: %s \n", err)
		return
	}

	var machines []model.Machine
	for _, keyedAcct := range resp {
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
			log.Printf("Database error: %s \n", dbResult.Error)
		}
	}
}

func addMachine(owner solana.PublicKey, uuid [16]uint8) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("machine"),
			owner[:],
			uuid[:],
		},
		DistriAIProgramID,
	)
	if err != nil {
		return
	}
	var m distri_ai.Machine
	if err := client.GetAccountDataBorshInto(context.TODO(), address, &m); err != nil {
		return
	}

	machine := buildMachineModel(m)
	if dbResult := common.Db.Create(&machine); dbResult.Error != nil {
		log.Printf("Database error: %s \n", dbResult.Error)
	}
}

func removeMachine(owner solana.PublicKey, uuid [16]uint8) {
	dbResult := common.Db.
		Where("owner = ?", owner.String()).
		Where("uuid = ?", fmt.Sprintf("%#x", uuid)).
		Delete(&model.Machine{})
	if dbResult.Error != nil {
		log.Printf("Database error: %s \n", dbResult.Error)
	}
}

func updateMachine(owner solana.PublicKey, uuid [16]uint8) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("machine"),
			owner[:],
			uuid[:],
		},
		DistriAIProgramID,
	)
	if err != nil {
		return
	}
	var m distri_ai.Machine
	if err := client.GetAccountDataBorshInto(context.TODO(), address, &m); err != nil {
		return
	}

	uuidStr := fmt.Sprintf("%#x", uuid)
	var machine model.Machine
	dbResult := common.Db.
		Where("owner = ?", owner.String()).
		Where("uuid = ?", uuidStr).
		Take(&machine)
	if dbResult.Error != nil {
		log.Printf("Database error: %s \n", dbResult.Error)
		return
	}

	updateMachine := buildMachineModel(m)
	updateMachine.Id = machine.Id
	dbResult = common.Db.Save(&updateMachine)
	if dbResult.Error != nil {
		log.Printf("Database error: %s \n", dbResult.Error)
	}
}
