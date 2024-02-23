package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/resp"
	"github.com/gin-gonic/gin"
)

type ClaimableRewardListItem struct {
	Period          uint32
	MachineId       uint64
	PeriodicRewards uint64
}

type RewardTotalResponse struct {
	ClaimedPeriodicRewards   uint64
	ClaimedTaskRewards       uint64
	ClaimablePeriodicRewards uint64
	ClaimableTaskRewards     uint64
}

func RewardTotal(context *gin.Context) {
	var header HttpHeader
	err := context.ShouldBindHeader(&header)
	if err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response RewardTotalResponse
	machine := &model.Machine{Owner: header.Account}
	tx := common.Db.Model(&machine).
		Select("SUM(claimed_periodic_rewards) AS claimed_periodic_rewards, SUM(claimed_task_rewards) AS claimed_task_rewards").
		Where(&machine).
		Find(&response)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	tx = common.Db.Model(&model.RewardMachine{}).
		Select("SUM(rewards.pool DIV rewards.machine_num) AS claimable_periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", header.Account).
		Where("reward_machines.period < ?", currentPeriod()).
		Where("reward_machines.claimed", false).
		Find(&response)
	if tx.Error != nil {
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type RewardClaimableListItem struct {
	Period    uint32
	MachineId string
}

type RewardClaimableListResponse struct {
	List []RewardClaimableListItem
	PageResp
}

func RewardClaimableList(context *gin.Context) {
	var header HttpHeader
	err := context.ShouldBindHeader(&header)
	if err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response RewardClaimableListResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("period,machine_id").
		Where("owner = ?", header.Account).
		Where("period < ?", currentPeriod()).
		Where("claimed", false)
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

type RewardListReq struct {
	Claimed *bool
	PageReq
}

type RewardListItem struct {
	Period          uint32
	Pool            uint64
	PeriodicRewards uint64
}

type RewardListResponse struct {
	List []RewardListItem
	PageResp
}

func RewardPeriodList(context *gin.Context) {
	var header HttpHeader
	err := context.ShouldBindHeader(&header)
	if err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var response RewardListResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("reward_machines.period,rewards.pool,SUM(rewards.pool DIV rewards.machine_num) AS periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", header.Account).
		Where("reward_machines.period < ?", currentPeriod()).
		Where("reward_machines.claimed", false).
		Group("reward_machines.period")
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
