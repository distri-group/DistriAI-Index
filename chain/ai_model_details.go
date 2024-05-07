package chain

import (
	"distriai-index-solana/chain/distri_ai"
	"distriai-index-solana/model"
	"time"
)

func buildAiModelModel(m *distri_ai.AiModel) model.AiModel {
	return model.AiModel{
		Owner:      m.Owner.String(),
		Name:       m.Name,
		Framework:  m.Framework,
		License:    m.License,
		Type1:      m.Type1,
		Type2:      m.Type2,
		Tags:       m.Tags,
		CreateTime: time.Unix(m.CreateTime, 0),
		UpdateTime: time.Unix(m.UpdateTime, 0),
	}
}
