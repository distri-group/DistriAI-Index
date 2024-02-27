package handlers

import "time"

const (
	// Period 0 start time: 2024-02-27 00:00:00 UTC
	genesisTime    int64 = 1708992000
	periodDuration int64 = 86400
)

func currentPeriod() uint32 {
	return uint32((time.Now().Unix() - genesisTime) / periodDuration)
}
