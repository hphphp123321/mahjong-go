package mahjong

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"strconv"
)

type Result struct {
	YakuResult  *yaku.Result
	ScoreResult *score.Score
}

func GenerateResult(yakuResult *yaku.Result, scoreResult *score.Score) *Result {
	return &Result{
		YakuResult:  yakuResult,
		ScoreResult: scoreResult,
	}
}

func (r *Result) String() string {
	yakuS := r.YakuResult.String() + "; Total hans: " + strconv.Itoa(int(r.YakuResult.Sum()))
	fusS := r.YakuResult.Fus.String() + "; Total fus: " + strconv.Itoa(int(r.YakuResult.Fus.Sum()))
	return yakuS + "\n" + fusS
}
