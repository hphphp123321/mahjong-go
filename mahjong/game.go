package mahjong

import (
	"fmt"
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/hphphp123321/mahjong-go/common"
	"math/rand"
	"sort"
)

type Game struct {
	rule *Rule
	seed int64

	P0        *Player
	P1        *Player
	P2        *Player
	P3        *Player
	posPlayer map[Wind]*Player

	posEvents  map[Wind]Events
	ValidCalls map[Wind]Calls
	//posCall   map[Wind]*Call

	Tiles *MahjongTiles

	WindRound WindRound // {0:东一, 1:东二, 2:东三, 3:东四, 4:南一, 5:南二, 6:南三, 7:南四(all last), 8:西一(西入条件)...}
	NumGame   int

	NumRiichi int // riichi sticks num
	NumHonba  int // honba num

	nextRound bool // if wind round ++
	honbaPlus bool // if honba ++

	Position Wind

	State gameState
}

// NewMahjongGame
//
//	@Description: create a new game
//	@param playerSlice: player slice, len(playerSlice) == 4, playerSlice[0, 1, 2, 3] is east, south, west, north
//	@param seed: random seed
//	@param rule: game rule, nil for default rule
//	@return *Game
func NewMahjongGame(playerSlice []*Player, seed int64, rule *Rule) *Game {
	randP := rand.New(rand.NewSource(seed))
	game := Game{
		Tiles:     NewMahjongTiles(randP),
		WindRound: WindRoundDummy,
		seed:      seed,
	}
	game.Reset(playerSlice, nil)
	if rule == nil {
		game.rule = GetDefaultRule()
	} else {
		game.rule = rule
	}
	return &game
}

func (game *Game) Step() map[Wind]Calls {
	return game.State.step()
}

func (game *Game) Next(posCalls map[Wind]*Call) bool {
	if err := game.State.next(posCalls); err != nil {
		if err == ErrGameEnd {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func (game *Game) GetPosEvents(pos Wind, startIndex int) Events {
	return game.posEvents[pos][startIndex:]
}

func (game *Game) GetPosBoardState(pos Wind) (r *BoardState) {
	r = &BoardState{
		WindRound:      game.WindRound,
		NumHonba:       game.NumHonba,
		NumRiichi:      game.NumRiichi,
		DoraIndicators: game.Tiles.DoraIndicators(),
		PlayerWind:     pos,
		Position:       game.Position,
		HandTiles:      game.posPlayer[pos].HandTiles,
		//ValidActions:   nil,
		NumRemainTiles: game.Tiles.NumRemainTiles,
		PlayerEast: PlayerState{
			Points:         game.posPlayer[East].Points,
			Melds:          game.posPlayer[East].Melds,
			DiscardTiles:   game.posPlayer[East].DiscardTiles,
			TilesTsumoGiri: game.posPlayer[East].TilesTsumoGiri,
			IsRiichi:       game.posPlayer[East].IsRiichi,
		},
		PlayerSouth: PlayerState{
			Points:         game.posPlayer[South].Points,
			Melds:          game.posPlayer[South].Melds,
			DiscardTiles:   game.posPlayer[South].DiscardTiles,
			TilesTsumoGiri: game.posPlayer[South].TilesTsumoGiri,
			IsRiichi:       game.posPlayer[South].IsRiichi,
		},
		PlayerWest: PlayerState{
			Points:         game.posPlayer[West].Points,
			Melds:          game.posPlayer[West].Melds,
			DiscardTiles:   game.posPlayer[West].DiscardTiles,
			TilesTsumoGiri: game.posPlayer[West].TilesTsumoGiri,
			IsRiichi:       game.posPlayer[West].IsRiichi,
		},
		PlayerNorth: PlayerState{
			Points:         game.posPlayer[North].Points,
			Melds:          game.posPlayer[North].Melds,
			DiscardTiles:   game.posPlayer[North].DiscardTiles,
			TilesTsumoGiri: game.posPlayer[North].TilesTsumoGiri,
			IsRiichi:       game.posPlayer[North].IsRiichi,
		},
	}
	return r
}

func (game *Game) addPosEvent(posEvent map[Wind]Event) {
	for pos, event := range posEvent {
		game.posEvents[pos] = append(game.posEvents[pos], event)
	}
}

func (game *Game) NewGameRound() {
	game.NumGame += 1
	if game.nextRound {
		game.WindRound++
		game.nextRound = false
	}
	if game.honbaPlus {
		game.NumHonba++
		game.honbaPlus = false
	}
	game.Tiles.Reset()
	game.posEvents = map[Wind]Events{}
	game.Position = East
	game.posPlayer[Wind((16-game.WindRound)%4)] = game.P0
	game.posPlayer[Wind((17-game.WindRound)%4)] = game.P1
	game.posPlayer[Wind((18-game.WindRound)%4)] = game.P2
	game.posPlayer[Wind((19-game.WindRound)%4)] = game.P3
	game.P0.ResetForRound()
	game.P1.ResetForRound()
	game.P2.ResetForRound()
	game.P3.ResetForRound()
}

// Reset
//
//	@Description: reset game
//	@param playerSlice: player slice, len must be 4, and the order is East, South, West, North
//	@param tiles: tiles for game, if nil, will use default tiles
func (game *Game) Reset(playerSlice []*Player, tiles Tiles) {
	game.NumGame = -1
	game.WindRound = WindRoundEast1
	game.NumRiichi = 0
	game.NumHonba = 0
	game.Position = East
	game.State = &InitState{g: game, tiles: tiles}
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

func (game *Game) GetTileProcess(pMain *Player, tileID Tile) {
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
	if len(pMain.HandTiles) > 14 {
		panic("HandTiles len > 14")
	}
	pMain.JunNum++
	pMain.JunFuriten = false
}

// DealTile deal a tile
func (game *Game) DealTile(dealRinshan bool) Tile {
	return game.Tiles.DealTile(dealRinshan)
}

func (game *Game) DiscardTileProcess(pMain *Player, tileID Tile) {
	if !game.Tiles.allTiles[tileID].discardable {
		panic("Illegal Discard ID")
	}
	if tileID == pMain.HandTiles[len(pMain.HandTiles)-1] {
		pMain.TilesTsumoGiri = append(pMain.TilesTsumoGiri, true)
	} else {
		pMain.TilesTsumoGiri = append(pMain.TilesTsumoGiri, false)
	}
	pMain.HandTiles.Remove(tileID)
	if len(pMain.HandTiles) > 13 {
		panic("HandTiles len > 13")
	}
	pMain.DiscardTiles.Append(tileID)
	pMain.BoardTiles.Append(tileID)
	sort.Sort(&pMain.HandTiles)
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
	if len(validCalls) == 0 {
		panic("no valid action")
	}
	return validCalls
}

func (game *Game) JudgeSelfCalls(pMain *Player) Calls {
	var validCalls = make(Calls, 0)
	tsumo := game.judgeTsumo(pMain)
	riichi := game.judgeRiichi(pMain)
	shouMinKan := game.judgeShouMinKan(pMain)
	anKan := game.judgeAnKan(pMain)
	discard := game.JudgeDiscardCall(pMain)
	validCalls = append(validCalls, tsumo...)
	validCalls = append(validCalls, riichi...)
	validCalls = append(validCalls, shouMinKan...)
	validCalls = append(validCalls, anKan...)
	validCalls = append(validCalls, discard...)
	if len(validCalls) == 0 {
		panic("no valid action")
	}
	return validCalls
}

func (game *Game) JudgeOtherCalls(pMain *Player, tileID Tile) Calls {
	validCalls := Calls{SkipCall}
	daiMinKan := game.judgeDaiMinKan(pMain, tileID)
	pon := game.judgePon(pMain, tileID)
	chi := game.judgeChi(pMain, tileID)
	ron := game.judgeRon(pMain, tileID)
	validCalls = append(validCalls, daiMinKan...)
	validCalls = append(validCalls, pon...)
	validCalls = append(validCalls, chi...)
	validCalls = append(validCalls, ron...)
	if len(validCalls) == 1 {
		// only skip call
		return make(Calls, 0)
	}
	return validCalls
}

func (game *Game) processChanKan(pMain *Player, call *Call) *Result {
	winTile := call.CallTiles[0]
	pMain.IsChankan = true
	if pMain.IsRiichi && pMain.IppatsuStatus {
		pMain.IsIppatsu = true
	}
	result := game.getRonResult(pMain, winTile)
	if result == nil {
		panic("chan kan result error")
	}
	result.RonCall = call
	return result
}

func (game *Game) processTsumo(pMain *Player, call *Call) *Result {
	winTile := call.CallTiles[0]
	pMain.IsTsumo = true
	if game.Tiles.allTiles[winTile].isLast {
		pMain.IsHaitei = true
	}
	if game.Tiles.allTiles[winTile].isRinshan {
		pMain.IsRinshan = true
	}
	if pMain.IsRiichi && pMain.IppatsuStatus {
		pMain.IsIppatsu = true
	}
	if pMain.JunNum == 1 && pMain.IppatsuStatus {
		if pMain.Wind == 0 {
			pMain.IsTenhou = true
		} else {
			pMain.IsChiihou = true
		}
	}
	result := game.getRonResult(pMain, winTile)
	if result == nil {
		panic("tsumo result is nil")
	}
	result.RonCall = call
	return result
}

func (game *Game) processRon(pMain *Player, call *Call) *Result {
	winTile := call.CallTiles[0]
	if game.Tiles.allTiles[winTile].isLast {
		pMain.IsHoutei = true
	}
	if pMain.IsRiichi && pMain.IppatsuStatus {
		pMain.IsIppatsu = true
	}
	result := game.getRonResult(pMain, winTile)
	if result == nil {
		panic("ron result is nil")
	}
	result.RonCall = call
	return result
}

func (game *Game) processRiichiStep1(pMain *Player, call *Call) {
	if pMain.JunNum == 1 && pMain.IppatsuStatus {
		pMain.IsDaburuRiichi = true
	}
	pMain.RiichiStep = 1
	pMain.IsRiichi = true
	riichiTile := call.CallTiles[0]
	game.DiscardTileProcess(pMain, riichiTile)
	pMain.IppatsuStatus = true
}

func (game *Game) processRiichiStep2(pMain *Player) {
	pMain.RiichiStep = 2
	pMain.Points -= 1000
	game.NumRiichi++
}

func (game *Game) processChi(pMain *Player, call *Call) {
	pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles.Remove(call.CallTiles[1])
	tileID := call.CallTiles[2]
	subWind := call.CallTilesFromWho[2]
	game.posPlayer[subWind].BoardTiles.Remove(tileID)
	pMain.Melds = append(pMain.Melds, call)

	// 食替
	tileClass := tileID.Class()
	tile1Class := call.CallTiles[0].Class()
	tile2Class := call.CallTiles[1].Class()
	posClass := make(TileClasses, 0, 4)
	if tile1Class-tile2Class == 1 || tile2Class-tile1Class == 1 {
		if !common.Contain(tile1Class, TileClasses{0, 9, 18, 8, 17, 26}) &&
			!common.Contain(tile2Class, TileClasses{0, 9, 18, 8, 17, 26}) {
			posClass = TileClasses{tile1Class - 1, tile1Class + 1, tile2Class - 1, tile2Class + 1}
			posClass.Remove(tileClass)
			posClass.Remove(tile1Class)
			posClass.Remove(tile2Class)
		}
	}
	posClass.Append(tileClass)
	for _, tile := range pMain.HandTiles {
		if common.Contain(tile.Class(), posClass) {
			game.Tiles.allTiles[tile].discardable = false
		} else {
			game.Tiles.allTiles[tile].discardable = true
		}
	}
	pMain.JunNum++
}

func (game *Game) processPon(pMain *Player, call *Call) {
	pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles.Remove(call.CallTiles[1])
	tileID := call.CallTiles[2]
	subWind := call.CallTilesFromWho[2]
	game.posPlayer[subWind].BoardTiles.Remove(tileID)
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
	pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles.Remove(call.CallTiles[1])
	pMain.HandTiles.Remove(call.CallTiles[2])
	tileID := call.CallTiles[3]
	subWind := call.CallTilesFromWho[3]
	game.posPlayer[subWind].BoardTiles.Remove(tileID)

	pMain.Melds.Append(call)
}

func (game *Game) processAnKan(pMain *Player, call *Call) {
	pMain.HandTiles.Remove(call.CallTiles[0])
	pMain.HandTiles.Remove(call.CallTiles[1])
	pMain.HandTiles.Remove(call.CallTiles[2])
	pMain.HandTiles.Remove(call.CallTiles[3])
	pMain.Melds.Append(call)
}

func (game *Game) processShouMinKan(pMain *Player, call *Call) {
	tileID := call.CallTiles[3]
	pMain.HandTiles.Remove(tileID)
	for i, meld := range pMain.Melds {
		if meld.CallType != Pon {
			continue
		}
		if meld.CallTiles[0].Class() != tileID.Class() {
			continue
		}
		pMain.Melds[i] = call
		return
	}
	panic("ShouMinKan not success!")
}

func (game *Game) processKyuuShuKyuuHai() *Result {
	return &Result{
		RyuuKyokuReason: RyuuKyokuKyuuShuKyuuHai,
	}
}

func (game *Game) judgeRon(pMain *Player, tileID Tile) Calls {
	if pMain.IsFuriten() || !common.Contain(tileID.Class(), pMain.TenhaiSlice) {
		return make(Calls, 0)
	}
	if game.Tiles.allTiles[tileID].isLast {
		pMain.IsHoutei = true
	}
	result := game.getRonResult(pMain, tileID)
	if result == nil {
		pMain.IsHoutei = false
		return make(Calls, 0)
	}
	return Calls{&Call{
		CallType:         Ron,
		CallTiles:        Tiles{tileID, -1, -1, -1},
		CallTilesFromWho: []Wind{game.Tiles.allTiles[tileID].discardWind, WindDummy, WindDummy, WindDummy},
	}}
}

func (game *Game) judgeChanKan(pMain *Player, tileID Tile, isAnKan bool) Calls {
	if pMain.IsFuriten() || !common.Contain(tileID/4, pMain.TenhaiSlice) {
		return make(Calls, 0)
	}
	pMain.IsChankan = true
	result := game.getRonResult(pMain, tileID)
	if result == nil {
		pMain.IsChankan = false
		return make(Calls, 0)
	}
	if isAnKan &&
		common.Contain(yaku.YakumanKokushi, result.YakuResult) &&
		common.Contain(yaku.YakumanKokushi13, result.YakuResult) {
		pMain.IsChankan = false
		return make(Calls, 0)
	}
	return Calls{&Call{
		CallType:         ChanKan,
		CallTiles:        Tiles{tileID, -1, -1, -1},
		CallTilesFromWho: []Wind{WindDummy, WindDummy, WindDummy, WindDummy},
	}}
}

func (game *Game) judgeTsumo(pMain *Player) Calls {
	if pMain.ShantenNum != 0 {
		return make(Calls, 0)
	}
	winTile := pMain.HandTiles[len(pMain.HandTiles)-1]
	pMain.IsTsumo = true
	if game.Tiles.allTiles[winTile].isLast {
		pMain.IsHaitei = true
	}
	if game.Tiles.allTiles[winTile].isRinshan {
		pMain.IsRinshan = true
	}
	if pMain.JunNum == 1 && pMain.IppatsuStatus {
		if pMain.Wind == South {
			pMain.IsTenhou = true
		} else {
			pMain.IsChiihou = true
		}
	}
	result := game.getRonResult(pMain, winTile)
	if result == nil {
		pMain.IsTsumo = false
		pMain.IsHaitei = false
		pMain.IsRinshan = false
		pMain.IsTenhou = false
		pMain.IsChiihou = false
		return make(Calls, 0)
	}
	return Calls{&Call{
		CallType:         Tsumo,
		CallTiles:        Tiles{winTile, -1, -1, -1},
		CallTilesFromWho: []Wind{pMain.Wind, WindDummy, WindDummy, WindDummy},
	}}
}

func (game *Game) judgeRiichi(pMain *Player) Calls {
	if pMain.IsRiichi || (pMain.ShantenNum > 1 && pMain.JunNum > 1) || pMain.Points < 1000 || game.GetNumRemainTiles() < 4 {
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
			CallTilesFromWho: []Wind{pMain.Wind, WindDummy, WindDummy, WindDummy},
		})
	}
	return riichiCalls
}

func (game *Game) judgeChi(pMain *Player, tileID Tile) Calls {
	discardWind := game.Tiles.allTiles[tileID].discardWind
	chiClass := tileID.Class()
	if pMain.IsRiichi ||
		(pMain.Wind-discardWind+4)%4 != 1 ||
		chiClass >= 27 ||
		game.Tiles.allTiles[tileID].isLast ||
		game.GetNumRemainTiles() == 0 {
		return make(Calls, 0)
	}
	handTilesClass := pMain.GetHandTilesClass()
	if !(common.Contain(chiClass-1, handTilesClass) ||
		common.Contain(chiClass-2, handTilesClass) ||
		common.Contain(chiClass+1, handTilesClass) ||
		common.Contain(chiClass+2, handTilesClass)) {
		return make(Calls, 0)
	}
	var posCombinations []TileClasses
	if common.Contain(chiClass, TileClasses{0, 9, 18}) {
		posCombinations = append(posCombinations, TileClasses{chiClass + 1, chiClass + 2})
	} else if common.Contain(chiClass, TileClasses{1, 10, 19}) {
		posCombinations = append(posCombinations, TileClasses{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, TileClasses{chiClass + 1, chiClass + 2})
	} else if common.Contain(chiClass, TileClasses{7, 16, 25}) {
		posCombinations = append(posCombinations, TileClasses{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, TileClasses{chiClass - 2, chiClass - 1})
	} else if common.Contain(chiClass, TileClasses{8, 17, 26}) {
		posCombinations = append(posCombinations, TileClasses{chiClass - 2, chiClass - 1})
	} else {
		posCombinations = append(posCombinations, TileClasses{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, TileClasses{chiClass + 1, chiClass + 2})
		posCombinations = append(posCombinations, TileClasses{chiClass - 2, chiClass - 1})
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
			CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
		}
		posCalls.Append(&posCall)
		if common.Contain(tile1ID, Tiles{16, 52, 88}) {
			tile1Idx2 := handTilesClass.Index(tile1Class, tile1Idx1+1)
			if tile1Idx2 != -1 {
				tile1ID = pMain.HandTiles[tile1Idx2]
				posCall = Call{
					CallType:         Chi,
					CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
					CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
				}
				posCalls.Append(&posCall)
			}
		} else if common.Contain(tile2ID, Tiles{16, 52, 88}) {
			tile2Idx2 := handTilesClass.Index(tile2Class, tile2Idx1+1)
			if tile2Idx2 != -1 {
				tile2ID = pMain.HandTiles[tile2Idx2]
				posCall = Call{
					CallType:         Chi,
					CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
					CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
				}
				posCalls.Append(&posCall)
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
		handTilesCopy := pMain.HandTiles.Copy()
		handTilesCopy.Remove(tile1ID)
		handTilesCopy.Remove(tile2ID)
		tileClass := tile3ID.Class()
		tile1Class := tile1ID.Class()
		tile2Class := tile2ID.Class()
		posClass := make(TileClasses, 0, 4)
		if tile1Class-tile2Class == 1 || tile2Class-tile1Class == 1 {
			if !common.Contain(tile1Class, TileClasses{0, 9, 18, 8, 17, 26}) &&
				!common.Contain(tile2Class, TileClasses{0, 9, 18, 8, 17, 26}) {
				posClass = TileClasses{tile1Class - 1, tile1Class + 1, tile2Class - 1, tile2Class + 1}
				posClass.Remove(tileClass)
				posClass.Remove(tile1Class)
				posClass.Remove(tile2Class)
			}
		}
		posClass.Append(tileClass)
		flag := true
		for _, handTilesID := range handTilesCopy {
			if !common.Contain(handTilesID.Class(), posClass) {
				flag = false
				break
			}
		}
		if flag {
			delIdxSlice = append(delIdxSlice, i)
		}
	}
	if len(delIdxSlice) > 0 {
		posCalls = common.RemoveIndex(posCalls, delIdxSlice...)
	}
	return posCalls
}

func (game *Game) judgePon(pMain *Player, tileID Tile) Calls {
	if pMain.IsRiichi {
		return make(Calls, 0)
	}
	discardWind := game.Tiles.allTiles[tileID].discardWind
	ponClass := tileID.Class()
	tilesClass := pMain.GetHandTilesClass()
	tileCount := tilesClass.Count(ponClass)
	if tileCount < 2 || game.Tiles.allTiles[tileID].isLast || game.GetNumRemainTiles() == 0 {
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
		CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
	}
	posCalls = append(posCalls, &posCall)
	if tileCount == 3 {
		tile3Idx := tilesClass.Index(ponClass, tile2Idx+1)
		if tile3Idx == -1 {
			panic("no tile3")
		}
		tile3ID := pMain.HandTiles[tile3Idx]
		if common.Contain(tile1ID, TileClasses{16, 52, 88}) {
			posCall = Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile2ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
			}
			posCalls = append(posCalls, &posCall)
		} else if common.Contain(tile2ID, TileClasses{16, 52, 88}) {
			posCall = Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile1ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
			}
			posCalls = append(posCalls, &posCall)
		} else if common.Contain(tile3ID, TileClasses{16, 52, 88}) {
			posCall = Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile1ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
			}
			posCalls = append(posCalls, &posCall)
		}
	}
	return posCalls
}

func (game *Game) judgeDaiMinKan(pMain *Player, tileID Tile) Calls {
	discardWind := game.Tiles.allTiles[tileID].discardWind
	if pMain.IsRiichi || game.Tiles.allTiles[tileID].isLast || game.GetNumRemainTiles() == 0 {
		return make(Calls, 0)
	}
	kanClass := tileID.Class()
	class := pMain.GetHandTilesClass()
	tileCount := class.Count(kanClass)
	if tileCount < 2 {
		return make(Calls, 0)
	}
	posKanTiles := Tiles{Tile(kanClass * 4), Tile(kanClass*4 + 1), Tile(kanClass*4 + 2), Tile(kanClass*4 + 3)}
	posKanTiles.Remove(tileID)
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
	if len(pMain.HandTiles) == 2 || game.GetNumRemainTiles() == 0 {
		return make(Calls, 0)
	}
	tilesClass := pMain.GetHandTilesClass()
	var posClass TileClasses
	var posCalls = make(Calls, 0)
	for _, tileClass := range tilesClass {
		if common.Contain(tileClass, posClass) {
			continue
		}
		if tilesClass.Count(tileClass) == 4 {
			posClass.Append(tileClass)
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
			if pMain.IsRiichi {
				// judge ankan when player riichi
				// if a riichi player has 4 same tiles in hand but not draw the 4th tile, then this ankan is not valid
				if d != len(pMain.HandTiles)-1 {
					continue
				}
				// if a riichi player's tenhai changed after ankan, then this ankan is not valid
				tenhaiSlice := pMain.TenhaiSlice
				tmpHandTiles := common.RemoveIndex(pMain.HandTiles, a, b, c, d)
				melds := Calls{&posCall}
				tenhaiSliceAfterKan := GetTenhaiSlice(tmpHandTiles, melds)
				if !common.Equal(tenhaiSlice, tenhaiSliceAfterKan) {
					continue
				}
			}
			posCalls.Append(&posCall)
		}
	}
	return posCalls
}

func (game *Game) judgeShouMinKan(pMain *Player) Calls {
	if len(pMain.Melds) == 0 || game.GetNumRemainTiles() == 0 {
		return make(Calls, 0)
	}
	var posCalls Calls
	for _, call := range pMain.Melds {
		if call.CallType != Pon {
			continue
		}
		ponClass := call.CallTiles[0].Class()
		for _, tileID := range pMain.HandTiles {
			if tileID.Class() == ponClass && game.Tiles.allTiles[tileID].discardable {
				c := call.Copy()
				c.CallType = ShouMinKan
				c.CallTiles = append(c.CallTiles[:3], tileID)
				c.CallTilesFromWho = append(c.CallTilesFromWho[:3], pMain.Wind)
				posCalls = append(posCalls, c)
			}
		}
	}
	return posCalls
}

func (game *Game) judgeKyuShuKyuHai(pMain *Player) Calls {
	if pMain.JunNum > 1 || !pMain.RyuukyokuStatus {
		return make(Calls, 0)
	}
	kyuHai := make(map[TileClass]struct{})
	for _, tileID := range pMain.HandTiles {
		if common.Contain(tileID.Class(), TileClasses{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}) {
			kyuHai[tileID.Class()] = struct{}{}
		}
	}
	if len(kyuHai) < 9 {
		return make(Calls, 0)
	}
	return Calls{&Call{
		CallType:         KyuuShuKyuuHai,
		CallTiles:        Tiles{-1, -1, -1, -1},
		CallTilesFromWho: []Wind{WindDummy, WindDummy, WindDummy, WindDummy},
	}}
}

func (game *Game) getRonResult(pMain *Player, winTile Tile) (r *Result) {
	defer func() {
		if err := recover(); err != nil { // 如果recover返回为空，说明没报错
			r = nil
		}
	}()

	// remove win tile from hand tiles, because the calculator will add it back
	var handTiles Tiles
	for _, tile := range pMain.HandTiles {
		if tile != winTile {
			handTiles.Append(tile)
		}
	}

	ctx := &yaku.Context{
		Tile:        IntToInstance(int(winTile)),
		SelfWind:    base.Wind(pMain.Wind),
		RoundWind:   base.Wind(game.WindRound / 4),
		DoraTiles:   IntsToTiles(IndicatorsToDora(game.Tiles.DoraIndicators())),
		UraTiles:    IntsToTiles(IndicatorsToDora(game.Tiles.UraDoraIndicators())),
		Rules:       game.rule.YakuRule(),
		IsTsumo:     pMain.IsTsumo,
		IsRiichi:    pMain.IsRiichi,
		IsIpatsu:    pMain.IsIppatsu,
		IsDaburi:    pMain.IsDaburuRiichi,
		IsLastTile:  pMain.IsHaitei || pMain.IsHoutei,
		IsRinshan:   pMain.IsRinshan,
		IsFirstTake: pMain.IsTenhou,
		IsChankan:   pMain.IsChankan,
	}
	yakuResult := GetYakuResult(handTiles, pMain.Melds, ctx)
	if yakuResult == nil {
		return nil
	}
	scoreResult := score.GetScoreByResult(game.rule.ScoreRule(), yakuResult, score.Honba(game.NumHonba))
	return GenerateRonResult(yakuResult, &scoreResult)
}

func (game *Game) judgeSuuFonRenDa() bool {
	var tileClass = game.posPlayer[East].BoardTiles[len(game.posPlayer[South].BoardTiles)-1].Class()
	if tileClass != 27 && tileClass != 28 && tileClass != 29 && tileClass != 30 {
		return false
	}
	for _, player := range game.posPlayer {
		if player.JunNum > 1 || player.JunNum == 0 {
			return false
		}
		if !player.IppatsuStatus {
			return false
		}
		if player.BoardTiles[len(player.BoardTiles)-1].Class() != tileClass {
			return false
		}
	}
	return true
}

func (game *Game) judgeSuuChaRiichi() bool {
	var riiChiNum int
	for _, player := range game.posPlayer {
		if player.IsRiichi {
			riiChiNum++
		}
	}
	if riiChiNum < 4 {
		return false
	}
	return true
}

func (game *Game) judgeSuuKaiKan() bool {
	if game.Tiles.kanNum < 4 {
		return false
	}
	for _, player := range game.posPlayer {
		if player.KanNum == 4 {
			return false
		}
	}
	return true
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

// judgeNagashiMangan judge nagashi mangan(ryuu kyoku mangan for 8000 points)
func (game *Game) judgeNagashiMangan() []Wind {
	if game.Tiles.NumRemainTiles != 0 {
		panic("the number of remain tiles is not 0")
	}
	var retSlice []Wind
	winds := common.SortMapByKey(game.posPlayer)
	for _, wind := range winds {
		player := game.posPlayer[wind]
		if player.judgeNagashiMangan() {
			retSlice = append(retSlice, wind)
		}
	}
	return retSlice
}

func (game *Game) CheckGameEnd() bool {
	for _, player := range game.posPlayer {
		if player.Points < 0 {
			// player is bankrupt
			return true
		}
	}
	if game.WindRound < WindRound(game.rule.GameLength) {
		return false
	} else if game.WindRound > WindRound(game.rule.GameLength+4) || game.WindRound > WindRoundNorth4 {
		return true
	} else {
		for _, player := range game.posPlayer {
			// the first player's points over 30000
			if player.Points >= 30000 {
				return true
			}
		}
	}
	return false
}

func (game *Game) processNagashiMangan(winds []Wind) {
	otherWinds := []Wind{0, 1, 2, 3}
	for _, wind := range winds {
		for i, v := range otherWinds {
			if v == wind {
				otherWinds = append(otherWinds[:i], otherWinds[i+1:]...)
			}
		}
	}
	if len(winds) == 1 {
		// only one player nagashi mangan
		wind := winds[0]
		if wind == East {
			game.posPlayer[East].Points += 12000 + game.NumHonba*300 + game.NumRiichi*1000
			for _, w := range otherWinds {
				game.posPlayer[w].Points -= 4000 + game.NumHonba*100
			}
		} else {
			game.posPlayer[wind].Points += 8000 + game.NumHonba*300 + game.NumRiichi*1000
			for _, w := range otherWinds {
				if w == East {
					game.posPlayer[w].Points -= 4000 + game.NumHonba*100
				} else {
					game.posPlayer[w].Points -= 2000 + game.NumHonba*100
				}
			}
		}
	} else if len(winds) == 2 {
		// two players nagashi mangan
		wind0 := winds[0]
		wind1 := winds[1]
		if wind0 == 0 {
			game.posPlayer[wind0].Points += 12000 - 4000 + game.NumHonba*300 + game.NumRiichi*1000
			game.posPlayer[wind1].Points += 8000 - 4000 + game.NumHonba*300
			for _, w := range otherWinds {
				game.posPlayer[w].Points -= 4000 + 2000 + game.NumHonba*300
			}
		} else {
			game.posPlayer[wind0].Points += 8000 - 2000 + game.NumHonba*300 + game.NumRiichi*1000
			game.posPlayer[wind1].Points += 8000 - 2000 + game.NumHonba*300
			for _, w := range otherWinds {
				if w == East {
					game.posPlayer[w].Points -= 4000 + 4000 + game.NumHonba*300
				} else {
					game.posPlayer[w].Points -= 2000 + 2000 + game.NumHonba*300
				}
			}
		}
	} else {
		// three players nagashi mangan
		wind0 := winds[0]
		wind1 := winds[1]
		wind2 := winds[2]
		otherWind := otherWinds[0]
		if wind0 == 0 {
			game.posPlayer[wind0].Points += 12000 - 4000 - 4000 + game.NumHonba*300 + game.NumRiichi*1000
			game.posPlayer[wind1].Points += 8000 - 4000 - 2000 + game.NumHonba*300
			game.posPlayer[wind2].Points += 8000 - 4000 - 2000 + game.NumHonba*300
			game.posPlayer[otherWind].Points -= 4000 + 2000 + 2000 + game.NumHonba*900
		} else {
			game.posPlayer[wind0].Points += 8000 - 2000 - 2000 + game.NumHonba*300 + game.NumRiichi*1000
			game.posPlayer[wind1].Points += 8000 - 2000 - 2000 + game.NumHonba*300
			game.posPlayer[wind2].Points += 8000 - 2000 - 2000 + game.NumHonba*300
			game.posPlayer[otherWind].Points -= 4000 + 4000 + 4000 + game.NumHonba*900
		}
	}
}

// judgeTenHaiWinds returns the winds of players who have ten hai in the ryuukyoku situation.
func (game *Game) judgeTenHaiWinds() []Wind {
	var retSlice = make([]Wind, 0, 4)
	winds := common.SortMapByKey(game.posPlayer)
	for _, wind := range winds {
		player := game.posPlayer[wind]
		if len(player.TenhaiSlice) > 0 {
			retSlice = append(retSlice, wind)
		}
	}

	return retSlice
}

func (game *Game) processNormalRyuuKyoku(winds []Wind) {
	if len(winds) == 0 {
		// no player tenhai
		return
	}
	otherWinds := []Wind{0, 1, 2, 3}
	for _, wind := range winds {
		for i, v := range otherWinds {
			if v == wind {
				otherWinds = append(otherWinds[:i], otherWinds[i+1:]...)
			}
		}
	}
	if len(winds) == 1 {
		// one player tenhai
		wind := winds[0]
		game.posPlayer[wind].Points += 3000
		for _, w := range otherWinds {
			game.posPlayer[w].Points -= 1000
		}
	} else if len(winds) == 2 {
		// two players tenhai
		wind0 := winds[0]
		wind1 := winds[1]
		game.posPlayer[wind0].Points += 1500
		game.posPlayer[wind1].Points += 1500
		for _, w := range otherWinds {
			game.posPlayer[w].Points -= 1500
		}
	} else {
		// three players tenhai
		wind0 := winds[0]
		wind1 := winds[1]
		wind2 := winds[2]
		otherWind := otherWinds[0]
		game.posPlayer[wind0].Points += 1000
		game.posPlayer[wind1].Points += 1000
		game.posPlayer[wind2].Points += 1000
		game.posPlayer[otherWind].Points -= 3000
	}
}

// processRonResult processes the result of ron.
func (game *Game) processRonResult(results map[Wind]*Result) {
	var totalPoints = 0
	// sort results
	var winds []Wind
	for wind := range results {
		winds = append(winds, wind)
	}
	sort.Slice(winds, func(i, j int) bool {
		return winds[i] < winds[j]
	})
	for i, wind := range winds {
		result := results[wind]
		var pointsChange ScoreChanges
		if i == 0 {
			// the first player get riichi bonus
			pointsChange = result.ScoreResult.GetChanges(wind, game.Position, game.NumRiichi)
		} else {
			pointsChange = result.ScoreResult.GetChanges(wind, game.Position, 0)
		}
		game.posPlayer[wind].Points += pointsChange.TotalWin()
		totalPoints += pointsChange.TotalPayed()
	}
	game.posPlayer[game.Position].Points -= totalPoints
}

func (game *Game) addRonEvents(results map[Wind]*Result) {
	if len(results) > 1 {
		fmt.Println(1)
	}
	var posEvent = make(map[Wind]Event)
	for wind, result := range results {
		if result.RonCall.CallType == Ron {
			for w := range game.posPlayer {
				posEvent[w] = &EventRon{
					Who:       wind,
					HandTiles: game.posPlayer[wind].HandTiles,
					WinTile:   result.RonCall.CallTiles[0],
					Result:    result,
				}
			}
		} else {
			// chan kan
			for w := range game.posPlayer {
				posEvent[w] = &EventChanKan{
					Who:       wind,
					HandTiles: game.posPlayer[wind].HandTiles,
					WinTile:   result.RonCall.CallTiles[0],
					Result:    result,
				}
			}
		}
		game.addPosEvent(posEvent)
		posEvent = make(map[Wind]Event)
	}
}

func (game *Game) processTsumoResult(wind Wind, result *Result) {
	otherWinds := []Wind{0, 1, 2, 3}
	for i, v := range otherWinds {
		if v == wind {
			otherWinds = append(otherWinds[:i], otherWinds[i+1:]...)
		}
	}
	pointsChange := result.ScoreResult.GetChanges(wind, wind, game.NumRiichi)
	game.posPlayer[wind].Points += int(pointsChange.TotalWin())
	if wind == East {
		// dealer tsumo
		for _, w := range otherWinds {
			game.posPlayer[w].Points -= int(result.ScoreResult.PayTsumoDealer)
		}
	} else {
		for _, w := range otherWinds {
			if w == East {
				game.posPlayer[w].Points -= int(result.ScoreResult.PayTsumoDealer)
			} else {
				game.posPlayer[w].Points -= int(result.ScoreResult.PayTsumo)
			}
		}
	}
}

func (game *Game) GetGlobalEvents() Events {
	var events = make(Events, 0)
	// add all tiles event
	events = append(events, &EventGlobalInit{
		AllTiles:  game.Tiles.tiles,
		WindRound: game.WindRound,
		Seed:      game.seed,
		NumGame:   game.NumGame,
		NumHonba:  game.NumHonba,
		NumRiichi: game.NumRiichi,
		Rule:      game.rule,
	})
	// add all players event
	all := 0
	index := 0
	winds := common.SortMapByKey(game.posPlayer)
	for all < 4 {
		for _, wind := range winds {
			es := game.posEvents[wind]
			if len(es) <= index {
				all++
				continue
			}
			e := es[index]
			switch e.GetType() {
			case EventTypeStart:
				index++
			case EventTypeGet:
				if e.(*EventGet).Tile == TileDummy {
					continue
				} else {
					events = append(events, e)
					index++
				}
			case EventTypeFuriten:
				index++
			default:
				events = append(events, e)
				index++
			}
		}
	}
	var initPoints = make(map[Wind]int)
	if events[len(events)-1].GetType() == EventTypeEnd {
		for _, wind := range winds {
			initPoints[wind] = game.posPlayer[wind].Points - events[len(events)-1].(*EventEnd).PointsChange[wind]
		}
	} else {
		for _, wind := range winds {
			initPoints[wind] = game.posPlayer[wind].Points
		}
	}
	events[0].(*EventGlobalInit).InitPoints = initPoints
	return events
}
