package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestBoardState(t *testing.T) {
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
	eventIndex := 0
	for flag != mahjong.EndTypeGame {
		for wind, calls := range posCalls {
			posCall[wind] = calls[rand.Intn(len(calls))]
		}
		posCalls, flag = game.Step(posCall)
		posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
		eventIndex++
		boardState := game.GetPosBoardState(mahjong.East, nil)
		bs, _ := json.Marshal(&boardState)
		fmt.Println(string(bs))
		var board mahjong.BoardState
		err := json.Unmarshal(bs, &board)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}
