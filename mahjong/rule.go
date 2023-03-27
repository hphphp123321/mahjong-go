package mahjong

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
)

type Rule struct {
	yakuRule  *yaku.RulesStruct
	scoreRule *score.RulesStruct

	limitTime int // player play time limitation, 0 means no limit
}

func GetDefaultRule() *Rule {
	return &Rule{
		yakuRule:  yaku.RulesTenhouRed(),
		scoreRule: score.RulesTenhou(),
		limitTime: 0,
	}
}
