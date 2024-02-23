package handlers

import "time"

const (
	genesisTime    int64 = 1706745600
	periodDuration int64 = 86400
)

func currentPeriod() uint32 {
	return uint32((time.Now().Unix() - genesisTime) / periodDuration)
}
