package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestOneGame(t *testing.T) {
	var seed int64 = 14

	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)

	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(seed, nil)
	posCalls = game.Reset(players, nil)
	var flag = mahjong.EndTypeNone
	for flag != mahjong.EndTypeGame {
		for wind, calls := range posCalls {
			posCall[wind] = calls[rand.Intn(len(calls))]
		}
		posCalls, flag = game.Step(posCall)
		posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
	}
}

func TestMultiGames(t *testing.T) {
	var seed int64 = 14

	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)
	var flag = mahjong.EndTypeNone

	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(seed, nil)
	r := rand.New(rand.NewSource(seed))

	var maxNum = 0

	for i := 0; i < 100; i++ {
		//println(i)
		posCalls = game.Reset(players, nil)
		flag = mahjong.EndTypeNone
		eventIndex := 0
		for flag != mahjong.EndTypeGame {
			for wind, calls := range posCalls {
				posCall[wind] = calls[r.Intn(len(calls))]
			}
			posCalls, flag = game.Step(posCall)
			posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
			events := game.GetPosEvents(players[0].Wind, eventIndex)
			eventIndex += len(events)
			if len(posCalls) == 4 {
				for _, player := range players {
					if len(player.DiscardTiles) > maxNum {
						maxNum = len(player.DiscardTiles)
						fmt.Println(maxNum)
					}
				}
				eventIndex = 0
				if events[0].GetType() == mahjong.EventTypeTsumo ||
					events[0].GetType() == mahjong.EventTypeRon ||
					events[0].GetType() == mahjong.EventTypeChanKan {
					b, _ := json.Marshal(&events)
					fmt.Println(string(b))
				}
			}
		}
	}
}
