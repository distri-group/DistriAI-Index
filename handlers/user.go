package handlers

import (
	"distriai-index-solana/middleware"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Account   string `binding:"required,len=44"`
	Signature string `binding:"required,len=88"`
}

func Login(context *gin.Context) {
	var req LoginReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	publicKey, err := solana.PublicKeyFromBase58(req.Account)
	if err != nil {
		logs.Warn(fmt.Sprintf("Invalid account : %s \n", err))
		resp.Fail(context, "Invalid account")
		return
	}
	signature, err := solana.SignatureFromBase58(req.Signature)
	if err != nil {
		logs.Warn(fmt.Sprintf("Invalid signature : %s \n", err))
		resp.Fail(context, "Invalid signature")
		return
	}
	if !publicKey.Verify([]byte(req.Account+"@distri.ai"), signature) {
		resp.Fail(context, "Invalid signature")
		return
	}

	token, _ := middleware.GenToken(req.Account)
	resp.Success(context, token)
}
