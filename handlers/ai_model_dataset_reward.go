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

// AiModelDatasetRewardPoolTotal is a handler function that calculates the total reward pool for AI model dataset reward periods.
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

// AiModelDatasetRewardPeriodDetail is a handler function that retrieves details of a specific reward period for AI model datasets.
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
type AiModelDatasetRewardItem struct {
	model.AiModelDatasetReward
	model.AiModelDatasetRewardPeriod
}

type AiModelDatasetRewardListResponse struct {
	List []AiModelDatasetRewardItem
	PageResp
}

// AiModelDatasetRewardList handles the request to list rewards for AI model datasets.
func AiModelDatasetRewardList(context *gin.Context) {
func AiModelDatasetRewardList(context *gin.Context) {
	var req AiModelDatasetRewardListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response AiModelDatasetRewardListResponse
	tx := common.Db.Model(&model.AiModelDatasetReward{}).
		Select("ai_model_dataset_rewards.*, ai_model_dataset_reward_periods.*").
		Joins("LEFT JOIN ai_model_dataset_reward_periods ON ai_model_dataset_rewards.period = ai_model_dataset_reward_periods.period")
	if req.Owner != nil {
		tx.Where("ai_model_dataset_rewards.owner = ?", req.Owner).
			Order("ai_model_dataset_rewards.period DESC")
	}
	if req.Period != nil {
		tx.Where("ai_model_dataset_rewards.period = ?", req.Period).
			Order("ai_model_dataset_rewards.periodic_reward DESC")
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
	Owner  string `binding:"required"`
	Period uint32
}

// AiModelDatasetRewardDetail handles the request to get the detail of a specific AI model dataset reward.
func AiModelDatasetRewardDetail(context *gin.Context) {
	var req AiModelDatasetRewardDetailReq
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var reward AiModelDatasetRewardItem
	tx := common.Db.Model(&model.AiModelDatasetReward{}).
		Select("ai_model_dataset_rewards.*, ai_model_dataset_reward_periods.*").
		Joins("LEFT JOIN ai_model_dataset_reward_periods ON ai_model_dataset_rewards.period = ai_model_dataset_reward_periods.period").
		Where("ai_model_dataset_rewards.period = ? AND ai_model_dataset_rewards.owner = ?", req.Period, req.Owner).
		Take(&reward)
	if tx.Error != nil {
		resp.Fail(context, "Not found")
		return
	}

	resp.Success(context, reward)
}
