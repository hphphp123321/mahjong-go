package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestOneGame(t *testing.T) {
	// TODO
	var seed int64 = 14

	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)
	var flag = true

	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(players, seed, nil)
	r := rand.New(rand.NewSource(seed))

	//println(i)
	game.Reset(players)
	flag = true
	eventIndex := 0
	for flag {
		posCalls = game.Step()
		for wind, calls := range posCalls {
			posCall[wind] = calls[r.Intn(len(calls))]
		}
		flag = game.Next(posCall)
		posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
		eventIndex++
		boardState := game.GetPosBoardState(mahjong.East)
		bs, _ := json.Marshal(&boardState)
		fmt.Println(string(bs))
		var board mahjong.BoardState
		err := json.Unmarshal(bs, &board)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
	indicators := game.Tiles.DoraIndicators()
	fmt.Println(indicators)
	es := game.GetPosEvents(players[0].Wind, 0)
	b, _ := json.Marshal(&es)

	var events mahjong.Events
	err := json.Unmarshal(b, &events)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(string(b))

}

func TestGame(t *testing.T) {
	// TODO
	var seed int64 = 14

	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)
	var flag = true

	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(players, seed, nil)
	r := rand.New(rand.NewSource(seed))

	for i := 0; i < 100000; i++ {
		//println(i)
		game.Reset(players)
		flag = true
		eventIndex := 0
		for flag {
			posCalls = game.Step()
			for wind, calls := range posCalls {
				posCall[wind] = calls[r.Intn(len(calls))]
			}
			flag = game.Next(posCall)
			posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
			events := game.GetPosEvents(players[0].Wind, eventIndex)
			eventIndex += len(events)
			if len(posCalls) == 4 {
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
