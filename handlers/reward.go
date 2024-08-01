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

type RewardPeriodDetailReq struct {
	Period uint32
}

// RewardPeriodDetail is a handler function that retrieves details of a specific reward period.
func RewardPeriodDetail(context *gin.Context) {
	var req RewardPeriodDetailReq
	// Attempt to bind the JSON body of the request to the req variable.
	// If there is an error during binding, respond with a failure message.
	if err := context.ShouldBindJSON(&req); err != nil {
		resp.Fail(context, err.Error())
		return
	}

	var reward model.Reward
	tx := common.Db.Model(&model.Reward{}).Where("period = ?", req.Period).Take(&reward)
	if tx.Error != nil {
		resp.Fail(context, "Not found")
		return
	}

	resp.Success(context, reward)
}

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
	account, err := getAccount(context)
	if err != nil {
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
		Select("IFNULL(SUM(rewards.unit_periodic_reward), 0) AS claimed_periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", account).
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
		Select("IFNULL(SUM(rewards.unit_periodic_reward), 0) AS claimable_periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", account).
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
	Period          uint32
	MachineId       string
	PeriodicRewards uint64
}

type RewardClaimableListResponse struct {
	List []RewardClaimableListItem
	PageResp
}

// RewardClaimableList handles the request to list claimable rewards for an account
func RewardClaimableList(context *gin.Context) {
	account, err := getAccount(context)
	if err != nil {
		return
	}
	var req RewardClaimableListReq
	if err := context.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		logs.Warn(fmt.Sprintf("Body paramter missing: %s \n", err.Error()))
		resp.Fail(context, err.Error())
		return
	}

	var response RewardClaimableListResponse
	// Start a new database query using the common.Db model for RewardMachine
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("reward_machines.period,reward_machines.machine_id,rewards.unit_periodic_reward AS periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", account).
		Where("reward_machines.claimed", false).
		Where("reward_machines.period < ?", currentPeriod())
	if req.Period != nil {
		tx.Where("reward_machines.period = ?", req.Period)
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
	account, err := getAccount(context)
	if err != nil {
		return
	}

	var response RewardPeriodListResponse
	tx := common.Db.Model(&model.RewardMachine{}).
		Select("reward_machines.period,rewards.start_time,rewards.pool,IFNULL(SUM(rewards.unit_periodic_reward), 0) AS periodic_rewards").
		Joins("LEFT JOIN rewards on rewards.period = reward_machines.period").
		Where("reward_machines.owner = ?", account).
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
	account, err := getAccount(context)
	if err != nil {
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
		Where("reward_machines.owner = ?", account).
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
