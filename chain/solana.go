package chain

import (
	"context"
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/utils/logs"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"math"
	"time"
)

var (
	adminPrivateKey  solana.PrivateKey
	adminPublicKey   solana.PublicKey
	faucetPrivateKey solana.PrivateKey
	faucetPublicKey  solana.PublicKey
	dist             solana.PublicKey
	distDecimals     uint8
	distFaucetAmount uint64
)

func initSolana() {
	adminPrivateKey = solana.MustPrivateKeyFromBase58(common.Conf.Chain.AdminPrivateKey)
	adminPublicKey = adminPrivateKey.PublicKey()
	faucetPrivateKey = solana.MustPrivateKeyFromBase58(common.Conf.Chain.FaucetPrivateKey)
	faucetPublicKey = faucetPrivateKey.PublicKey()
	dist = solana.MustPublicKeyFromBase58(common.Conf.Chain.Dist)
	distDecimals = common.Conf.Chain.DistDecimals
	distFaucetAmount = common.Conf.Chain.DistFaucetAmount * uint64(math.Pow10(int(distDecimals)))
}

// FaucetDist facilitates the distribution of a token from a faucet account to a recipient's associated token account on Solana blockchain.
func FaucetDist(publicKey solana.PublicKey) (string, error) {
	// Find the associated token address for the faucet account.
	faucetAta, _, err := solana.FindAssociatedTokenAddress(faucetPublicKey, dist)
	if err != nil {
		logs.Error(fmt.Sprintf("error finding associated token address: %s \n", err))
		return "", fmt.Errorf("error finding associated token address: %s \n", err)
	}
	// Find the associated token address for the recipient's public key.
	receiverAta, _, err := solana.FindAssociatedTokenAddress(publicKey, dist)
	if err != nil {
		logs.Error(fmt.Sprintf("error finding associated token address: %s \n", err))
		return "", fmt.Errorf("error finding associated token address: %s \n", err)
	}

	// Prepare a list of instructions to be included in the transaction.
	var instructions []solana.Instruction
	_, err = rpcClient.GetAccountInfo(context.TODO(), receiverAta)
	if errors.Is(err, rpc.ErrNotFound) {
		instructions = append(instructions,
			associatedtokenaccount.NewCreateInstruction(
				faucetPublicKey,
				publicKey,
				dist,
			).Build(),
		)
	}
	
	// Add a transfer instruction to transfer tokens from the faucet to the recipient.
	instructions = append(instructions,
		token.NewTransferCheckedInstruction(
			distFaucetAmount,
			distDecimals,
			faucetAta,
			dist,
			receiverAta,
			faucetPublicKey,
			[]solana.PublicKey{},
		).Build(),
	)

	// Get the recent blockhash for the transaction.
	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		logs.Error(fmt.Sprintf("Creating transaction error: %s \n", err))
		return "", fmt.Errorf("error creating transaction: %s \n", err)
	}

	// Create a new transaction using the collected instructions and the recent blockhash.
	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(faucetPublicKey),
	)

	if err != nil {
		logs.Error(fmt.Sprintf("Creating transaction error: %s \n", err))
		return "", fmt.Errorf("error creating transaction: %s \n", err)
	}

	// Sign the transaction with the faucet's private key.
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if faucetPublicKey.Equals(key) {
				return &faucetPrivateKey
			}
			return nil
		},
	)
	if err != nil {
		logs.Error(fmt.Sprintf("Signing transaction error: %s \n", err))
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	// Display the transaction details using spew (debugging tool).
	spew.Dump(tx)

	// Send the signed transaction to the Solana blockchain.
	sig, err := rpcClient.SendTransaction(context.TODO(), tx)
	if err != nil {
		spew.Dump(err)
		logs.Error(fmt.Sprintf("Sending transaction error: %s \n", err))
		return "", fmt.Errorf("error sending transaction: %v", err)
	}
	if _, err := waitForConfirm(sig); err != nil {
		logs.Error(fmt.Sprintf("Tx waitForConfirm error: %s \n", err))
		return "", fmt.Errorf("error sending transaction: %v", err)
	}

	logs.Info(fmt.Sprintf("[FaucetDist] Tx confirmed : %v", sig.String()))

	return sig.String(), nil
}

// Function to report AI model dataset reward for a given owner and amount.
func ReportAiModelDatasetReward(owner solana.PublicKey, amount uint64) (string, error) {
	address, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("statistics"),
			owner.Bytes(),
		},
		distriProgramID,
	)
	if err != nil {
		logs.Error(fmt.Sprintf("FindProgramAddress error: %s \n", err))
		return "", fmt.Errorf("FindProgramAddress error: %s \n", err)
	}

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		logs.Error(fmt.Sprintf("Creating transaction error: %s \n", err))
		return "", fmt.Errorf("error creating transaction: %s \n", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			distri_ai.NewReportAiModelDatasetRewardInstruction(
				amount,
				address,
				adminPublicKey,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(adminPublicKey),
	)

	if err != nil {
		logs.Error(fmt.Sprintf("Creating transaction error: %s \n", err))
		return "", fmt.Errorf("error creating transaction: %s \n", err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if adminPublicKey.Equals(key) {
				return &adminPrivateKey
			}
			return nil
		},
	)
	if err != nil {
		logs.Error(fmt.Sprintf("Signing transaction error: %s \n", err))
		return "", fmt.Errorf("error signing transaction: %v", err)
	}

	spew.Dump(tx)

	sig, err := rpcClient.SendTransaction(context.TODO(), tx)
	if err != nil {
		spew.Dump(err)
		logs.Error(fmt.Sprintf("Sending transaction error: %s \n", err))
		return "", fmt.Errorf("error sending transaction: %v", err)
	}
	if _, err := waitForConfirm(sig); err != nil {
		logs.Error(fmt.Sprintf("Tx waitForConfirm error: %s \n", err))
		return "", fmt.Errorf("error sending transaction: %v", err)
	}

	logs.Info(fmt.Sprintf("[ReportAiModelDatasetReward] Tx confirmed : %v", sig.String()))

	return sig.String(), nil
}

// waitForConfirm is a function that waits for a Solana transaction signature to be confirmed.
// It takes a Solana.Signature as input and returns a boolean indicating if the transaction is confirmed,
// and an error if one occurred during the process.
func waitForConfirm(sig solana.Signature) (confirmed bool, err error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	tryTimes := 0
	for {
		select {
		case <-ticker.C:
			_, err := rpcClient.GetTransaction(
				context.TODO(),
				sig,
				&rpc.GetTransactionOpts{
					Commitment: rpc.CommitmentConfirmed,
				},
			)
			if err == nil {
				return true, nil
			}
			tryTimes++
			if tryTimes > 20 {
				return false, err
			}
		}
	}
}
