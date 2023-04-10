package mahjong

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
)

type Rule struct {
	YakuRule  *yaku.RulesStruct
	ScoreRule *score.RulesStruct

	GameLength    int  // game length 1 for 1 round only, 4 for tonpuusen, 8 for hanchan, etc.
	SanChaRon     bool // can san chan ron, true for can, false for can't -> ryuu kyoku
	NagashiMangan bool // can nagashi/ryuukyoku mangan, true for can, false for can't

	//limitTime int // player play time limitation, 0 means no limit
}

func GetDefaultRule() *Rule {
	return &Rule{
		YakuRule:  yaku.RulesTenhouRed(),
		ScoreRule: score.RulesTenhou(),

		GameLength:    8,
		SanChaRon:     false,
		NagashiMangan: true,
		//limitTime: 0,
	}
}
