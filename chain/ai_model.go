package chain

import (
	"context"
	"crypto/sha256"
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
func fetchAllAiModel(out rpc.GetProgramAccountsResult) {
	var aiModels []model.AiModel
	
	// Iterate through each keyed account in the 'out' slice
	for _, keyedAcct := range out {
		acct := keyedAcct.Account
		m := new(distri_ai.AiModel)
		
		// Attempt to decode the account data into an AiModel object
		if err := m.UnmarshalWithDecoder(bin.NewBorshDecoder(acct.Data.GetBinary())); err != nil {
			continue
		}

		// Build an AiModel database model object
		machine := buildAiModelModel(m)
		aiModels = append(aiModels, machine)
	}

	// If aiModels slice is not empty, attempt to bulk insert into the database
	if len(aiModels) > 0 {
		if dbResult := common.Db.Create(&aiModels); dbResult.Error != nil {
			logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		}
		createAiModelHeats(aiModels)
	}
}

// createAiModelHeats processes a slice of AiModel and creates a new slice of AiModelHeat
func createAiModelHeats(aiModels []model.AiModel) {
	var heats []model.AiModelHeat
	for _, aiModel := range aiModels {
		heat := model.AiModelHeat{Owner: aiModel.Owner, Name: aiModel.Name}
		var count int64
		common.Db.Model(&heat).Where(&heat).Count(&count)
		if count == 0 {
			heats = append(heats, heat)
		}
	}
	if len(heats) > 0 {
		if dbResult := common.Db.Create(&heats); dbResult.Error != nil {
			logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		}
	}
}

// Create a new  account
func addAiModel(owner solana.PublicKey, name string) {
	nameHash := sha256.Sum256([]byte(name))
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("ai_model"),
			owner[:],
			nameHash[:],
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
	m := new(distri_ai.AiModel)
	if err := m.UnmarshalWithDecoder(bin.NewBorshDecoder(resp.Value.Data.GetBinary())); err != nil {
		return
	}

	aiModel := buildAiModelModel(m)
	if dbResult := common.Db.Create(&aiModel); dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}

	heat := model.AiModelHeat{Owner: aiModel.Owner, Name: aiModel.Name}
	if dbResult := common.Db.Create(&heat); dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}

// remove account
func removeAiModel(owner solana.PublicKey, name string) {
	dbResult := common.Db.
		Where("owner = ?", owner.String()).
		Where("name = ?", name).
		Delete(&model.AiModel{})
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}

	dbResult = common.Db.
		Where("owner = ?", owner.String()).
		Where("name = ?", name).
		Delete(&model.AiModelHeat{})
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}

	dbResult = common.Db.
		Where("owner = ?", owner.String()).
		Where("name = ?", name).
		Delete(&model.AiModelLike{})
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
	}
}
