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
		&Dataset{},
		&DatasetHeat{},
		&DatasetLike{},
		&Log{},
		&Mailbox{},
		&Machine{},
		&Order{},
		&Reward{},
		&RewardMachine{},
	)
}
