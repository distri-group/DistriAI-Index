package handlers

import (
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
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
