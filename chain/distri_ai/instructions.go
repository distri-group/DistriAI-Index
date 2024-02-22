// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package distri_ai

import (
	"bytes"
	"fmt"
	ag_spew "github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_text "github.com/gagliardetto/solana-go/text"
	ag_treeout "github.com/gagliardetto/treeout"
)

var ProgramID ag_solanago.PublicKey

func SetProgramID(pubkey ag_solanago.PublicKey) {
	ProgramID = pubkey
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const ProgramName = "DistriAi"

func init() {
	if !ProgramID.IsZero() {
		ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
	}
}

var (
	Instruction_AddMachine = ag_binary.TypeID([8]byte{148, 26, 70, 80, 42, 110, 107, 230})

	Instruction_RemoveMachine = ag_binary.TypeID([8]byte{85, 41, 207, 236, 20, 250, 8, 97})

	Instruction_MakeOffer = ag_binary.TypeID([8]byte{214, 98, 97, 35, 59, 12, 44, 178})

	Instruction_CancelOffer = ag_binary.TypeID([8]byte{92, 203, 223, 40, 92, 89, 53, 119})

	Instruction_SubmitTask = ag_binary.TypeID([8]byte{148, 183, 26, 116, 107, 213, 118, 213})

	Instruction_Claim = ag_binary.TypeID([8]byte{62, 198, 214, 193, 213, 159, 108, 210})

	Instruction_PlaceOrder = ag_binary.TypeID([8]byte{51, 194, 155, 175, 109, 130, 96, 106})

	Instruction_RenewOrder = ag_binary.TypeID([8]byte{216, 180, 12, 76, 71, 44, 165, 151})

	Instruction_RefundOrder = ag_binary.TypeID([8]byte{164, 168, 47, 144, 154, 1, 241, 255})

	Instruction_OrderCompleted = ag_binary.TypeID([8]byte{60, 28, 38, 17, 211, 99, 139, 226})

	Instruction_OrderFailed = ag_binary.TypeID([8]byte{27, 173, 43, 153, 198, 108, 109, 66})

	Instruction_RemoveOrder = ag_binary.TypeID([8]byte{118, 116, 244, 40, 144, 211, 242, 51})
)

// InstructionIDToName returns the name of the instruction given its ID.
func InstructionIDToName(id ag_binary.TypeID) string {
	switch id {
	case Instruction_AddMachine:
		return "AddMachine"
	case Instruction_RemoveMachine:
		return "RemoveMachine"
	case Instruction_MakeOffer:
		return "MakeOffer"
	case Instruction_CancelOffer:
		return "CancelOffer"
	case Instruction_SubmitTask:
		return "SubmitTask"
	case Instruction_Claim:
		return "Claim"
	case Instruction_PlaceOrder:
		return "PlaceOrder"
	case Instruction_RenewOrder:
		return "RenewOrder"
	case Instruction_RefundOrder:
		return "RefundOrder"
	case Instruction_OrderCompleted:
		return "OrderCompleted"
	case Instruction_OrderFailed:
		return "OrderFailed"
	case Instruction_RemoveOrder:
		return "RemoveOrder"
	default:
		return ""
	}
}

type Instruction struct {
	ag_binary.BaseVariant
}

func (inst *Instruction) EncodeToTree(parent ag_treeout.Branches) {
	if enToTree, ok := inst.Impl.(ag_text.EncodableToTree); ok {
		enToTree.EncodeToTree(parent)
	} else {
		parent.Child(ag_spew.Sdump(inst))
	}
}

var InstructionImplDef = ag_binary.NewVariantDefinition(
	ag_binary.AnchorTypeIDEncoding,
	[]ag_binary.VariantType{
		{
			"add_machine", (*AddMachine)(nil),
		},
		{
			"remove_machine", (*RemoveMachine)(nil),
		},
		{
			"make_offer", (*MakeOffer)(nil),
		},
		{
			"cancel_offer", (*CancelOffer)(nil),
		},
		{
			"submit_task", (*SubmitTask)(nil),
		},
		{
			"claim", (*Claim)(nil),
		},
		{
			"place_order", (*PlaceOrder)(nil),
		},
		{
			"renew_order", (*RenewOrder)(nil),
		},
		{
			"refund_order", (*RefundOrder)(nil),
		},
		{
			"order_completed", (*OrderCompleted)(nil),
		},
		{
			"order_failed", (*OrderFailed)(nil),
		},
		{
			"remove_order", (*RemoveOrder)(nil),
		},
	},
)

func (inst *Instruction) ProgramID() ag_solanago.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() (out []*ag_solanago.AccountMeta) {
	return inst.Impl.(ag_solanago.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ag_binary.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) TextEncode(encoder *ag_text.Encoder, option *ag_text.Option) error {
	return encoder.Encode(inst.Impl, option)
}

func (inst *Instruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return inst.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionImplDef)
}

func (inst *Instruction) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.WriteBytes(inst.TypeID.Bytes(), false)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}

func registryDecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (*Instruction, error) {
	inst := new(Instruction)
	if err := ag_binary.NewBorshDecoder(data).Decode(inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction: %w", err)
	}
	if v, ok := inst.Impl.(ag_solanago.AccountsSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}
	return inst, nil
}
