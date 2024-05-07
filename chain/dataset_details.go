package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"time"
)

func buildDatasetModel(m *distri_ai.Dataset) model.Dataset {
	return model.Dataset{
		Owner:      m.Owner.String(),
		Name:       m.Name,
		Scale:      m.Scale,
		License:    m.License,
		Type1:      m.Type1,
		Type2:      m.Type2,
		Tags:       m.Tags,
		CreateTime: time.Unix(m.CreateTime, 0),
		UpdateTime: time.Unix(m.UpdateTime, 0),
	}
}
