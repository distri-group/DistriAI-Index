package handlers

import (
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func getAuthAccount(c *gin.Context) string {
	return c.MustGet("account").(string)
}

func getAccount(context *gin.Context) (string, error) {
	var header HttpHeader
	if err := context.ShouldBindHeader(&header); err != nil {
		logs.Warn(fmt.Sprintf("Header paramter missing: %s \n", err))
		resp.Fail(context, err.Error())
		return "", err
	}
	return header.Account, nil
}

const (
	// Period 0 start time: 2024-02-27 00:00:00 UTC
	genesisTime    int64 = 1708992000
	periodDuration int64 = 86400
)

func currentPeriod() uint32 {
	return uint32((time.Now().Unix() - genesisTime) / periodDuration)
}
