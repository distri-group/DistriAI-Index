package chain

import (
	"context"
	"distriai-index-solana/common"
	"distriai-index-solana/utils/logs"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"slices"
	"strings"
)

const (
	// Event
	_Instruction    = "Program log: Instruction: "
	_Data           = "Program data: "
	_AddMachine     = "AddMachine"
	_RemoveMachine  = "RemoveMachine"
	_MakeOffer      = "MakeOffer"
	_CancelOffer    = "CancelOffer"
	_SubmitTask     = "SubmitTask"
	_Claim          = "Claim"
	_PlaceOrder     = "PlaceOrder"
	_RenewOrder     = "RenewOrder"
	_StartOrder     = "StartOrder"
	_RefundOrder    = "RefundOrder"
	_OrderCompleted = "OrderCompleted"
	_OrderFailed    = "OrderFailed"
	_RemoveOrder    = "RemoveOrder"
	_CreateAiModel  = "CreateAiModel"
	_RemoveAiModel  = "RemoveAiModel"
	_CreateDataset  = "CreateDataset"
	_RemoveDataset  = "RemoveDataset"
)

var (
	distriProgramID    solana.PublicKey
	rpcClient          *rpc.Client
	distriInstructions = []string{
		_AddMachine, _RemoveMachine, _MakeOffer, _CancelOffer, _SubmitTask, _Claim,
		_PlaceOrder, _RenewOrder, _StartOrder, _RefundOrder, _OrderCompleted, _OrderFailed, _RemoveOrder,
		_CreateAiModel, _RemoveAiModel, _CreateDataset, _RemoveDataset,
	}
)

func initChain() {
	distriProgramID = solana.MustPublicKeyFromBase58(common.Conf.Chain.ProgramId)
	rpcClient = rpc.New(common.Conf.Chain.Rpc)
	initSolana()
}

func Sync() {
	initChain()
	//sync machine、order、reward、rewardMachine data from chain
	fetchAll()
}

func fetchAll() {
	out, err := rpcClient.GetProgramAccountsWithOpts(
		context.TODO(),
		distriProgramID,
		&rpc.GetProgramAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		logs.Error(fmt.Sprintf("GetProgramAccounts error: %s \n", err))
		return
	}

	fetchAllAiModel(out)
	fetchAllDataset(out)
	fetchAllMachine(out)
	fetchAllOrder(out)
	fetchAllReward(out)
	fetchAllRewardMachine(out)
}

func HandleEventLogs(eventLogs []string) {
	spew.Dump(eventLogs)
	var instruction, data string
	for _, l := range eventLogs {
		// Find first HashrateMarket Instruction in event
		if instruction == "" {
			if after, found := strings.CutPrefix(l, _Instruction); found {
				if i := slices.Index(distriInstructions, after); i >= 0 {
					instruction = after
					continue
				}
			}
		}
		if after, found := strings.CutPrefix(l, _Data); found {
			data = after
			break
		}
	}
	if instruction == "" || data == "" {
		return
	}
	logs.Info(fmt.Sprintf("[Webhook] Receive instruction: %s \n", instruction))

	switch instruction {
	case _AddMachine:
		var event MachineEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		addMachine(event.Owner, event.Uuid)
	case _RemoveMachine:
		var event MachineEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		removeMachine(event.Owner, event.Uuid)
	case _MakeOffer, _CancelOffer:
		var event MachineEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		updateMachine(event.Owner, event.Uuid)
	case _SubmitTask:
		var event RewardEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		saveReward(event.Period)
		saveRewardMachine(event.Period, event.Owner, event.MachineId)
	case _Claim:
		var event RewardEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		updateMachine(event.Owner, event.MachineId)
		saveRewardMachine(event.Period, event.Owner, event.MachineId)
	case _PlaceOrder:
		var event OrderEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		updateMachine(event.Seller, event.MachineId)
		addOrder(event.OrderId, event.Buyer)
	case _RenewOrder, _StartOrder:
		var event OrderEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		updateOrder(event.OrderId, event.Buyer)
	case _RefundOrder, _OrderCompleted, _OrderFailed:
		var event OrderEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		updateMachine(event.Seller, event.MachineId)
		updateOrder(event.OrderId, event.Buyer)
	case _RemoveOrder:
		var event OrderEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		removeOrder(event.OrderId)
	case _CreateAiModel:
		var event AiModelEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		addAiModel(event.Owner, event.Name)
	case _RemoveAiModel:
		var event AiModelEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		removeAiModel(event.Owner, event.Name)
	case _CreateDataset:
		var event DatasetEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		addDataset(event.Owner, event.Name)
	case _RemoveDataset:
		var event DatasetEvent
		if decodeDistriEvent(data, &event) != nil {
			return
		}
		removeDataset(event.Owner, event.Name)
	}

	logs.Info(fmt.Sprintf("[Webhook] Handle instruction: %s \n", instruction))
}
