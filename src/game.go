package mahjong

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/hphphp123321/mahjong-goserver/common"
	"sort"
)

type Game struct {
	rule *Rule

	P0        *Player
	P1        *Player
	P2        *Player
	P3        *Player
	posPlayer map[Wind]*Player

	posEvents map[Wind]Events
	posCall   map[Wind]*Call

	Tiles *MahjongTiles

	WindRound int // {0:东一, 1:东二, 2:东三, 3:东四, 4:南一, 5:南二, 6:南三, 7:南四(all last), 8:西一(西入条件)...}
	NumGame   int
	NumRiichi int
	NumHonba  int

	Position Wind

	State gameState
}

func NewMahjongGame(playerSlice []*Player, rule *Rule) *Game {
	game := Game{Tiles: NewMahjongTiles()}
	game.Reset(playerSlice)
	if rule == nil {
		game.rule = GetDefaultRule()
	} else {
		game.rule = rule
	}
	return &game
}

func (game *Game) NewGameRound(windRound int) {
	game.WindRound = windRound
	game.NumGame += 1
	game.Tiles.Reset()
	game.Position = 0
	game.posEvents = map[Wind]Events{}
	game.posPlayer[Wind((16-windRound)%4)] = game.P0
	game.posPlayer[Wind((17-windRound)%4)] = game.P1
	game.posPlayer[Wind((18-windRound)%4)] = game.P2
	game.posPlayer[Wind((19-windRound)%4)] = game.P3
	game.P0.ResetForRound()
	game.P1.ResetForRound()
	game.P2.ResetForRound()
	game.P3.ResetForRound()
}

func (game *Game) Reset(playerSlice []*Player) {
	game.Tiles.Reset()
	game.NumGame = 0
	game.WindRound = 0
	game.NumRiichi = 0
	game.NumHonba = 0
	game.Position = 0
	game.posEvents = map[Wind]Events{}

	game.P0 = playerSlice[0]
	game.P1 = playerSlice[1]
	game.P2 = playerSlice[2]
	game.P3 = playerSlice[3]

	game.P0.ResetForGame()
	game.P1.ResetForGame()
	game.P2.ResetForGame()
	game.P3.ResetForGame()
	game.posPlayer = map[Wind]*Player{0: playerSlice[0], 1: playerSlice[1], 2: playerSlice[2], 3: playerSlice[3]}
}

func (game *Game) ProcessOtherCall(pMain *Player, call *Call) {
	switch call.CallType {
	case Skip:
		return
	case Chi:
		game.processChi(pMain, call)
		game.Position = pMain.Wind
		game.breakIppatsu()
		game.breakRyuukyoku()
	case Pon:
		game.processPon(pMain, call)
		game.Position = pMain.Wind
		game.breakIppatsu()
		game.breakRyuukyoku()
	case DaiMinKan:
		game.processDaiMinKan(pMain, call)
		game.Position = pMain.Wind
		game.breakIppatsu()
		game.breakRyuukyoku()
	default:
		panic("unknown call type")
	}
}

func (game *Game) ProcessSelfCall(pMain *Player, call *Call) {
	switch call.CallType {
	case Discard:
		game.DiscardTileProcess(pMain, call.CallTiles[0])
	case ShouMinKan:
		game.processShouMinKan(pMain, call)
		game.breakIppatsu()
		game.breakRyuukyoku()
	case AnKan:
		game.processAnKan(pMain, call)
		game.breakIppatsu()
		game.breakRyuukyoku()
	case Riichi:
		game.processRiichi(pMain, call)
		game.breakIppatsu()
	default:
		panic("unknown call type")
	}
}

func (game *Game) GetTileProcess(pMain *Player, tileID int) {
	if pMain.IsRiichi {
		for _, tile := range pMain.HandTiles {
			game.Tiles.allTiles[tile].discardable = false
		}
	} else {
		for _, tile := range pMain.HandTiles {
			game.Tiles.allTiles[tile].discardable = true
		}
	}
	pMain.HandTiles = append(pMain.HandTiles, tileID)
	pMain.JunNum++
	pMain.JunFuriten = false
}

// DealTile TODO deal tile
func (game *Game) DealTile() {
	game.Tiles.DealTile(false)
}

func (game *Game) DiscardTileProcess(pMain *Player, tileID int) {
	if !game.Tiles.allTiles[tileID].discardable {
		panic("Illegal Discard ID")
	}
	if tileID == pMain.HandTiles[len(pMain.HandTiles)-1] {
		pMain.TilesTsumoGiri = append(pMain.TilesTsumoGiri, 1)
	} else {
		pMain.TilesTsumoGiri = append(pMain.TilesTsumoGiri, 0)
	}
	pMain.HandTiles = pMain.HandTiles.Remove(tileID)
	pMain.DiscardTiles = append(pMain.DiscardTiles, tileID)
	pMain.BoardTiles = append(pMain.BoardTiles, tileID)
	sort.Ints(pMain.HandTiles)
	pMain.IppatsuStatus = false
	game.Tiles.allTiles[tileID].discardWind = pMain.Wind
	pMain.ShantenNum = pMain.GetShantenNum()
	if pMain.ShantenNum == 0 {
		pMain.TenhaiSlice = pMain.GetTenhaiSlice()
		flag := false
		for _, tile := range pMain.DiscardTiles {
			if common.Contain(tile/4, pMain.TenhaiSlice) {
				flag = true
			}
		}
		if flag {
			pMain.DiscardFuriten = true
		} else {
			pMain.DiscardFuriten = false
		}
	}
	otherWinds := game.getOtherWinds()
	for _, wind := range otherWinds {
		if common.Contain(tileID/4, game.posPlayer[wind].TenhaiSlice) {
			game.posPlayer[wind].FuritenStatus = true
		} else {
			game.posPlayer[wind].FuritenStatus = false
		}
	}
}

func (game *Game) GetDiscardableSlice(handTiles Tiles) Tiles {
	tiles := Tiles{}
	for _, tile := range handTiles {
		if game.Tiles.allTiles[tile].discardable {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

func (game *Game) JudgeDiscardCall(pMain *Player) Calls {
	var validCalls = make(Calls, 0)
	discardableSlice := game.GetDiscardableSlice(pMain.HandTiles)
	for _, tile := range discardableSlice {
		call := Call{
			CallType:         Discard,
			CallTiles:        Tiles{tile, -1, -1, -1},
			CallTilesFromWho: []Wind{pMain.Wind, -1, -1, -1},
		}
		validCalls = append(validCalls, &call)
	}
	return validCalls
}

func (game *Game) JudgeSelfCalls(pMain *Player) Calls {
	var validCalls Calls
	riichi := game.judgeRiichi(pMain)
	shouMinKan := game.judgeShouMinKan(pMain)
	anKan := game.judgeAnKan(pMain)
	discard := game.JudgeDiscardCall(pMain)
	validCalls = append(validCalls, riichi...)
	validCalls = append(validCalls, shouMinKan...)
	validCalls = append(validCalls, anKan...)
	validCalls = append(validCalls, discard...)
	return validCalls
}

func (game *Game) JudgeOtherCalls(pMain *Player, tileID int) Calls {
	validCalls := Calls{&Call{
		CallType:         Skip,
		CallTiles:        Tiles{-1, -1, -1, -1},
		CallTilesFromWho: []Wind{-1, -1, -1, -1},
	}}
	daiMinKan := game.judgeDaiMinKan(pMain, tileID)
	pon := game.judgePon(pMain, tileID)
	chi := game.judgeChi(pMain, tileID)
	validCalls = append(validCalls, daiMinKan...)
	validCalls = append(validCalls, pon...)
	validCalls = append(validCalls, chi...)
	return validCalls
}

func (game *Game) processRiichi(pMain *Player, call *Call) {
	if pMain.JunNum == 1 && pMain.IppatsuStatus {
		pMain.IsDaburuRiichi = true
	}
	pMain.IsRiichi = true
	riichiTile := call.CallTiles[0]
	pMain.Points -= 1000
	game.DiscardTileProcess(pMain, riichiTile)
	pMain.IppatsuStatus = true
}

func (game *Game) processChi(pMain *Player, call *Call) {
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[1])
	tileID := call.CallTiles[2]
	subWind := call.CallTilesFromWho[2]
	game.posPlayer[subWind].BoardTiles = game.posPlayer[subWind].BoardTiles.Remove(tileID)
	pMain.Melds = append(pMain.Melds, call)

	// 食替
	tileClass := tileID / 4
	tile1Class := call.CallTiles[0] / 4
	tile2Class := call.CallTiles[1] / 4
	posClass := make(Tiles, 0, 4)
	if tile1Class-tile2Class == 1 || tile2Class-tile1Class == 1 {
		if !common.Contain(tile1Class, []int{0, 9, 18, 8, 17, 26}) && !common.Contain(tile2Class, []int{0, 9, 18, 8, 17, 26}) {
			posClass = Tiles{tile1Class - 1, tile1Class + 1, tile2Class - 1, tile2Class + 1}
			posClass = posClass.Remove(tileClass)
			posClass = posClass.Remove(tile1Class)
			posClass = posClass.Remove(tile2Class)
		}
	}
	posClass = append(posClass, tileClass)
	for _, tile := range pMain.HandTiles {
		if common.Contain(tile/4, posClass) {
			game.Tiles.allTiles[tile].discardable = false
		} else {
			game.Tiles.allTiles[tile].discardable = true
		}
	}
	pMain.JunNum++
}

func (game *Game) processPon(pMain *Player, call *Call) {
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[1])
	tileID := call.CallTiles[2]
	subWind := call.CallTilesFromWho[2]
	game.posPlayer[subWind].BoardTiles = game.posPlayer[subWind].BoardTiles.Remove(tileID)
	pMain.Melds = append(pMain.Melds, call)
	// 食替
	tileClass := tileID / 4
	for _, tile := range pMain.HandTiles {
		if tile/4 == tileClass {
			game.Tiles.allTiles[tile].discardable = false
		} else {
			game.Tiles.allTiles[tile].discardable = true
		}
	}
	pMain.JunNum++
}

func (game *Game) processDaiMinKan(pMain *Player, call *Call) {
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[1])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[2])
	tileID := call.CallTiles[3]
	subWind := call.CallTilesFromWho[3]
	game.posPlayer[subWind].BoardTiles = game.posPlayer[subWind].BoardTiles.Remove(tileID)

	pMain.Melds = append(pMain.Melds, call)
}

func (game *Game) processAnKan(pMain *Player, call *Call) {
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[1])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[2])
	pMain.HandTiles = pMain.HandTiles.Remove(call.CallTiles[3])
	pMain.Melds = append(pMain.Melds, call)
}

func (game *Game) processShouMinKan(pMain *Player, call *Call) {
	tileID := call.CallTiles[3]
	pMain.HandTiles = pMain.HandTiles.Remove(tileID)
	for i, meld := range pMain.Melds {
		if meld.CallType != Pon {
			continue
		}
		if meld.CallTiles[0]/4 != tileID/4 {
			continue
		}
		pMain.Melds[i] = call
		return
	}
	panic("ShouMinKan not success!")
}

func (game *Game) judgeRiichi(pMain *Player) Calls {
	if pMain.IsRiichi || (pMain.ShantenNum > 1 && pMain.JunNum > 1) || pMain.Points < 1000 {
		return make(Calls, 0)
	}
	if len(pMain.HandTiles) != 14 {
		for _, meld := range pMain.Melds {
			if meld.CallType != AnKan {
				return make(Calls, 0)
			}
		}
	}
	tiles := pMain.GetRiichiTiles()
	riichiCalls := make(Calls, 0, len(tiles))
	for _, tileID := range tiles {
		riichiCalls = append(riichiCalls, &Call{
			CallType:         Riichi,
			CallTiles:        Tiles{tileID, -1, -1, -1},
			CallTilesFromWho: []Wind{pMain.Wind, -1, -1, -1},
		})
	}
	return riichiCalls
}

func (game *Game) judgeChi(pMain *Player, tileID int) Calls {
	discardWind := game.Tiles.allTiles[tileID].discardWind
	chiClass := tileID / 4
	if pMain.IsRiichi || (pMain.Wind-discardWind+4)%4 != 1 || chiClass > 27 || game.Tiles.allTiles[tileID].isLast {
		return make(Calls, 0)
	}
	handTilesClass := Tiles(pMain.GetHandTilesClass())
	if !(common.Contain(chiClass-1, handTilesClass) ||
		common.Contain(chiClass-2, handTilesClass) ||
		common.Contain(chiClass+1, handTilesClass) ||
		common.Contain(chiClass+2, handTilesClass)) {
		return make(Calls, 0)
	}
	var posCombinations [][]int
	if common.Contain(chiClass, []int{0, 9, 18}) {
		posCombinations = append(posCombinations, []int{chiClass + 1, chiClass + 2})
	} else if common.Contain(chiClass, []int{1, 10, 19}) {
		posCombinations = append(posCombinations, []int{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, []int{chiClass + 1, chiClass + 2})
	} else if common.Contain(chiClass, []int{7, 16, 25}) {
		posCombinations = append(posCombinations, []int{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, []int{chiClass - 2, chiClass - 1})
	} else if common.Contain(chiClass, []int{8, 17, 26}) {
		posCombinations = append(posCombinations, []int{chiClass - 2, chiClass - 1})
	} else {
		posCombinations = append(posCombinations, []int{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, []int{chiClass + 1, chiClass + 2})
		posCombinations = append(posCombinations, []int{chiClass - 2, chiClass - 1})
	}
	var posCalls Calls
	for _, posCom := range posCombinations {
		tile1Class := posCom[0]
		tile2Class := posCom[1]
		if !(common.Contain(tile1Class, handTilesClass) && common.Contain(tile2Class, handTilesClass)) {
			continue
		}
		tile1Idx1 := handTilesClass.Index(tile1Class, 0)
		tile1ID := pMain.HandTiles[tile1Idx1]
		tile2Idx1 := handTilesClass.Index(tile2Class, 0)
		tile2ID := pMain.HandTiles[tile2Idx1]
		posCall := Call{
			CallType:         Chi,
			CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
			CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
		}
		posCalls = append(posCalls, &posCall)
		if common.Contain(tile1ID, []int{16, 52, 88}) {
			tile1Idx2 := handTilesClass.Index(tile1Class, tile1Idx1+1)
			if tile1Idx2 != -1 {
				tile1ID = pMain.HandTiles[tile1Idx2]
				posCall = Call{
					CallType:         Chi,
					CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
					CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
				}
				posCalls = append(posCalls, &posCall)
			}
		} else if common.Contain(tile2ID, []int{16, 52, 88}) {
			tile2Idx2 := handTilesClass.Index(tile2Class, tile2Idx1+1)
			if tile2Idx2 != -1 {
				tile2ID = pMain.HandTiles[tile2Idx2]
				posCall = Call{
					CallType:         Chi,
					CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
					CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
				}
				posCalls = append(posCalls, &posCall)
			}
		}
	}
	// 食替
	if len(pMain.HandTiles) > 7 {
		return posCalls
	}
	delIdxSlice := make([]int, 0, len(posCalls))
	for i, call := range posCalls {
		tile1ID := call.CallTiles[0]
		tile2ID := call.CallTiles[1]
		tile3ID := call.CallTiles[2]
		handTilesCopy := make(Tiles, len(pMain.HandTiles), len(pMain.HandTiles))
		copy(handTilesCopy, pMain.HandTiles)
		handTilesCopy = handTilesCopy.Remove(tile1ID)
		handTilesCopy = handTilesCopy.Remove(tile2ID)
		tileClass := tile3ID / 4
		tile1Class := tile1ID / 4
		tile2Class := tile2ID / 4
		posClass := make(Tiles, 0, 4)
		if tile1Class-tile2Class == 1 || tile2Class-tile1Class == 1 {
			if !common.Contain(tile1Class, []int{0, 9, 18, 8, 17, 26}) && !common.Contain(tile2Class, []int{0, 9, 18, 8, 17, 26}) {
				posClass = Tiles{tile1Class - 1, tile1Class + 1, tile2Class - 1, tile2Class + 1}
				posClass = posClass.Remove(tileClass)
				posClass = posClass.Remove(tile1Class)
				posClass = posClass.Remove(tile2Class)
			}
		}
		posClass = append(posClass, tileClass)
		flag := true
		for _, handTIlesID := range pMain.HandTiles {
			if !common.Contain(handTIlesID/4, posClass) {
				flag = false
			}
		}
		if flag {
			delIdxSlice = append(delIdxSlice, i)
		}
	}
	for idx := len(delIdxSlice) - 1; idx >= 0; idx-- {
		delIdx := delIdxSlice[idx]
		posCalls = append(posCalls[:delIdx], posCalls[delIdx+1:]...)
	}
	return posCalls
}

func (game *Game) judgePon(pMain *Player, tileID int) Calls {
	if pMain.IsRiichi {
		return make(Calls, 0)
	}
	discardWind := game.Tiles.allTiles[tileID].discardWind
	ponClass := tileID / 4
	tilesClass := Tiles(pMain.GetHandTilesClass())
	tileCount := tilesClass.Count(ponClass)
	if tileCount < 2 || game.Tiles.allTiles[tileID].isLast {
		return make(Calls, 0)
	}
	var posCalls Calls
	tile1Idx := tilesClass.Index(ponClass, 0)
	tile1ID := pMain.HandTiles[tile1Idx]
	tile2Idx := tilesClass.Index(ponClass, tile1Idx+1)
	tile2ID := pMain.HandTiles[tile2Idx]
	posCall := Call{
		CallType:         Pon,
		CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
		CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
	}
	posCalls = append(posCalls, &posCall)
	if tileCount == 3 {
		tile3Idx := tilesClass.Index(ponClass, tile2Idx+1)
		if tile3Idx == -1 {
			panic("no tile3")
		}
		tile3ID := pMain.HandTiles[tile3Idx]
		if common.Contain(tile1ID, []int{16, 52, 88}) {
			posCall = Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile2ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
			}
			posCalls = append(posCalls, &posCall)
		} else if common.Contain(tile2ID, []int{16, 52, 88}) {
			posCall = Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile1ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
			}
			posCalls = append(posCalls, &posCall)
		} else if common.Contain(tile3ID, []int{16, 52, 88}) {
			posCall = Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile1ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, -1},
			}
			posCalls = append(posCalls, &posCall)
		}
	}
	return posCalls
}

func (game *Game) judgeDaiMinKan(pMain *Player, tileID int) Calls {
	discardWind := game.Tiles.allTiles[tileID].discardWind
	if pMain.IsRiichi || game.Tiles.allTiles[tileID].isLast || game.Tiles.NumRemainTiles == 0 {
		return make(Calls, 0)
	}
	kanClass := tileID / 4
	tileCount := Tiles(pMain.GetHandTilesClass()).Count(kanClass)
	if tileCount < 2 {
		return make(Calls, 0)
	}
	posKanTiles := Tiles{kanClass * 4, kanClass*4 + 1, kanClass*4 + 2, kanClass*4 + 3}.Remove(tileID)
	tile0 := posKanTiles[0]
	tile1 := posKanTiles[1]
	tile2 := posKanTiles[2]
	if !common.Contain(tile0, pMain.HandTiles) ||
		!common.Contain(tile1, pMain.HandTiles) ||
		!common.Contain(tile2, pMain.HandTiles) {
		return make(Calls, 0)
	}
	var posCalls Calls
	posCall := Call{
		CallType:         DaiMinKan,
		CallTiles:        Tiles{tile0, tile1, tile2, tileID},
		CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, pMain.Wind, discardWind},
	}
	posCalls = append(posCalls, &posCall)
	return posCalls
}

func (game *Game) judgeAnKan(pMain *Player) Calls {
	if len(pMain.HandTiles) == 2 || game.Tiles.NumRemainTiles == 0 {
		return make(Calls, 0)
	}
	tilesClass := Tiles(pMain.GetHandTilesClass())
	var posClass []int
	var posCalls = make(Calls, 0)
	for _, tileClass := range tilesClass {
		if common.Contain(tileClass, posClass) {
			continue
		}
		if tilesClass.Count(tileClass) == 4 {
			posClass = append(posClass, tileClass)
			a := tilesClass.Index(tileClass, 0)
			b := tilesClass.Index(tileClass, a+1)
			c := tilesClass.Index(tileClass, b+1)
			d := tilesClass.Index(tileClass, c+1)
			if a == -1 || b == -1 || c == -1 || d == -1 {
				panic("index error")
			}
			posCall := Call{
				CallType:         AnKan,
				CallTiles:        Tiles{pMain.HandTiles[a], pMain.HandTiles[b], pMain.HandTiles[c], pMain.HandTiles[d]},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, pMain.Wind, pMain.Wind},
			}
			posCalls = append(posCalls, &posCall)
		}
	}
	return posCalls
}

func (game *Game) judgeShouMinKan(pMain *Player) Calls {
	if len(pMain.Melds) == 0 || game.Tiles.NumRemainTiles == 0 {
		return make(Calls, 0)
	}
	var posCalls Calls
	for _, call := range pMain.Melds {
		if call.CallType != Pon {
			continue
		}
		ponClass := call.CallTiles[0] / 4
		for _, tileID := range pMain.HandTiles {
			if tileID/4 == ponClass && game.Tiles.allTiles[tileID].discardable {
				posCall := Call{
					CallType:         ShouMinKan,
					CallTiles:        append(call.CallTiles[:3], tileID),
					CallTilesFromWho: append(call.CallTilesFromWho[:3], pMain.Wind),
				}
				posCalls = append(posCalls, &posCall)
			}
		}
	}
	return posCalls
}

func (game *Game) getRonResult(pMain *Player, winTile int) (r *Result) {
	defer func() {
		if err := recover(); err != nil { // 如果recover返回为空，说明没报错
			r = nil
		}
	}()
	ctx := &yaku.Context{
		Tile:        IntToInstance(winTile),
		SelfWind:    base.Wind(pMain.Wind),
		RoundWind:   base.Wind(game.WindRound / 4),
		DoraTiles:   IntsToTiles(IndicatorsToDora(game.Tiles.DoraIndicators())),
		UraTiles:    IntsToTiles(IndicatorsToDora(game.Tiles.UraDoraIndicators())),
		Rules:       game.rule.yakuRule,
		IsTsumo:     pMain.IsTsumo,
		IsRiichi:    pMain.IsRiichi,
		IsIpatsu:    pMain.IsIppatsu,
		IsDaburi:    pMain.IsDaburuRiichi,
		IsLastTile:  pMain.IsHaitei || pMain.IsHoutei,
		IsRinshan:   pMain.IsRinshan,
		IsFirstTake: pMain.IsTenhou,
		IsChankan:   pMain.IsChankan,
	}
	yakuResult := GetYakuResult(pMain.HandTiles, pMain.Melds, ctx)
	if yakuResult == nil {
		return nil
	}
	scoreResult := score.GetScoreByResult(game.rule.scoreRule, yakuResult, score.Honba(game.NumHonba))
	return GenerateResult(yakuResult, &scoreResult)
}

func (game *Game) GetNumRemainTiles() int {
	return game.Tiles.NumRemainTiles
}

func (game *Game) getOtherWinds() []Wind {
	otherWinds := []Wind{0, 1, 2, 3}
	for i, v := range otherWinds {
		if v == game.Position {
			otherWinds = append(otherWinds[:i], otherWinds[i+1:]...)
			break
		}
	}
	return otherWinds
}

func (game *Game) breakIppatsu() {
	for wind, player := range game.posPlayer {
		if wind == game.Position {
			continue
		}
		player.IppatsuStatus = false
	}
}

func (game *Game) breakRyuukyoku() {
	for _, player := range game.posPlayer {
		player.RyuukyokuStatus = false
	}
}

func (game *Game) getCurrentRiichiNum() int {
	num := 0
	for _, player := range game.posPlayer {
		if player.IsRiichi {
			num++
		}
	}
	return num
}
