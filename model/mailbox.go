package model

import "time"

type Mailbox struct {
	Id        int    `gorm:"primarykey"`
	MailBox   string `gorm:"size:64;not null;unique;uniqueIndex"`
	CreatedAt time.Time
}
