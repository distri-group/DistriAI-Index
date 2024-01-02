package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/resp"
	"github.com/gin-gonic/gin"
)

type LogAddReq struct {
	OrderUuid string `binding:"required"`
	Content   string `binding:"required"`
}

func LogAdd(context *gin.Context) {
	var req LogAddReq
	err := context.ShouldBindJSON(&req)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}

	log := &model.Log{OrderUuid: req.OrderUuid, Content: req.Content}
	dbResult := common.Db.Create(&log)
	if dbResult.Error != nil {
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
	err := context.ShouldBindJSON(&req)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}

	log := &model.Log{OrderUuid: req.OrderUuid}
	var response LogListResponse
	tx := common.Db.Model(&log).Where(&log)
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Order("created_at DESC").Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}
