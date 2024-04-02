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

type OrderMineReq struct {
	Direction string
	Status    *uint8
	PageReq
}

type OrderListResponse struct {
	List []model.Order
	PageResp
}

func OrderMine(context *gin.Context) {
	account, err := getAccount(context)
	if err != nil {
		return
	}
	var req OrderMineReq
	err = context.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		logs.Warn(fmt.Sprintf("body missing,error: %s \n", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	tx := common.Db.Model(&model.Order{})
	if req.Status != nil {
		tx.Where("status = ?", *req.Status)
	}
	if req.Direction == "buy" {
		tx.Where("buyer = ?", account)
	} else if req.Direction == "sell" {
		tx.Where("seller = ?", account)
	} else {
		tx.Where("buyer = ? OR seller = ?", account, account)
	}
	var response OrderListResponse
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("db count error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Order("order_time DESC").Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("db find error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type OrderAllReq struct {
	Status *uint8
	PageReq
}

func OrderAll(context *gin.Context) {
	var req OrderAllReq
	err := context.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		logs.Warn(fmt.Sprintf("Body paramter missing,error: %s \n", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	tx := common.Db.Model(&model.Order{})
	if req.Status != nil {
		tx.Where("status = ?", *req.Status)
	}

	var response OrderListResponse
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("db count error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Order("order_time DESC").Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("db find error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type OrderGetReq struct {
	Uuid string `binding:"required"`
}

func OrderGet(context *gin.Context) {
	var req OrderGetReq
	if err := context.ShouldBindUri(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var order model.Order
	tx := common.Db.Where("uuid = ?", req.Uuid).Take(&order)
	if tx.Error != nil {
		resp.Fail(context, "Order not found")
		return
	}

	resp.Success(context, order)
}
