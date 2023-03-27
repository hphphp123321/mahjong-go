package mahjong

import (
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
)

var YaoKyuTiles = [...]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

func CalculateShantenNum(handTiles Tiles, melds Calls) int {
	instances, meldsOpt := TilesCallsToCalc(handTiles, melds)
	res := shanten.Calculate(instances, meldsOpt)
	return res.Total.Value
}

func GetTenhaiSlice(handTiles Tiles, melds Calls) []int {
	var tenhaiSlice []int

	instances, meldsOpt := TilesCallsToCalc(handTiles, melds)
	res := tempai.Calculate(instances, meldsOpt)
	tiles := tempai.GetWaits(res).Tiles()
	for _, t := range tiles {
		tenhaiSlice = append(tenhaiSlice, int(t)-1)
	}
	return tenhaiSlice
}

func IndicatorsToDora(indicators Tiles) Tiles {
	var doraTiles Tiles
	for _, indicator := range indicators {
		doraTiles = append(doraTiles, IndicatorToDora(indicator))
	}
	return doraTiles
}

func IndicatorToDora(indicator int) int {
	switch indicator {
	case 32, 33, 34, 35, 68, 69, 70, 71, 104, 105, 106, 107:
		return indicator - 32
	case 120, 121, 122, 123:
		return indicator - 12
	case 132, 133, 134, 135:
		return indicator - 8
	default:
		return indicator + 1
	}
}

func GetYakuResult(handTiles Tiles, melds Calls, ctx *yaku.Context) *yaku.Result {
	instances, meldsOpt := TilesCallsToCalc(handTiles, melds)
	res := tempai.Calculate(instances, meldsOpt)
	if res == nil {
		return nil
	}
	yakuResult := yaku.Win(res, ctx, nil)
	return yakuResult
}

func GetScoreResult(scoreRule *score.RulesStruct, yakuResult *yaku.Result, honba int) score.Score {
	return score.GetScoreByResult(scoreRule, yakuResult, score.Honba(honba))
}
