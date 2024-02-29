package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type RewardTotalReq struct {
	Period *uint32
}

type RewardTotalResponse struct {
	ClaimedPeriodicRewards   uint64
	ClaimedTaskRewards       uint64
	ClaimablePeriodicRewards uint64
	ClaimableTaskRewards     uint64
}

func RewardTotal(context *gin.Context) {
	var header HttpHeader
	if err := context.ShouldBindHeader(&header); err != nil {
		logs.Warn(fmt.Sprintf("Header paramter missing: %s \n", err))
		resp.Fail(context, err.Error())
		return
	}
	var req RewardTotalReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		logs.Warn(fmt.Sprintf("Body paramter missing: %s \n", err))
		resp.Fail(context, err.Error())
		return
	}

	var response RewardTotalResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("SUM(rewards.unit_periodic_reward) AS claimed_periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", header.Account).
		Where("reward_machines.claimed", true).
		Where("reward_machines.period < ?", currentPeriod())
	if req.Period != nil {
		tx.Where("reward_machines.period = ?", req.Period)
	}
	tx.Find(&response)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		resp.Fail(context, "Database error")
		return
	}

	tx = common.Db.Model(&model.RewardMachine{}).
		Select("SUM(rewards.unit_periodic_reward) AS claimable_periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", header.Account).
		Where("reward_machines.claimed", false).
		Where("reward_machines.period < ?", currentPeriod())
	if req.Period != nil {
		tx.Where("reward_machines.period = ?", req.Period)
	}
	tx.Find(&response)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type RewardClaimableListReq struct {
	Period *uint32
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
	if err := context.ShouldBindHeader(&header); err != nil {
		logs.Warn(fmt.Sprintf("Header paramter missing: %s \n", err.Error()))
		resp.Fail(context, err.Error())
		return
	}
	var req RewardClaimableListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		logs.Warn(fmt.Sprintf("Body paramter missing: %s \n", err.Error()))
		resp.Fail(context, err.Error())
		return
	}

	var response RewardClaimableListResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("period,machine_id").
		Where("owner = ?", header.Account).
		Where("claimed", false).
		Where("period < ?", currentPeriod())
	if req.Period != nil {
		tx.Where("period = ?", req.Period)
	}
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database count error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database find error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type RewardPeriodListItem struct {
	Period          uint32
	StartTime       time.Time
	Pool            uint64
	PeriodicRewards uint64
}

type RewardPeriodListResponse struct {
	List []RewardPeriodListItem
	PageResp
}

func RewardPeriodList(context *gin.Context) {
	var header HttpHeader
	if err := context.ShouldBindHeader(&header); err != nil {
		logs.Warn(fmt.Sprintf("Header paramter missing: %s \n", err.Error()))
		resp.Fail(context, err.Error())
		return
	}

	var response RewardPeriodListResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("reward_machines.period,rewards.start_time,rewards.pool,SUM(rewards.unit_periodic_reward) AS periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", header.Account).
		Where("reward_machines.period < ?", currentPeriod()).
		Group("reward_machines.period").
		Order("reward_machines.period DESC")
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database count error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database find error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}

type RewardMachineListReq struct {
	Period *uint32 `binding:"required"`
}

type RewardMachineListItem struct {
	Period          uint32
	StartTime       time.Time
	Pool            uint64
	MachineNum      uint32
	PeriodicRewards uint64
	model.Machine
}

type RewardMachineListResponse struct {
	List []RewardMachineListItem
	PageResp
}

func RewardMachineList(context *gin.Context) {
	var header HttpHeader
	if err := context.ShouldBindHeader(&header); err != nil {
		logs.Warn(fmt.Sprintf("Header paramter missing: %s \n", err))
		resp.Fail(context, err.Error())
		return
	}
	var req RewardMachineListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		logs.Warn(fmt.Sprintf("Body paramter missing: %s \n", err))
		resp.Fail(context, err.Error())
		return
	}

	var response RewardMachineListResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("reward_machines.period,rewards.start_time,rewards.pool,rewards.machine_num,rewards.unit_periodic_reward AS periodic_rewards,machines.*").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Joins("LEFT JOIN machines on machines.owner = reward_machines.owner AND machines.uuid = reward_machines.machine_id").
		Where("reward_machines.owner = ?", header.Account).
		Where("reward_machines.period < ?", currentPeriod()).
		Where("reward_machines.period = ?", *req.Period)
	dbResult := tx.Count(&response.Total)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database count error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	dbResult = tx.Scopes(Paginate(context)).Find(&response.List)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Database find error: %s \n", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, response)
}
