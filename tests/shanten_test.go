package tests

import (
	"fmt"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"testing"
)

// 1, 4, 5, 11, 17, 22, 28, 30, 32, 34  	27
func TestShanten(t *testing.T) {
	generator := compact.NewTileGenerator()
	a := compact.NewInstances()
	handTiles := tile.Tiles{11, 12, 13, 28, 28, 28, 24, 24, 24, 27}
	hand := generator.Tiles(handTiles)
	a.Add(hand)

	declared := []calc.Meld{calc.Open(calc.Chi(tile.Tile(2)))}
	melds := calc.Melds{}
	melds = append(melds, declared...)
	cal := calc.Declared(melds)

	res := shanten.Calculate(a, cal)

	//results := tempai.Calculate(a, cal)

	fmt.Printf("Total shanten value is: %v\n", res.Total.Value)

	//fmt.Printf("Waits are %s\n", tempai.GetWaits(results).Tiles())

	//var tenhaiSlice []int
	//tiles := tempai.GetWaits(results).Tiles()
	//for _, tileID := range tiles {
	//	tenhaiSlice = append(tenhaiSlice, int(tileID))
	//}
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
