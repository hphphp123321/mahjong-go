package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"testing"
)

func TestJson(t *testing.T) {
	var e mahjong.Events
	e1 := mahjong.EventAnKan{
		Who: mahjong.South,
		Call: &mahjong.Call{
			CallType:         mahjong.AnKan,
			CallTiles:        mahjong.Tiles{1, 2, 3, 4},
			CallTilesFromWho: []mahjong.Wind{mahjong.South, mahjong.South, mahjong.South, mahjong.South},
		},
	}
	e2 := mahjong.EventRiichi{
		Who: mahjong.South,
	}
	e3 := mahjong.EventGet{
		Who:  mahjong.North,
		Tile: 10,
	}
	e4 := mahjong.EventStart{
		WindRound: mahjong.WindRoundEast2,
		InitWind:  mahjong.South,
		Seed:      0,
		NumGame:   0,
		NumHonba:  0,
		NumRiichi: 0,
		InitTiles: mahjong.Tiles{1, 3, 5, 6, 8, 9},
	}
	e5 := mahjong.EventRyuuKyoku{
		Who: mahjong.South,
		//HandTiles: mahjong.Tiles{12, 58},
		Reason: mahjong.RyuuKyokuKyuuShuKyuuHai,
	}
	e6 := mahjong.EventEnd{
		PointsChange: map[mahjong.Wind]int{
			mahjong.South: 1000,
			mahjong.North: -1000,
			mahjong.East:  0,
			mahjong.West:  0,
		},
	}
	e = append(e, &e1)
	e = append(e, &e2)
	e = append(e, &e3)
	e = append(e, &e4)
	e = append(e, &e5)
	e = append(e, &e6)
	b, err := json.Marshal(&e)
	fmt.Println("json: " + string(b))
	var a mahjong.Events
	err = json.Unmarshal(b, &a)
	if err != nil {
		panic(err)
	}
}
