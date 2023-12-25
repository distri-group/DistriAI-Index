package handlers

import (
	"distriai-backend-solana/common"
	"distriai-backend-solana/model"
	"distriai-backend-solana/utils/resp"
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
	var header HttpHeader
	err := context.ShouldBindHeader(&header)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}
	var req OrderMineReq
	err = context.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}

	tx := common.Db.Model(&model.Order{})
	if req.Status != nil {
		tx.Where("status = ?", *req.Status)
	}
	if req.Direction == "buy" {
		tx.Where("buyer = ?", header.Account)
	} else if req.Direction == "sell" {
		tx.Where("seller = ?", header.Account)
	} else {
		tx.Where("buyer = ? OR seller = ?", header.Account, header.Account)
	}
	var response OrderListResponse
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		fmt.Println(dbResult.Error)
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Order("order_time DESC").Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
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
		fmt.Println(dbResult.Error)
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Order("order_time DESC").Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}
