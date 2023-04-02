package tests

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"testing"
)

func TestTsumo(t *testing.T) {
	player := mahjong.NewMahjongPlayer()
	game := mahjong.NewMahjongGame([]*mahjong.Player{
		player,
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer(),
		mahjong.NewMahjongPlayer(),
	}, 1, mahjong.GetDefaultRule())

	player.ShantenNum = 0
	player.HandTiles = mahjong.Tiles{1, 2, 3, 5, 6, 7, 9, 10, 11, 13, 14, 15, 17, 110}
	calls := game.JudgeTsumo(player)
	println(calls)
}
