package model

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {
	db.Migrator().DropTable(
		&AiModel{},
		&Dataset{},
		&Machine{},
		&Order{},
		&Reward{},
		&RewardMachine{},
	)
	db.AutoMigrate(
		&AiModel{},
		&AiModelHeat{},
		&AiModelLike{},
		&AiModelDatasetRewardPeriod{},
		&AiModelDatasetReward{},
		&Dataset{},
		&DatasetHeat{},
		&DatasetLike{},
		&Log{},
		&Mailbox{},
		&Machine{},
		&MachineInfo{},
		&Order{},
		&Reward{},
		&RewardMachine{},
	)
}
