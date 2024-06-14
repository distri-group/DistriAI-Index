package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"time"
)

const (
	// Period 0 start time: 2024-02-27 00:00:00 UTC
	genesisTime          int64 = 1708992000
	periodDuration       int64 = 86400
	decayPeriods               = 4
	decayRateNumerator         = 9737
	decayRateDenominator       = 10000
)

var poolCheckpoints = []uint64{
	65_750_000_000_000,
	50_367_000_000_000,
	38_583_000_000_000,
	29_556_000_000_000,
	22_641_000_000_000,
	17_344_000_000_000,
	13_286_000_000_000,
	10_178_000_000_000,
	7_797_000_000_000,
	5_973_000_000_000,
	4_575_000_000_000,
	3_505_000_000_000,
	2_685_000_000_000,
	2_057_000_000_000,
	1_576_000_000_000,
	1_207_000_000_000,
	925_000_000_000,
	708_000_000_000,
	543_000_000_000,
	416_000_000_000,
	318_000_000_000,
	244_000_000_000,
	187_000_000_000,
	143_000_000_000,
	110_000_000_000,
	84_000_000_000,
	64_000_000_000,
	49_000_000_000,
	38_000_000_000,
	29_000_000_000,
	22_000_000_000,
}

func currentPeriod() uint32 {
	return uint32((time.Now().Unix() - genesisTime) / periodDuration)
}

func startTime(period uint32) time.Time {
	return time.Unix(genesisTime+periodDuration*int64(period), 0)
}

func pool(period uint32) (total uint64, aiModelDataset uint64) {
	decayTimes := int(period / decayPeriods)
	checkpointIndex := decayTimes / 10
	if checkpointIndex > len(poolCheckpoints)-1 {
		checkpointIndex = len(poolCheckpoints) - 1
	}
	remainingDecayTimes := decayTimes - checkpointIndex*10
	total = poolCheckpoints[checkpointIndex]
	for i := 0; i < remainingDecayTimes; i++ {
		total = total * decayRateNumerator / decayRateDenominator
	}
	aiModelDataset = total * 2 / 10
	return total, aiModelDataset
}

func StartRewardCron() {
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		logs.Error(fmt.Sprintf("Error scheduling job: %s \n", err))
		return
	}

	_, err = scheduler.NewJob(
		gocron.CronJob("58 23 * * *", false),
		gocron.NewTask(func() {
			createAiModelDatasetRewardPeriod()
		}),
	)

	if err != nil {
		logs.Error(fmt.Sprintf("Error scheduling job: %s \n", err))
		return
	}
	scheduler.Start()
}

func createAiModelDatasetRewardPeriod() {
	var aiModelCount int64
	tx := common.Db.Table("ai_models").
		Joins("LEFT JOIN ai_model_heats ON ai_models.owner = ai_model_heats.owner AND ai_models.name = ai_model_heats.name").
		Where("ai_model_heats.review = 1").
		Count(&aiModelCount)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		return
	}
	var datasetCount int64
	tx = common.Db.Table("datasets").
		Joins("LEFT JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name").
		Where("dataset_heats.review = 1").
		Count(&datasetCount)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		return
	}

	period := currentPeriod()
	logs.Info(fmt.Sprintf("Create AiModel & Dataset Reward Period: %d \n", period))
	_, rewardPool := pool(period)
	rewardPeriod := model.AiModelDatasetRewardPeriod{
		Period:     period,
		StartTime:  startTime(period),
		Pool:       rewardPool,
		AiModelNum: uint32(aiModelCount),
		DatasetNum: uint32(datasetCount),
	}
	totalNum := aiModelCount + datasetCount
	if totalNum > 0 {
		rewardPeriod.UnitPeriodicReward = rewardPool / uint64(totalNum)
	}

	tx = common.Db.Create(&rewardPeriod)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		return
	}

	createAiModelDatasetReward(rewardPeriod)
}

func createAiModelDatasetReward(rewardPeriod model.AiModelDatasetRewardPeriod) {
	var aiModelRewards []model.AiModelDatasetReward
	tx := common.Db.Table("ai_models").
		Select("ai_models.owner AS owner, COUNT(ai_models.owner) AS ai_model_num").
		Joins("LEFT JOIN ai_model_heats ON ai_models.owner = ai_model_heats.owner AND ai_models.name = ai_model_heats.name").
		Where("ai_model_heats.review = 1").
		Group("ai_models.owner").
		Find(&aiModelRewards)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		return
	}

	var datasetRewards []model.AiModelDatasetReward
	tx = common.Db.Table("datasets").
		Select("datasets.owner AS owner, COUNT(datasets.owner) AS dataset_num").
		Joins("LEFT JOIN dataset_heats ON datasets.owner = dataset_heats.owner AND datasets.name = dataset_heats.name").
		Where("dataset_heats.review = 1").
		Group("datasets.owner").
		Find(&datasetRewards)
	if tx.Error != nil {
		logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
		return
	}

	m := make(map[string]model.AiModelDatasetReward)
	for _, reward := range aiModelRewards {
		m[reward.Owner] = reward
	}
	for _, reward := range datasetRewards {
		value, ok := m[reward.Owner]
		if ok {
			value.DatasetNum = reward.DatasetNum
			m[reward.Owner] = value
		} else {
			m[reward.Owner] = reward
		}
	}

	var aiModelDatasetRewards []model.AiModelDatasetReward
	for _, reward := range m {
		reward.Period = rewardPeriod.Period
		reward.PeriodicReward = uint64(reward.AiModelNum+reward.DatasetNum) * rewardPeriod.UnitPeriodicReward
		aiModelDatasetRewards = append(aiModelDatasetRewards, reward)
	}
	if len(aiModelDatasetRewards) > 0 {
		tx = common.Db.Create(&aiModelDatasetRewards)
		if tx.Error != nil {
			logs.Error(fmt.Sprintf("Database error: %s \n", tx.Error))
			return
		}
	}
}
