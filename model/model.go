package model

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {
	db.Migrator().DropTable(
		&Sync{},
		&Machine{},
		&Order{},
		&Reward{},
		&RewardMachine{},
	)
	db.AutoMigrate(
		&Log{},
		&Mailbox{},
		&Sync{},
		&Machine{},
		&Order{},
		&Reward{},
		&RewardMachine{},
	)
}
