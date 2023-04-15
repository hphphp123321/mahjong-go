package tests

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
	"testing"
)

func TestRule(t *testing.T) {
	var seed int64 = 14

	rule := &mahjong.Rule{
		GameLength:           1,
		IsOpenTanyao:         false,
		HasAkaDora:           false,
		RenhouLimit:          mahjong.LimitMangan,
		IsHaiteiFromLiveOnly: false,
		IsUra:                false,
		IsIpatsu:             true,
		IsGreenRequired:      false,
		IsRinshanFu:          false,
		IsManganRound:        false,
		IsKazoeYakuman:       false,
		IsDoubleYakumans:     false,
		IsYakumanSum:         false,
		HonbaValue:           100,
		IsSanChaHou:          false,
		IsNagashiMangan:      false,
	}

	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)

	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(seed, rule)
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
