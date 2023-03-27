package tests

import (
	"fmt"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"strconv"
	"testing"
)

func TestTenhai(t *testing.T) {
	generator := compact.NewTileGenerator()
	a := compact.NewInstances()
	handTiles := tile.Tiles{5, 6, 14, 15, 17, 17, 33, 33, 4, 16}
	hand := generator.Tiles(handTiles)
	a.Add(hand)

	declared := []calc.Meld{calc.Kan(tile.Tile(21))}
	melds := calc.Melds{}
	melds = append(melds, declared...)
	cal := calc.Declared(melds)
	results := tempai.Calculate(a, cal)
	winTile := generator.Instance(tile.Tile(33))
	ctx := &yaku.Context{
		Tile:        winTile,
		SelfWind:    0,
		RoundWind:   0,
		DoraTiles:   tile.Tiles{33},
		UraTiles:    tile.Tiles{17},
		Rules:       yaku.RulesTenhouRed(),
		IsTsumo:     true,
		IsRiichi:    false,
		IsIpatsu:    false,
		IsDaburi:    false,
		IsLastTile:  false,
		IsRinshan:   false,
		IsFirstTake: false,
		IsChankan:   false,
	}

	yakuResult := yaku.Win(results, ctx, nil)

	scoreResult := score.GetScoreByResult(score.RulesTenhou(), yakuResult, 0)
	fmt.Println(scoreResult.PayRon, scoreResult.PayTsumo, scoreResult.PayRonDealer, scoreResult.PayTsumoDealer, scoreResult.Special)

	fmt.Println(yakuResult.Yakumans.String())
	fmt.Println(yakuResult.Yaku.String())
	fmt.Println(yakuResult.Bonuses.String())
	fmt.Println(yakuResult.Fus.String())
	fmt.Println(yakuResult.Sum())
	fmt.Println(strconv.FormatBool(yakuResult.IsClosed))
	fmt.Printf("%v\n", yakuResult.String())
}
