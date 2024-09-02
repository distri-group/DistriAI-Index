package chain

import (
	"context"
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/utils/logs"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
	"slices"
	"strings"
	"time"
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

// initChain initializes the blockchain-related configurations.
// This function sets the distribution program's public key and creates an RPC client for further blockchain operations.
func initChain() {
	// Converts the chain's program ID from Base58 to PublicKey, ensuring the correct program ID is configured.
	distriProgramID = solana.MustPublicKeyFromBase58(common.Conf.Chain.ProgramId)
	distri_ai.SetProgramID(distriProgramID)
	// Creates a new RPC client based on the RPC address in the configuration file, for communication with the blockchain.
	rpcClient = rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		common.Conf.Chain.Rpc,
		rate.Every(time.Second),
		10,
	))
	// Calls the initSolana function to perform additional initialization tasks for the Solana chain.
	initSolana()
}

func Sync() {
	initChain()
	//sync machine、order、reward、rewardMachine data from chain
	fetchAll()
}

// fetchAll retrieves all program accounts data and processes them.
func fetchAll() {
	// Fetch program accounts data using the default context and confirmed commitment level
	out, err := rpcClient.GetProgramAccountsWithOpts(
		context.TODO(),
		distriProgramID,
		&rpc.GetProgramAccountsOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	// If an error occurs during retrieval, log the error and terminate the function
	if err != nil {
		logs.Error(fmt.Sprintf("GetProgramAccounts error: %s \n", err))
		return
	}
	// Call handling functions to process fetched account data for AI models, datasets, machines, orders, rewards, and reward machines
	fetchAllAiModel(out)
	fetchAllDataset(out)
	fetchAllMachine(out)
	fetchAllOrder(out)
	fetchAllReward(out)
	fetchAllRewardMachine(out)
}

// HandleEventLogs processes a slice of event logs, identifying specific instructions and data to execute corresponding operations.
// It iterates through the logs, extracts instruction and data, then performs actions based on the instruction type.
func HandleEventLogs(eventLogs []string) {
	spew.Dump(eventLogs) // Debug print of incoming event logs
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
		// Extract data once the instruction is identified
		if after, found := strings.CutPrefix(l, _Data); found {
			data = after
			break
		}
	}
	// If neither instruction nor data is found, exit early
	if instruction == "" || data == "" {
		return
	}
	logs.Info(fmt.Sprintf("[Webhook] Receive instruction: %s \n", instruction)) // Log received instruction

	// Switch statement to handle various instruction cases
	switch instruction {
	// Cases handle different types of machine, task, reward, order, AI model, and dataset events
	// Each case decodes the event data, then performs an action such as adding, removing, updating records, etc.
	// ...

	// Default or additional cases would be placed here following the pattern above
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
		var event TaskEvent
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
	// Log the handled instruction upon completion
	logs.Info(fmt.Sprintf("[Webhook] Handle instruction: %s \n", instruction))
}
