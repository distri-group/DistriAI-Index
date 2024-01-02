package handlers

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MachineFilterResponse struct {
	Gpu      []string
	GpuCount []uint32
	Region   []string
}

func MachineFilter(context *gin.Context) {
	var response MachineFilterResponse
	dbResult := common.Db.Model(&model.Machine{}).Select("gpu").Group("gpu").Find(&response.Gpu)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	dbResult = common.Db.Model(&model.Machine{}).Select("gpu_count").Group("gpu_count").Find(&response.GpuCount)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	dbResult = common.Db.Model(&model.Machine{}).Select("region").Group("region").Find(&response.Region)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type MachineListReq struct {
	Gpu      string
	GpuCount uint32
	Region   string
	Status   *uint8
	OrderBy  string
	PageReq
}

type MachineListResponse struct {
	List []model.Machine
	PageResp
}

func MachineMarket(context *gin.Context) {
	var req MachineListReq
	err := context.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}

	machine := &model.Machine{Gpu: req.Gpu, GpuCount: req.GpuCount, Region: req.Region}
	var response MachineListResponse
	tx := common.Db.Model(&machine).Where(&machine)
	if req.Status == nil {
		tx.Where("status != ?", uint8(distri_ai.MachineStatusIdle))
	} else {
		tx.Where("status = ?", *req.Status)
	}
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	switch req.OrderBy {
	case "price", "price DESC", "score DESC", "tflops DESC":
		tx.Order(req.OrderBy)
	case "reliability":
		tx.Order("`completed_count`/(`completed_count` + `failed_count`) DESC")
	}
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

func MachineMine(context *gin.Context) {
	var header HttpHeader
	err := context.ShouldBindHeader(&header)
	if err != nil {
		resp.Fail(context, "Parameter missing")
		return
	}

	machine := &model.Machine{Owner: header.Account}
	var response MachineListResponse
	tx := common.Db.Model(&machine).Where(&machine)
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}
