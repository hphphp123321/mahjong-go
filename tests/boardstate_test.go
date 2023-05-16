package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestBoardState(t *testing.T) {
	var seed = rand.Int63()

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
		boardState := game.GetPosBoardState(mahjong.East, posCalls[mahjong.East])
		events := game.GetPosEvents(mahjong.East, 0)
		//println(boardState.UTF8())

		if flag != mahjong.EndTypeRound {
			nb := mahjong.NewBoardState()
			nb.DecodeEvents(events)
			if !boardState.Equal(nb) {
				panic("boardState not equal")
			}
		}

		// test json
		bs, _ := json.Marshal(&boardState)
		//fmt.Println(string(bs))
		var board mahjong.BoardState
		err := json.Unmarshal(bs, &board)
		if err != nil {
			fmt.Println("error:", err)
		}
		posCalls, flag = game.Step(posCall)
		posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
	}
}
