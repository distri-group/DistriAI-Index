package handlers

import (
	"context"
	"distriai-index-solana/chain"
	"distriai-index-solana/common"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
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
		logs.Warn(fmt.Sprintf("Faucet parameter missing error : %s \n", err))
		resp.Fail(c, "Parameter missing")
		return
	}

	reqIp := c.ClientIP()
	logs.Info(fmt.Sprintf("Faucet request IP : %s \n", reqIp))
	// Check request Ip in Redis
	ctx := context.Background()
	exists, err := common.Rdb.Exists(ctx, reqIp).Result()
	if err == nil && exists == 1 {
		logs.Warn(fmt.Sprintf("Too many requests from this IP. IP : %s", reqIp))
		resp.Fail(c, "You have requested too many airdrops. Wait 24 hours for a refill.")
		return
	}

	publicKey, err := solana.PublicKeyFromBase58(req.Account)
	if err != nil {
		logs.Warn(fmt.Sprintf("Invalid account : %s \n", err))
		resp.Fail(c, "Invalid account")
		return
	}
	txHash, err := chain.FaucetDist(publicKey)
	if err != nil {
		logs.Warn(fmt.Sprintf("Transaction error : %s \n", err))
		resp.Fail(c, "Transaction fail")
		return
	}

	// Save request Ip in Redis
	err = common.Rdb.SetEx(ctx, reqIp, req.Account, time.Hour*24).Err()
	if err != nil {
		logs.Error(fmt.Sprintf("Redis error : %s \n", err))
		resp.Fail(c, "Redis error")
		return
	}

	response := FaucetResp{TxHash: txHash}
	resp.Success(c, response)
}
