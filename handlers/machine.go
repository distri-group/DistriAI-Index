package handlers

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
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
		logs.Error(fmt.Sprintf("Database error based on gpu query: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = common.Db.Model(&model.Machine{}).Select("gpu_count").Group("gpu_count").Find(&response.GpuCount)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error based on gpu_count query: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = common.Db.Model(&model.Machine{}).Select("region").Group("region").Find(&response.Region)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error based on region query: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type MachineListReq struct {
	Gpu        string
	GpuCount   uint32
	Region     string
	Status     *uint8
	OrderBy    string
	PriceOrder int
	PageReq
}

type MachineListResponse struct {
	List []model.Machine
	PageResp
}

func MachineMarket(context *gin.Context) {
	// bind request params
	var req MachineListReq
	err := context.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		logs.Warn(fmt.Sprintf("RequestBody Parameter missing,error: %s \n", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	// build sql query params
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

	if req.PriceOrder == 1 {
		tx.Order("price DESC")
	} else if req.PriceOrder == 2 {
		tx.Order("price ASC")
	}
	switch req.OrderBy {
	case "price", "price DESC", "score DESC", "tflops DESC", "max_duration DESC", "disk DESC", "ram DESC":
		tx.Order(req.OrderBy)
	case "reliability":
		tx.Order("`completed_count`/(`completed_count` + `failed_count`) DESC")
	case "":
		tx.Order("status ASC,tflops DESC")
	}

	// execute pagination query
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

func MachineMine(context *gin.Context) {
	account, err := getAccount(context)
	if err != nil {
		return
	}

	machine := &model.Machine{Owner: account}
	var response MachineListResponse
	tx := common.Db.Model(&machine).Where(&machine)
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", err))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", err))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type MachineGetReq struct {
	Owner string `binding:"required"`
	Uuid  string `binding:"required"`
}

func MachineGet(context *gin.Context) {
	var req MachineGetReq
	if err := context.ShouldBindUri(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var machine model.Machine
	tx := common.Db.
		Where("owner = ?", req.Owner).
		Where("uuid = ?", req.Uuid).
		Take(&machine)
	if tx.Error != nil {
		resp.Fail(context, "Machine not found")
		return
	}

	resp.Success(context, machine)
}
