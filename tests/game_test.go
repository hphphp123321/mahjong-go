package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestGame(t *testing.T) {
	// TODO
	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)
	var flag = true

	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(players, 0, nil)

	eventIndex := 0
	for flag {
		posCalls = game.Step()
		for wind, calls := range posCalls {
			posCall[wind] = calls[rand.Intn(len(calls))]
		}
		flag = game.Next(posCall)
		posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
		events := game.GetPosEvents(players[0].Wind, eventIndex)
		b, _ := json.Marshal(&events)
		fmt.Println(string(b))
		eventIndex += len(events)
		if len(posCalls) == 4 {
			eventIndex = 0
		}
	}
}
