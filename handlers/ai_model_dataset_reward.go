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

func AiModelDatasetRewardPoolTotal(context *gin.Context) {
	var total uint64
	tx := common.Db.Model(&model.AiModelDatasetRewardPeriod{}).
		Select("IFNULL(SUM(pool), 0) AS total").
		Find(&total)
	if tx.Error != nil {
		resp.Fail(context, "Not found")
		return
	}

	resp.Success(context, total)
}

type AiModelDatasetRewardPeriodDetailReq struct {
	Period *uint32
}

func AiModelDatasetRewardPeriodDetail(context *gin.Context) {
	var req AiModelDatasetRewardPeriodDetailReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var reward model.AiModelDatasetRewardPeriod
	tx := common.Db.Model(&model.AiModelDatasetRewardPeriod{})
	if req.Period == nil {
		tx.Order("period DESC")
	} else {
		tx.Where("period = ?", req.Period)
	}
	if tx.Take(&reward).Error != nil {
		resp.Fail(context, "Not found")
		return
	}

	resp.Success(context, reward)
}

type AiModelDatasetRewardListReq struct {
	Owner  *string
	Period *uint32
}

type AiModelDatasetRewardListResponse struct {
	List []model.AiModelDatasetReward
	PageResp
}

func AiModelDatasetRewardList(context *gin.Context) {
	var req AiModelDatasetRewardListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response AiModelDatasetRewardListResponse
	tx := common.Db.Model(&model.AiModelDatasetReward{})
	if req.Owner != nil {
		tx.Where("owner = ?", req.Owner).
			Order("period DESC")
	}
	if req.Period != nil {
		tx.Where("period = ?", req.Period).
			Order("period_reward DESC")
	}

	tx = tx.Count(&response.Total)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database count error: %s \n", tx.Error))
		resp.Fail(context, "Database error")
		return
	}
	tx = tx.Scopes(Paginate(context)).Find(&response.List)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database find error: %s \n", tx.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type AiModelDatasetRewardDetailReq struct {
	Period uint32 `binding:"required"`
}

func AiModelDatasetRewardDetail(context *gin.Context) {
	account, err := getAccount(context)
	if err != nil {
		return
	}
	var req AiModelDatasetRewardDetailReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var reward model.AiModelDatasetReward
	tx := common.Db.Model(&model.AiModelDatasetReward{}).
		Where("period = ? AND owner = ?", req.Period, account).
		Take(&reward)
	if tx.Error != nil {
		resp.Fail(context, "Not found")
		return
	}

	resp.Success(context, reward)
}
