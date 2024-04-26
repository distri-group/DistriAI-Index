package chain

import (
	"context"
	"distriai-index-solana/common"
	"distriai-index-solana/utils/logs"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
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
)

var (
	distriProgramID    solana.PublicKey
	rpcClient          *rpc.Client
	wsClient           *ws.Client
	sub                *ws.LogSubscription
	distriInstructions = []string{
		_AddMachine, _RemoveMachine, _MakeOffer, _CancelOffer, _SubmitTask, _Claim,
		_PlaceOrder, _RenewOrder, _StartOrder, _RefundOrder, _OrderCompleted, _OrderFailed, _RemoveOrder,
	}
)

func initChain() {
	distriProgramID = solana.MustPublicKeyFromBase58(common.Conf.Chain.ProgramId)
	rpcClient = rpc.New(common.Conf.Chain.Rpc)

	initSolana()
}

func subLogs() {
	var err error
	wsClient, err = ws.Connect(context.Background(), common.Conf.Chain.Ws)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to '%s': %s", rpc.DevNet_WS, err))
	}
	sub, err = wsClient.LogsSubscribeMentions(
		distriProgramID,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		logs.Error(fmt.Sprintf("LogsSubscribe error: '%s' \n", err))
	}
}

func Sync() {
	initChain()
	//sync machine、order、reward、rewardMachine data from chain
	fetchAll()
	// subscribe smart contract events
	//go subEvents()
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

	fetchAllMachine(out)
	fetchAllOrder(out)
	fetchAllReward(out)
	fetchAllRewardMachine(out)
}

func subEvents() {
	subLogs()
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			logs.Error(fmt.Sprintf("SubEvents error: %v \n", err))
			subLogs()
			continue
		}

		eventLogs := got.Value.Logs
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
			continue
		}
		logs.Info(fmt.Sprintf("[WebSocket] Receive instruction: %s \n", instruction))

		switch instruction {
		case _AddMachine:
			event, err := decodeMachineEvent(data)
			if err != nil {
				continue
			}
			addMachine(event.Owner, event.Uuid)
		case _RemoveMachine:
			event, err := decodeMachineEvent(data)
			if err != nil {
				continue
			}
			removeMachine(event.Owner, event.Uuid)
		case _MakeOffer, _CancelOffer:
			event, err := decodeMachineEvent(data)
			if err != nil {
				continue
			}
			updateMachine(event.Owner, event.Uuid)
		case _SubmitTask:
			event, err := decodeTaskEvent(data)
			if err != nil {
				continue
			}
			saveReward(event.Period)
			saveRewardMachine(event.Period, event.Owner, event.MachineId)
		case _Claim:
			event, err := decodeRewardEvent(data)
			if err != nil {
				continue
			}
			updateMachine(event.Owner, event.MachineId)
			saveRewardMachine(event.Period, event.Owner, event.MachineId)
		case _PlaceOrder:
			event, err := decodeOrderEvent(data)
			if err != nil {
				continue
			}
			updateMachine(event.Seller, event.MachineId)
			addOrder(event.OrderId, event.Buyer)
		case _RenewOrder, _StartOrder:
			event, err := decodeOrderEvent(data)
			if err != nil {
				continue
			}
			updateOrder(event.OrderId, event.Buyer)
		case _RefundOrder, _OrderCompleted, _OrderFailed:
			event, err := decodeOrderEvent(data)
			if err != nil {
				continue
			}
			updateMachine(event.Seller, event.MachineId)
			updateOrder(event.OrderId, event.Buyer)
		case _RemoveOrder:
			event, err := decodeOrderEvent(data)
			if err != nil {
				continue
			}
			removeOrder(event.OrderId)
		}
	}
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
		event, err := decodeMachineEvent(data)
		if err != nil {
			return
		}
		addMachine(event.Owner, event.Uuid)
	case _RemoveMachine:
		event, err := decodeMachineEvent(data)
		if err != nil {
			return
		}
		removeMachine(event.Owner, event.Uuid)
	case _MakeOffer, _CancelOffer:
		event, err := decodeMachineEvent(data)
		if err != nil {
			return
		}
		updateMachine(event.Owner, event.Uuid)
	case _SubmitTask:
		event, err := decodeTaskEvent(data)
		if err != nil {
			return
		}
		saveReward(event.Period)
		saveRewardMachine(event.Period, event.Owner, event.MachineId)
	case _Claim:
		event, err := decodeRewardEvent(data)
		if err != nil {
			return
		}
		updateMachine(event.Owner, event.MachineId)
		saveRewardMachine(event.Period, event.Owner, event.MachineId)
	case _PlaceOrder:
		event, err := decodeOrderEvent(data)
		if err != nil {
			return
		}
		updateMachine(event.Seller, event.MachineId)
		addOrder(event.OrderId, event.Buyer)
	case _RenewOrder, _StartOrder:
		event, err := decodeOrderEvent(data)
		if err != nil {
			return
		}
		updateOrder(event.OrderId, event.Buyer)
	case _RefundOrder, _OrderCompleted, _OrderFailed:
		event, err := decodeOrderEvent(data)
		if err != nil {
			return
		}
		updateMachine(event.Seller, event.MachineId)
		updateOrder(event.OrderId, event.Buyer)
	case _RemoveOrder:
		event, err := decodeOrderEvent(data)
		if err != nil {
			return
		}
		removeOrder(event.OrderId)
	}

	logs.Info(fmt.Sprintf("[Webhook] Handle instruction: %s \n", instruction))
}

func decodeMachineEvent(data string) (*MachineEvent, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(bytes) < 8 {
		return nil, errors.New("data too short")
	}
	event := new(MachineEvent)
	if err := event.UnmarshalWithDecoder(bin.NewBorshDecoder(bytes[8:])); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeOrderEvent(data string) (*OrderEvent, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(bytes) < 8 {
		return nil, errors.New("data too short")
	}
	event := new(OrderEvent)
	if err := event.UnmarshalWithDecoder(bin.NewBorshDecoder(bytes[8:])); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeTaskEvent(data string) (*TaskEvent, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(bytes) < 8 {
		return nil, errors.New("data too short")
	}
	event := new(TaskEvent)
	if err := event.UnmarshalWithDecoder(bin.NewBorshDecoder(bytes[8:])); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeRewardEvent(data string) (*RewardEvent, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(bytes) < 8 {
		return nil, errors.New("data too short")
	}
	event := new(RewardEvent)
	if err := event.UnmarshalWithDecoder(bin.NewBorshDecoder(bytes[8:])); err != nil {
		return nil, err
	}
	return event, nil
}
