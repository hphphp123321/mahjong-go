package mahjong

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/tile"
	common "github.com/hphphp123321/go-common"
)

func IntToInstance(t int) tile.Instance {
	return tile.Instance(t + 1)
}

func IntsToInstances(tiles Tiles) tile.Instances {
	instances := tile.Instances{}
	for _, num := range tiles {
		instances = append(instances, IntToInstance(int(num)))
	}
	return instances
}

func IntsToTiles(tiles Tiles) tile.Tiles {
	tilesT := tile.Tiles{}
	for _, num := range tiles {
		tilesT = append(tilesT, tile.Tile(num/4+1))
	}
	return tilesT
}

func CallToMeld(call *Call) calc.Meld {
	var meld calc.Meld
	switch call.CallType {
	case Chi:
		var tileClass = common.MinNum(call.CallTiles[:3]).Class() + 1
		meld = calc.Open(calc.Chi(tile.Tile(tileClass)))
	case Pon:
		var tileClass = call.CallTiles[0].Class() + 1
		meld = calc.Open(calc.Pon(tile.Tile(tileClass)))
	case DaiMinKan:
		var tileClass = call.CallTiles[0].Class() + 1
		meld = calc.Open(calc.Kan(tile.Tile(tileClass)))
	case ShouMinKan:
		var tileClass = call.CallTiles[0].Class() + 1
		meld = calc.Open(calc.Kan(tile.Tile(tileClass)))
	case AnKan:
		var tileClass = call.CallTiles[0].Class() + 1
		meld = calc.Kan(tile.Tile(tileClass))
	}
	return meld
}

func CallsToMelds(melds Calls) calc.Melds {
	meldsT := calc.Melds{}
	for _, v := range melds {
		meldsT = append(meldsT, CallToMeld(v))
	}
	return meldsT
}

func TilesCallsToCalc(tiles Tiles, calls Calls) (compact.Instances, calc.Option) {
	hand := IntsToInstances(tiles)
	instances := compact.NewInstances()
	instances.Add(hand)

	var meldsOpt calc.Option = nil
	if calls != nil {
		meldsT := CallsToMelds(calls)
		meldsOpt = calc.Declared(meldsT)
	}
	return instances, meldsOpt
}
