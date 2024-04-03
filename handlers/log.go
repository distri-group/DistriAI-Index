package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LogAddReq struct {
	OrderUuid string `binding:"required"`
	Content   string `binding:"required"`
}

func LogAdd(context *gin.Context) {
	var req LogAddReq
	err := context.ShouldBindJSON(&req)
	if err != nil {
		logs.Warn(fmt.Sprintf("Parameter missing,error:%s \n", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	log := &model.Log{OrderUuid: req.OrderUuid, Content: req.Content}
	dbResult := common.Db.Create(&log)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error,error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, "")
}

type LogListReq struct {
	OrderUuid string `binding:"required"`
	PageReq
}

type LogListResponse struct {
	List []model.Log
	PageResp
}

func LogList(context *gin.Context) {
	var req LogListReq
	err := context.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		logs.Warn(fmt.Sprintf("Parameter missing,error: %s \n", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	log := &model.Log{OrderUuid: req.OrderUuid}
	var response LogListResponse
	tx := common.Db.Model(&log).Where(&log)
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Order("created_at DESC").Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Warn(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}
