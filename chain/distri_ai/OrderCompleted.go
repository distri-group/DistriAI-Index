// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package distri_ai

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// OrderCompleted is the `orderCompleted` instruction.
type OrderCompleted struct {
	Metadata *string
	Score    *uint8

	// [0] = [WRITE] machine
	//
	// [1] = [WRITE] order
	//
	// [2] = [WRITE, SIGNER] seller
	//
	// [3] = [WRITE] sellerAta
	//
	// [4] = [WRITE] model1OwnerAta
	//
	// [5] = [WRITE] model2OwnerAta
	//
	// [6] = [WRITE] model3OwnerAta
	//
	// [7] = [WRITE] model4OwnerAta
	//
	// [8] = [WRITE] model5OwnerAta
	//
	// [9] = [WRITE] statisticsSeller
	//
	// [10] = [WRITE] statisticsModel1Owner
	//
	// [11] = [WRITE] statisticsModel2Owner
	//
	// [12] = [WRITE] statisticsModel3Owner
	//
	// [13] = [WRITE] statisticsModel4Owner
	//
	// [14] = [WRITE] statisticsModel5Owner
	//
	// [15] = [WRITE] vault
	//
	// [16] = [] mint
	//
	// [17] = [] tokenProgram
	//
	// [18] = [] associatedTokenProgram
	//
	// [19] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewOrderCompletedInstructionBuilder creates a new `OrderCompleted` instruction builder.
func NewOrderCompletedInstructionBuilder() *OrderCompleted {
	nd := &OrderCompleted{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 20),
	}
	return nd
}

// SetMetadata sets the "metadata" parameter.
func (inst *OrderCompleted) SetMetadata(metadata string) *OrderCompleted {
	inst.Metadata = &metadata
	return inst
}

// SetScore sets the "score" parameter.
func (inst *OrderCompleted) SetScore(score uint8) *OrderCompleted {
	inst.Score = &score
	return inst
}

// SetMachineAccount sets the "machine" account.
func (inst *OrderCompleted) SetMachineAccount(machine ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(machine).WRITE()
	return inst
}

// GetMachineAccount gets the "machine" account.
func (inst *OrderCompleted) GetMachineAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetOrderAccount sets the "order" account.
func (inst *OrderCompleted) SetOrderAccount(order ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(order).WRITE()
	return inst
}

// GetOrderAccount gets the "order" account.
func (inst *OrderCompleted) GetOrderAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSellerAccount sets the "seller" account.
func (inst *OrderCompleted) SetSellerAccount(seller ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(seller).WRITE().SIGNER()
	return inst
}

// GetSellerAccount gets the "seller" account.
func (inst *OrderCompleted) GetSellerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSellerAtaAccount sets the "sellerAta" account.
func (inst *OrderCompleted) SetSellerAtaAccount(sellerAta ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(sellerAta).WRITE()
	return inst
}

// GetSellerAtaAccount gets the "sellerAta" account.
func (inst *OrderCompleted) GetSellerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetModel1OwnerAtaAccount sets the "model1OwnerAta" account.
func (inst *OrderCompleted) SetModel1OwnerAtaAccount(model1OwnerAta ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(model1OwnerAta).WRITE()
	return inst
}

// GetModel1OwnerAtaAccount gets the "model1OwnerAta" account.
func (inst *OrderCompleted) GetModel1OwnerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetModel2OwnerAtaAccount sets the "model2OwnerAta" account.
func (inst *OrderCompleted) SetModel2OwnerAtaAccount(model2OwnerAta ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(model2OwnerAta).WRITE()
	return inst
}

// GetModel2OwnerAtaAccount gets the "model2OwnerAta" account.
func (inst *OrderCompleted) GetModel2OwnerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetModel3OwnerAtaAccount sets the "model3OwnerAta" account.
func (inst *OrderCompleted) SetModel3OwnerAtaAccount(model3OwnerAta ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(model3OwnerAta).WRITE()
	return inst
}

// GetModel3OwnerAtaAccount gets the "model3OwnerAta" account.
func (inst *OrderCompleted) GetModel3OwnerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetModel4OwnerAtaAccount sets the "model4OwnerAta" account.
func (inst *OrderCompleted) SetModel4OwnerAtaAccount(model4OwnerAta ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(model4OwnerAta).WRITE()
	return inst
}

// GetModel4OwnerAtaAccount gets the "model4OwnerAta" account.
func (inst *OrderCompleted) GetModel4OwnerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetModel5OwnerAtaAccount sets the "model5OwnerAta" account.
func (inst *OrderCompleted) SetModel5OwnerAtaAccount(model5OwnerAta ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(model5OwnerAta).WRITE()
	return inst
}

// GetModel5OwnerAtaAccount gets the "model5OwnerAta" account.
func (inst *OrderCompleted) GetModel5OwnerAtaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetStatisticsSellerAccount sets the "statisticsSeller" account.
func (inst *OrderCompleted) SetStatisticsSellerAccount(statisticsSeller ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(statisticsSeller).WRITE()
	return inst
}

// GetStatisticsSellerAccount gets the "statisticsSeller" account.
func (inst *OrderCompleted) GetStatisticsSellerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetStatisticsModel1OwnerAccount sets the "statisticsModel1Owner" account.
func (inst *OrderCompleted) SetStatisticsModel1OwnerAccount(statisticsModel1Owner ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(statisticsModel1Owner).WRITE()
	return inst
}

// GetStatisticsModel1OwnerAccount gets the "statisticsModel1Owner" account.
func (inst *OrderCompleted) GetStatisticsModel1OwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetStatisticsModel2OwnerAccount sets the "statisticsModel2Owner" account.
func (inst *OrderCompleted) SetStatisticsModel2OwnerAccount(statisticsModel2Owner ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(statisticsModel2Owner).WRITE()
	return inst
}

// GetStatisticsModel2OwnerAccount gets the "statisticsModel2Owner" account.
func (inst *OrderCompleted) GetStatisticsModel2OwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

// SetStatisticsModel3OwnerAccount sets the "statisticsModel3Owner" account.
func (inst *OrderCompleted) SetStatisticsModel3OwnerAccount(statisticsModel3Owner ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[12] = ag_solanago.Meta(statisticsModel3Owner).WRITE()
	return inst
}

// GetStatisticsModel3OwnerAccount gets the "statisticsModel3Owner" account.
func (inst *OrderCompleted) GetStatisticsModel3OwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(12)
}

// SetStatisticsModel4OwnerAccount sets the "statisticsModel4Owner" account.
func (inst *OrderCompleted) SetStatisticsModel4OwnerAccount(statisticsModel4Owner ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[13] = ag_solanago.Meta(statisticsModel4Owner).WRITE()
	return inst
}

// GetStatisticsModel4OwnerAccount gets the "statisticsModel4Owner" account.
func (inst *OrderCompleted) GetStatisticsModel4OwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(13)
}

// SetStatisticsModel5OwnerAccount sets the "statisticsModel5Owner" account.
func (inst *OrderCompleted) SetStatisticsModel5OwnerAccount(statisticsModel5Owner ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[14] = ag_solanago.Meta(statisticsModel5Owner).WRITE()
	return inst
}

// GetStatisticsModel5OwnerAccount gets the "statisticsModel5Owner" account.
func (inst *OrderCompleted) GetStatisticsModel5OwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(14)
}

// SetVaultAccount sets the "vault" account.
func (inst *OrderCompleted) SetVaultAccount(vault ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[15] = ag_solanago.Meta(vault).WRITE()
	return inst
}

// GetVaultAccount gets the "vault" account.
func (inst *OrderCompleted) GetVaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(15)
}

// SetMintAccount sets the "mint" account.
func (inst *OrderCompleted) SetMintAccount(mint ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[16] = ag_solanago.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
func (inst *OrderCompleted) GetMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(16)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *OrderCompleted) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[17] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *OrderCompleted) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(17)
}

// SetAssociatedTokenProgramAccount sets the "associatedTokenProgram" account.
func (inst *OrderCompleted) SetAssociatedTokenProgramAccount(associatedTokenProgram ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[18] = ag_solanago.Meta(associatedTokenProgram)
	return inst
}

// GetAssociatedTokenProgramAccount gets the "associatedTokenProgram" account.
func (inst *OrderCompleted) GetAssociatedTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(18)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *OrderCompleted) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *OrderCompleted {
	inst.AccountMetaSlice[19] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *OrderCompleted) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(19)
}

func (inst OrderCompleted) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_OrderCompleted,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst OrderCompleted) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *OrderCompleted) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Metadata == nil {
			return errors.New("Metadata parameter is not set")
		}
		if inst.Score == nil {
			return errors.New("Score parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Machine is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Order is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Seller is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SellerAta is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.Model1OwnerAta is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.Model2OwnerAta is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.Model3OwnerAta is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.Model4OwnerAta is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.Model5OwnerAta is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.StatisticsSeller is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.StatisticsModel1Owner is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.StatisticsModel2Owner is not set")
		}
		if inst.AccountMetaSlice[12] == nil {
			return errors.New("accounts.StatisticsModel3Owner is not set")
		}
		if inst.AccountMetaSlice[13] == nil {
			return errors.New("accounts.StatisticsModel4Owner is not set")
		}
		if inst.AccountMetaSlice[14] == nil {
			return errors.New("accounts.StatisticsModel5Owner is not set")
		}
		if inst.AccountMetaSlice[15] == nil {
			return errors.New("accounts.Vault is not set")
		}
		if inst.AccountMetaSlice[16] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[17] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[18] == nil {
			return errors.New("accounts.AssociatedTokenProgram is not set")
		}
		if inst.AccountMetaSlice[19] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *OrderCompleted) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("OrderCompleted")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=2]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Metadata", *inst.Metadata))
						paramsBranch.Child(ag_format.Param("   Score", *inst.Score))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=20]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("               machine", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("                 order", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("                seller", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("             sellerAta", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("        model1OwnerAta", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("        model2OwnerAta", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("        model3OwnerAta", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("        model4OwnerAta", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("        model5OwnerAta", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("      statisticsSeller", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta(" statisticsModel1Owner", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta(" statisticsModel2Owner", inst.AccountMetaSlice.Get(11)))
						accountsBranch.Child(ag_format.Meta(" statisticsModel3Owner", inst.AccountMetaSlice.Get(12)))
						accountsBranch.Child(ag_format.Meta(" statisticsModel4Owner", inst.AccountMetaSlice.Get(13)))
						accountsBranch.Child(ag_format.Meta(" statisticsModel5Owner", inst.AccountMetaSlice.Get(14)))
						accountsBranch.Child(ag_format.Meta("                 vault", inst.AccountMetaSlice.Get(15)))
						accountsBranch.Child(ag_format.Meta("                  mint", inst.AccountMetaSlice.Get(16)))
						accountsBranch.Child(ag_format.Meta("          tokenProgram", inst.AccountMetaSlice.Get(17)))
						accountsBranch.Child(ag_format.Meta("associatedTokenProgram", inst.AccountMetaSlice.Get(18)))
						accountsBranch.Child(ag_format.Meta("         systemProgram", inst.AccountMetaSlice.Get(19)))
					})
				})
		})
}

func (obj OrderCompleted) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Metadata` param:
	err = encoder.Encode(obj.Metadata)
	if err != nil {
		return err
	}
	// Serialize `Score` param:
	err = encoder.Encode(obj.Score)
	if err != nil {
		return err
	}
	return nil
}
func (obj *OrderCompleted) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Metadata`:
	err = decoder.Decode(&obj.Metadata)
	if err != nil {
		return err
	}
	// Deserialize `Score`:
	err = decoder.Decode(&obj.Score)
	if err != nil {
		return err
	}
	return nil
}

// NewOrderCompletedInstruction declares a new OrderCompleted instruction with the provided parameters and accounts.
func NewOrderCompletedInstruction(
	// Parameters:
	metadata string,
	score uint8,
	// Accounts:
	machine ag_solanago.PublicKey,
	order ag_solanago.PublicKey,
	seller ag_solanago.PublicKey,
	sellerAta ag_solanago.PublicKey,
	model1OwnerAta ag_solanago.PublicKey,
	model2OwnerAta ag_solanago.PublicKey,
	model3OwnerAta ag_solanago.PublicKey,
	model4OwnerAta ag_solanago.PublicKey,
	model5OwnerAta ag_solanago.PublicKey,
	statisticsSeller ag_solanago.PublicKey,
	statisticsModel1Owner ag_solanago.PublicKey,
	statisticsModel2Owner ag_solanago.PublicKey,
	statisticsModel3Owner ag_solanago.PublicKey,
	statisticsModel4Owner ag_solanago.PublicKey,
	statisticsModel5Owner ag_solanago.PublicKey,
	vault ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	associatedTokenProgram ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *OrderCompleted {
	return NewOrderCompletedInstructionBuilder().
		SetMetadata(metadata).
		SetScore(score).
		SetMachineAccount(machine).
		SetOrderAccount(order).
		SetSellerAccount(seller).
		SetSellerAtaAccount(sellerAta).
		SetModel1OwnerAtaAccount(model1OwnerAta).
		SetModel2OwnerAtaAccount(model2OwnerAta).
		SetModel3OwnerAtaAccount(model3OwnerAta).
		SetModel4OwnerAtaAccount(model4OwnerAta).
		SetModel5OwnerAtaAccount(model5OwnerAta).
		SetStatisticsSellerAccount(statisticsSeller).
		SetStatisticsModel1OwnerAccount(statisticsModel1Owner).
		SetStatisticsModel2OwnerAccount(statisticsModel2Owner).
		SetStatisticsModel3OwnerAccount(statisticsModel3Owner).
		SetStatisticsModel4OwnerAccount(statisticsModel4Owner).
		SetStatisticsModel5OwnerAccount(statisticsModel5Owner).
		SetVaultAccount(vault).
		SetMintAccount(mint).
		SetTokenProgramAccount(tokenProgram).
		SetAssociatedTokenProgramAccount(associatedTokenProgram).
		SetSystemProgramAccount(systemProgram)
}
