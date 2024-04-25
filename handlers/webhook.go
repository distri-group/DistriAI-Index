package handlers

import (
	"distriai-index-solana/chain"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Tx struct {
	Meta Meta `binding:"required" json:"meta"`
}

type Meta struct {
	LogMessages []string `binding:"required" json:"logMessages"`
}

func Webhook(context *gin.Context) {
	var req []Tx
	err := context.ShouldBindJSON(&req)
	if err != nil {
		logs.Warn(fmt.Sprintf("[Webhook] payload error:%s \n", err))
		resp.Fail(context, err.Error())
		return
	}

	for _, tx := range req {
		chain.HandleEventLogs(tx.Meta.LogMessages)
	}

	resp.Success(context, "")
}
