package handlers

import (
	"distriai-index-solana/chain"
	"distriai-index-solana/utils/resp"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
	"log"
)

type FaucetReq struct {
	Account string `binding:"required"`
}

type FaucetResp struct {
	TxHash string
}

func Faucet(context *gin.Context) {
	var req FaucetReq
	err := context.ShouldBindJSON(&req)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}

	publicKey, err := solana.PublicKeyFromBase58(req.Account)
	if err != nil {
		resp.Fail(context, "Invalid account")
		return
	}
	txHash, err := chain.FaucetDist(publicKey)
	if err != nil {
		log.Printf("Transaction error : %v \n", err)
		resp.Fail(context, "Transaction fail")
		return
	}

	response := FaucetResp{TxHash: txHash}
	resp.Success(context, response)
}
