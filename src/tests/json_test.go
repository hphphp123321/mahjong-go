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
	e = append(e, &e1)
	e = append(e, &e2)
	e = append(e, &e3)
	b, err := json.Marshal(&e)
	fmt.Println(string(b))
	var a mahjong.Events
	err = json.Unmarshal(b, &a)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
