package tests

import (
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestReConstruct(t *testing.T) {
	var seed int64 = 17
	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)
	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(seed, nil)
	r := rand.New(rand.NewSource(seed))

	for i := 0; i < 100; i++ {
		posCalls = game.Reset(players, nil)
		flag := mahjong.EndTypeNone

		for flag != mahjong.EndTypeGame {
			for wind, calls := range posCalls {
				posCall[wind] = calls[r.Intn(len(calls))]
			}
			posCalls, flag = game.Step(posCall)
			posCall = make(map[mahjong.Wind]*mahjong.Call, 4)

			events := game.GetGlobalEvents()
			if len(events) > 1 {
				//b, _ := json.Marshal(&events)
				//fmt.Println(string(b))
				pSlice := make([]*mahjong.Player, 4)
				for i := 0; i < 4; i++ {
					pSlice[i] = mahjong.NewMahjongPlayer()
				}

				cGame := mahjong.ReConstructGame(pSlice, events)
				//fmt.Println("re:   " + cGame.State.String() + " player: " + cGame.Position.String())
				fmt.Println("game: " + cGame.State.String())
			}
			if len(posCalls) == 4 {
				fmt.Println("next round")
			}
		}
	}
}
