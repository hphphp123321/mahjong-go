package mahjong

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
)

type Result struct {
	YakuResult      *YakuResult  `json:"yakuResult,omitempty"`
	ScoreResult     *ScoreResult `json:"scoreResult,omitempty"`
	RonCall         *Call
	RyuuKyokuReason RyuuKyokuReason
}

func GenerateRonResult(yakuResult *yaku.Result, scoreResult *score.Score) (r *Result) {
	r = &Result{}
	yakuR := make(map[Yaku]int)
	yakuMans := make([]Yakuman, 0)
	bonuses := make(map[Yaku]int)
	fus := make(Fus, 0)
	for k, v := range yakuResult.Yaku {
		yakuR[Yaku(k)] = int(v)
	}
	for _, v := range yakuResult.Yakumans {
		yakuMans = append(yakuMans, Yakuman(v))
	}
	for k, v := range yakuResult.Bonuses {
		bonuses[Yaku(k)] = int(v)
	}
	for _, v := range yakuResult.Fus {
		fus = append(fus, &FuInfo{
			Fu:     Fu(v.Fu),
			Points: int(v.Points),
		})
	}
	r.YakuResult = &YakuResult{
		Yaku:     yakuR,
		Yakumans: yakuMans,
		Bonuses:  bonuses,
		Fus:      fus,
		IsClosed: yakuResult.IsClosed,
	}

	r.ScoreResult = &ScoreResult{
		PayRon:         int(scoreResult.PayRon),
		PayRonDealer:   int(scoreResult.PayRonDealer),
		PayTsumo:       int(scoreResult.PayTsumo),
		PayTsumoDealer: int(scoreResult.PayTsumoDealer),
		Special:        Limit(scoreResult.Special),
		Han:            int(scoreResult.Han),
		Fu:             int(scoreResult.Fu),
	}
	return r
}

//func (r *Result) String() string {
//	if r.RyuuKyokuReason != NoRyuuKyoku {
//		return r.RyuuKyokuReason.String()
//	}
//	yakuS := r.YakuResult.String() + "; Total hans: " + strconv.Itoa(int(r.YakuResult.Sum()))
//	fusS := r.YakuResult.Fus.String() + "; Total fus: " + strconv.Itoa(int(r.YakuResult.Fus.Sum()))
//	return yakuS + "\n" + fusS
//}

type ScoreResult struct {
	PayRon         int   `json:"payRon,omitempty"`
	PayRonDealer   int   `json:"payRonDealer,omitempty"`
	PayTsumo       int   `json:"payTsumo,omitempty"`
	PayTsumoDealer int   `json:"payTsumoDealer,omitempty"`
	Special        Limit `json:"yaku_limit,omitempty"`
	Han            int   `json:"han,omitempty"`
	Fu             int   `json:"fu,omitempty"`
}

type YakuResult struct {
	Yaku     YakuSet  `json:"yaku,omitempty"`
	Yakumans Yakumans `json:"yakumans,omitempty"`
	Bonuses  YakuSet  `json:"bonuses,omitempty"`
	Fus      Fus      `json:"fus,omitempty"`
	IsClosed bool     `json:"isClosed,omitempty"`
}

type ScoreChanges map[Wind]int

func (sc ScoreChanges) TotalWin() int {
	for _, v := range sc {
		if v > 0 {
			return v
		}
	}
	return 0
}

func (sc ScoreChanges) TotalPayed() (total int) {
	for _, v := range sc {
		if v < 0 {
			total -= v
		}
	}
	return
}

func (s *ScoreResult) GetChanges(selfWind, otherWind Wind, sticks int) ScoreChanges {
	scoreChanges := make(ScoreChanges)
	scoreS := &score.Score{
		PayRon:         score.Money(s.PayRon),
		PayRonDealer:   score.Money(s.PayRonDealer),
		PayTsumo:       score.Money(s.PayTsumo),
		PayTsumoDealer: score.Money(s.PayTsumoDealer),
		Special:        yaku.Limit(s.Special),
		Han:            yaku.HanPoints(s.Han),
		Fu:             yaku.FuPoints(s.Fu),
	}
	scoreC := scoreS.GetChanges(base.Wind(selfWind), base.Wind(otherWind), score.RiichiSticks(sticks))
	for i, v := range scoreC {
		scoreChanges[Wind(i)] = int(v)
	}
	return scoreChanges
}
