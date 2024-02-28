// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package distri_ai

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// SubmitTask is the `submitTask` instruction.
type SubmitTask struct {
	Uuid     *[16]uint8
	Period   *uint32
	Metadata *string

	// [0] = [WRITE] machine
	//
	// [1] = [WRITE] task
	//
	// [2] = [WRITE] reward
	//
	// [3] = [WRITE] rewardMachine
	//
	// [4] = [WRITE, SIGNER] owner
	//
	// [5] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewSubmitTaskInstructionBuilder creates a new `SubmitTask` instruction builder.
func NewSubmitTaskInstructionBuilder() *SubmitTask {
	nd := &SubmitTask{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 6),
	}
	return nd
}

// SetUuid sets the "uuid" parameter.
func (inst *SubmitTask) SetUuid(uuid [16]uint8) *SubmitTask {
	inst.Uuid = &uuid
	return inst
}

// SetPeriod sets the "period" parameter.
func (inst *SubmitTask) SetPeriod(period uint32) *SubmitTask {
	inst.Period = &period
	return inst
}

// SetMetadata sets the "metadata" parameter.
func (inst *SubmitTask) SetMetadata(metadata string) *SubmitTask {
	inst.Metadata = &metadata
	return inst
}

// SetMachineAccount sets the "machine" account.
func (inst *SubmitTask) SetMachineAccount(machine ag_solanago.PublicKey) *SubmitTask {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(machine).WRITE()
	return inst
}

// GetMachineAccount gets the "machine" account.
func (inst *SubmitTask) GetMachineAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetTaskAccount sets the "task" account.
func (inst *SubmitTask) SetTaskAccount(task ag_solanago.PublicKey) *SubmitTask {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(task).WRITE()
	return inst
}

// GetTaskAccount gets the "task" account.
func (inst *SubmitTask) GetTaskAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetRewardAccount sets the "reward" account.
func (inst *SubmitTask) SetRewardAccount(reward ag_solanago.PublicKey) *SubmitTask {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(reward).WRITE()
	return inst
}

// GetRewardAccount gets the "reward" account.
func (inst *SubmitTask) GetRewardAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetRewardMachineAccount sets the "rewardMachine" account.
func (inst *SubmitTask) SetRewardMachineAccount(rewardMachine ag_solanago.PublicKey) *SubmitTask {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(rewardMachine).WRITE()
	return inst
}

// GetRewardMachineAccount gets the "rewardMachine" account.
func (inst *SubmitTask) GetRewardMachineAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetOwnerAccount sets the "owner" account.
func (inst *SubmitTask) SetOwnerAccount(owner ag_solanago.PublicKey) *SubmitTask {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(owner).WRITE().SIGNER()
	return inst
}

// GetOwnerAccount gets the "owner" account.
func (inst *SubmitTask) GetOwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *SubmitTask) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *SubmitTask {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *SubmitTask) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

func (inst SubmitTask) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_SubmitTask,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SubmitTask) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SubmitTask) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Uuid == nil {
			return errors.New("Uuid parameter is not set")
		}
		if inst.Period == nil {
			return errors.New("Period parameter is not set")
		}
		if inst.Metadata == nil {
			return errors.New("Metadata parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Machine is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Task is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Reward is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.RewardMachine is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *SubmitTask) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SubmitTask")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=3]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("    Uuid", *inst.Uuid))
						paramsBranch.Child(ag_format.Param("  Period", *inst.Period))
						paramsBranch.Child(ag_format.Param("Metadata", *inst.Metadata))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=6]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("      machine", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("         task", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("       reward", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("rewardMachine", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("        owner", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(5)))
					})
				})
		})
}

func (obj SubmitTask) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Uuid` param:
	err = encoder.Encode(obj.Uuid)
	if err != nil {
		return err
	}
	// Serialize `Period` param:
	err = encoder.Encode(obj.Period)
	if err != nil {
		return err
	}
	// Serialize `Metadata` param:
	err = encoder.Encode(obj.Metadata)
	if err != nil {
		return err
	}
	return nil
}
func (obj *SubmitTask) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Uuid`:
	err = decoder.Decode(&obj.Uuid)
	if err != nil {
		return err
	}
	// Deserialize `Period`:
	err = decoder.Decode(&obj.Period)
	if err != nil {
		return err
	}
	// Deserialize `Metadata`:
	err = decoder.Decode(&obj.Metadata)
	if err != nil {
		return err
	}
	return nil
}

// NewSubmitTaskInstruction declares a new SubmitTask instruction with the provided parameters and accounts.
func NewSubmitTaskInstruction(
	// Parameters:
	uuid [16]uint8,
	period uint32,
	metadata string,
	// Accounts:
	machine ag_solanago.PublicKey,
	task ag_solanago.PublicKey,
	reward ag_solanago.PublicKey,
	rewardMachine ag_solanago.PublicKey,
	owner ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *SubmitTask {
	return NewSubmitTaskInstructionBuilder().
		SetUuid(uuid).
		SetPeriod(period).
		SetMetadata(metadata).
		SetMachineAccount(machine).
		SetTaskAccount(task).
		SetRewardAccount(reward).
		SetRewardMachineAccount(rewardMachine).
		SetOwnerAccount(owner).
		SetSystemProgramAccount(systemProgram)
}