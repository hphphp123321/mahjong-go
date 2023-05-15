package mahjong

import (
	"github.com/hphphp123321/go-common"
	"sort"
)

type Player struct {
	Points          int
	Wind            Wind
	JunNum          int
	KanNum          int
	HandTiles       Tiles
	DiscardTiles    Tiles
	TilesTsumoGiri  []bool
	BoardTiles      Tiles
	Melds           Calls
	TenpaiTiles     Tiles
	ShantenNum      int
	TenpaiSlice     TileClasses
	JunFuriten      bool
	DiscardFuriten  bool
	RiichiFuriten   bool
	IppatsuStatus   bool
	RyuukyokuStatus bool
	FuritenStatus   bool
	IsTsumo         bool
	IsRiichi        bool
	IsIppatsu       bool
	IsRinshan       bool
	IsChankan       bool
	IsHaitei        bool
	IsHoutei        bool
	IsDaburuRiichi  bool
	IsTenhou        bool
	IsChiihou       bool
	RiichiStep      int
}

func NewMahjongPlayer() *Player {
	p := Player{}
	p.ResetForGame()
	return &p
}

//func (player *Player) GetMelds() calc.Melds {
//	return CallsToMelds(player.melds)
//}

func (player *Player) InitTilesWind(tiles Tiles, wind Wind) {
	player.HandTiles = tiles
	sort.Sort(&player.HandTiles)
	player.Wind = wind
}

func (player *Player) GetShantenNum() int {
	return CalculateShantenNum(player.HandTiles, player.Melds)
}

// GetTenpaiSlice
//
//	@Description: Get player's tenpai slice after discard a tile
//	@receiver player
//	@return []TileClass
func (player *Player) GetTenpaiSlice() []TileClass {
	return GetTenpaiSlice(player.HandTiles.Copy(), player.Melds.Copy())
}

// GetPossibleTenpaiTiles
//
//	@Description: Get player's possible tenpai tiles after get a tile
//	@receiver player
//	@return Tiles
func (player *Player) GetPossibleTenpaiTiles() Tiles {
	if player.ShantenNum > 1 && player.JunNum > 1 {
		return Tiles{}
	}
	rTiles := make(Tiles, 0, len(player.HandTiles))
	handTilesCopy := make(Tiles, len(player.HandTiles)-1, len(player.HandTiles))
	for _, tile := range player.HandTiles {
		handTilesCopy = append(handTilesCopy, -1)
		copy(handTilesCopy, player.HandTiles)
		handTilesCopy.Remove(tile)
		shantenNum := CalculateShantenNum(handTilesCopy, player.Melds)
		if shantenNum == 0 {
			rTiles = append(rTiles, tile)
		}
	}
	return rTiles
}

func (player *Player) GetHandTilesClass() TileClasses {
	tilesClass := make([]TileClass, 0, len(player.HandTiles))
	for _, tile := range player.HandTiles {
		tilesClass = append(tilesClass, tile.Class())
	}
	return tilesClass
}

func (player *Player) IsNagashiMangan() bool {
	if len(player.BoardTiles) != len(player.DiscardTiles) {
		return false
	}
	for _, tileID := range player.DiscardTiles {
		if !common.SliceContain(YaoKyuTileClasses, tileID.Class()) {
			return false
		}
	}
	return true
}

func (player *Player) IsFuriten() bool {
	return player.JunFuriten || player.RiichiFuriten || player.DiscardFuriten
}

func (player *Player) ResetForRound() {
	player.Wind = WindDummy
	player.JunNum = 0
	player.KanNum = 0
	player.HandTiles = make(Tiles, 0, 14)
	player.DiscardTiles = make(Tiles, 0, 25)
	player.TilesTsumoGiri = make([]bool, 0, 25)
	player.BoardTiles = make(Tiles, 0, 25)
	player.Melds = make(Calls, 0, 4)
	player.TenpaiTiles = make(Tiles, 0, 13)
	player.ShantenNum = 7
	player.TenpaiSlice = []TileClass{}
	player.JunFuriten = false
	player.DiscardFuriten = false
	player.RiichiFuriten = false
	player.IppatsuStatus = true
	player.RyuukyokuStatus = true
	player.FuritenStatus = false
	player.IsTsumo = false
	player.IsRiichi = false
	player.IsIppatsu = false
	player.IsRinshan = false
	player.IsChankan = false
	player.IsHaitei = false
	player.IsHoutei = false
	player.IsDaburuRiichi = false
	player.IsTenhou = false
	player.IsChiihou = false
	player.RiichiStep = 0
}

func (player *Player) ResetForGame() {
	player.Points = 25000
	player.ResetForRound()
}

func (player *Player) judgeNagashiMangan() bool {
	if len(player.BoardTiles) != len(player.DiscardTiles) {
		return false
	}
	for _, tileID := range player.DiscardTiles {
		if !common.SliceContain(YaoKyuTileClasses, tileID.Class()) {
			return false
		}
	}
	return true
}

func (player *Player) Copy() *Player {
	p := Player{}
	p.Points = player.Points
	p.Wind = player.Wind
	p.JunNum = player.JunNum
	p.KanNum = player.KanNum
	p.HandTiles = player.HandTiles.Copy()
	p.DiscardTiles = player.DiscardTiles.Copy()
	p.TilesTsumoGiri = player.TilesTsumoGiri
	p.BoardTiles = player.BoardTiles.Copy()
	p.Melds = player.Melds.Copy()
	p.TenpaiTiles = player.TenpaiTiles.Copy()
	p.ShantenNum = player.ShantenNum
	p.TenpaiSlice = player.TenpaiSlice
	p.JunFuriten = player.JunFuriten
	p.DiscardFuriten = player.DiscardFuriten
	p.RiichiFuriten = player.RiichiFuriten
	p.IppatsuStatus = player.IppatsuStatus
	p.RyuukyokuStatus = player.RyuukyokuStatus
	p.FuritenStatus = player.FuritenStatus
	p.IsTsumo = player.IsTsumo
	p.IsRiichi = player.IsRiichi
	p.IsIppatsu = player.IsIppatsu
	p.IsRinshan = player.IsRinshan
	p.IsChankan = player.IsChankan
	p.IsHaitei = player.IsHaitei
	p.IsHoutei = player.IsHoutei
	p.IsDaburuRiichi = player.IsDaburuRiichi
	p.IsTenhou = player.IsTenhou
	p.IsChiihou = player.IsChiihou
	p.RiichiStep = player.RiichiStep
	return &p
}
