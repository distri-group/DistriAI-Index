package chain

import (
	"context"
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
	faucetPrivateKey solana.PrivateKey
	faucetPublicKey  solana.PublicKey
	dist             solana.PublicKey
	distDecimals     uint8
	distFaucetAmount uint64
)

func initSolana() {
	faucetPrivateKey = solana.MustPrivateKeyFromBase58(common.Conf.Chain.FaucetPrivateKey)
	faucetPublicKey = faucetPrivateKey.PublicKey()
	dist = solana.MustPublicKeyFromBase58(common.Conf.Chain.Dist)
	distDecimals = common.Conf.Chain.DistDecimals
	distFaucetAmount = common.Conf.Chain.DistFaucetAmount * uint64(math.Pow10(int(distDecimals)))
}

func FaucetDist(publicKey solana.PublicKey) (string, error) {
	faucetAta, _, err := solana.FindAssociatedTokenAddress(faucetPublicKey, dist)
	if err != nil {
		logs.Error(fmt.Sprintf("error finding associated token address: %s \n", err))
		return "", fmt.Errorf("error finding associated token address: %s \n", err)
	}
	receiverAta, _, err := solana.FindAssociatedTokenAddress(publicKey, dist)
	if err != nil {
		logs.Error(fmt.Sprintf("error finding associated token address: %s \n", err))
		return "", fmt.Errorf("error finding associated token address: %s \n", err)
	}

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

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		logs.Error(fmt.Sprintf("Creating transaction error: %s \n", err))
		return "", fmt.Errorf("error creating transaction: %s \n", err)
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(faucetPublicKey),
	)

	if err != nil {
		logs.Error(fmt.Sprintf("Creating transaction error: %s \n", err))
		return "", fmt.Errorf("error creating transaction: %s \n", err)
	}

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

	logs.Info(fmt.Sprintf("[FaucetDist] Tx confirmed : %v", sig.String()))

	return sig.String(), nil
}

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
