package tests

import (
	"fmt"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"strconv"
	"testing"
)

func TestShanten(t *testing.T) {
	generator := compact.NewTileGenerator()
	a := compact.NewInstances()
	handTiles := tile.Tiles{5, 6, 14, 15, 17, 17, 33, 33, 4, 16}
	hand := generator.Tiles(handTiles)
	a.Add(hand)

	declared := []calc.Meld{calc.Kan(tile.Tile(21))}
	melds := calc.Melds{}
	melds = append(melds, declared...)
	cal := calc.Declared(melds)

	res := shanten.Calculate(a, cal)

	results := tempai.Calculate(a, cal)

	fmt.Printf("Total shanten value is: %v\n", res.Total.Value)
	fmt.Printf("Waits are %s\n", tempai.GetWaits(results).Tiles())

	var tenhaiSlice []int
	tiles := tempai.GetWaits(results).Tiles()
	for _, tileID := range tiles {
		tenhaiSlice = append(tenhaiSlice, int(tileID))
	}
}

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

func TestCalculate(t *testing.T) {
	p1 := mahjong.NewMahjongPlayer()
	p1.ResetForGame()
	p1.HandTiles = mahjong.Tiles{14, 19, 21, 55, 57, 63, 64, 65, 88, 129, 131}
	p1.Melds = mahjong.Calls{&mahjong.Call{
		CallType:         mahjong.AnKan,
		CallTiles:        mahjong.Tiles{80, 81, 82, 83},
		CallTilesFromWho: []mahjong.Wind{3, 3, 3, 3},
	}}
	p1.GetRiichiTiles()

	//handTiles = tile.Tiles{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5}
	//hand = generator.Tiles(handTiles)
	//a = compact.NewInstances()
	//a.Add(hand)

}
