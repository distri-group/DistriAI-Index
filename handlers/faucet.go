package handlers

import (
	"context"
	"distriai-index-solana/chain"
	"distriai-index-solana/common"
	"distriai-index-solana/utils/resp"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type FaucetReq struct {
	Account string `binding:"required"`
}

type FaucetResp struct {
	TxHash string
}

func Faucet(c *gin.Context) {
	var req FaucetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.Fail(c, "Parameter missing")
		return
	}

	reqIp := c.ClientIP()
	log.Printf("Faucet request IP : %s \n", reqIp)
	// Check request Ip in Redis
	ctx := context.Background()
	exists, err := common.Rdb.Exists(ctx, reqIp).Result()
	if err == nil && exists == 1 {
		resp.Fail(c, "You have requested too many airdrops. Wait 24 hours for a refill.")
		return
	}

	publicKey, err := solana.PublicKeyFromBase58(req.Account)
	if err != nil {
		resp.Fail(c, "Invalid account")
		return
	}
	txHash, err := chain.FaucetDist(publicKey)
	if err != nil {
		log.Printf("Transaction error : %v \n", err)
		resp.Fail(c, "Transaction fail")
		return
	}

	// Save request Ip in Redis
	err = common.Rdb.SetEx(ctx, reqIp, req.Account, time.Hour*24).Err()
	if err != nil {
		resp.Fail(c, "Redis error")
		return
	}

	response := FaucetResp{TxHash: txHash}
	resp.Success(c, response)
}
