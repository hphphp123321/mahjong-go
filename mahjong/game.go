package mahjong

import (
	"fmt"
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/hphphp123321/go-common"
	"math/rand"
	"sort"
)

type Game struct {
	Rule *Rule
	Seed int64

	P0 *Player
	P1 *Player
	P2 *Player
	P3 *Player

	PosPlayer map[Wind]*Player
	posEvents map[Wind]Events

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
//	@param Seed: random Seed
//	@param Rule: game Rule, nil for default Rule
//	@return *Game
func NewMahjongGame(seed int64, rule *Rule) *Game {
	randP := rand.New(rand.NewSource(seed))
	game := Game{
		Tiles:     NewMahjongTiles(randP),
		WindRound: WindRoundDummy,
		Seed:      seed,
	}
	if rule == nil {
		game.Rule = GetDefaultRule()
	} else {
		game.Rule = rule
	}
	return &game
}

// Reset
//
//	@Description: reset game for new game
//	@param playerSlice: player slice, len must be 4, and the order is East, South, West, North
//	@param tiles: tiles for game, if nil, will use default random tiles
//	@return map[Wind]Calls: calls for East player
func (game *Game) Reset(playerSlice []*Player, tiles Tiles) map[Wind]Calls {
	game.NumGame = -1
	game.WindRound = WindRoundEast1
	game.NumRiichi = 0
	game.NumHonba = 0
	game.Position = East
	game.State = &InitState{g: game, tiles: tiles}
	game.posEvents = map[Wind]Events{}
	game.honbaPlus = false
	game.nextRound = false

	game.P0 = playerSlice[0]
	game.P1 = playerSlice[1]
	game.P2 = playerSlice[2]
	game.P3 = playerSlice[3]

	game.P0.ResetForGame()
	game.P1.ResetForGame()
	game.P2.ResetForGame()
	game.P3.ResetForGame()
	game.PosPlayer = map[Wind]*Player{0: playerSlice[0], 1: playerSlice[1], 2: playerSlice[2], 3: playerSlice[3]}

	posCalls, _ := game.Step(make(map[Wind]*Call, 4))
	return posCalls
}

// Step
//
//	@Description: game step
//	@receiver game
//	@param map[Wind]*Call: player action, if len(posCall) == 0, game will auto call
//	@return map[Wind]Calls: player valid actions
//	@return EndType: game end type, EndTypeNone for not end, EndTypeRound for round end, EndTypeGame for game end
func (game *Game) Step(posCall map[Wind]*Call) (map[Wind]Calls, EndType) {
	var posCalls = make(map[Wind]Calls, 4)
	if len(posCall) == 0 {
		posCalls = game.State.step()
	}
	for len(posCalls) == 0 {
		if err := game.State.next(posCall); err != nil {
			if err == ErrGameEnd {
				return posCalls, EndTypeGame
			} else {
				panic(err)
			}
		}
		posCall = make(map[Wind]*Call, 4)
		posCalls = game.State.step()
	}
	if len(posCalls) == 4 {
		return posCalls, EndTypeRound
	}
	return posCalls, EndTypeNone
}

// GetPosEvents
//
//	@Description: get player events
//	@receiver game
//	@param pos: player wind
//	@param startIndex: start index, 0 for all events
//	@return Events: player events
func (game *Game) GetPosEvents(pos Wind, startIndex int) Events {
	return game.posEvents[pos][startIndex:]
}

// GetPosBoardState
//
//	@Description: get player's wind board state
//	@receiver game
//	@param pos: player wind
//	@return r: board state
func (game *Game) GetPosBoardState(pos Wind, validActions Calls) (r *BoardState) {
	playerStates := make(map[Wind]*PlayerState, 4)
	for wind, p := range game.PosPlayer {
		playerStates[wind] = &PlayerState{
			Points:         p.Points,
			Melds:          p.Melds,
			DiscardTiles:   p.DiscardTiles,
			TilesTsumoGiri: p.TilesTsumoGiri,
			IsRiichi:       p.IsRiichi,
		}
	}
	r = &BoardState{
		WindRound:      game.WindRound,
		NumHonba:       game.NumHonba,
		NumRiichi:      game.NumRiichi,
		DoraIndicators: game.Tiles.DoraIndicators(),
		PlayerWind:     pos,
		Position:       game.Position,
		HandTiles:      game.PosPlayer[pos].HandTiles,
		ValidActions:   validActions,
		NumRemainTiles: game.Tiles.NumRemainTiles,
		PlayerStates:   playerStates,
	}
	return r
}

func (game *Game) addPosEvent(posEvent map[Wind]Event) {
	for pos, event := range posEvent {
		if event == nil {
			panic("event is nil")
		}
		game.posEvents[pos] = append(game.posEvents[pos], event)
	}
}

func (game *Game) newGameRound() {
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
	game.PosPlayer[Wind((16-game.WindRound)%4)] = game.P0
	game.PosPlayer[Wind((17-game.WindRound)%4)] = game.P1
	game.PosPlayer[Wind((18-game.WindRound)%4)] = game.P2
	game.PosPlayer[Wind((19-game.WindRound)%4)] = game.P3
	game.P0.ResetForRound()
	game.P1.ResetForRound()
	game.P2.ResetForRound()
	game.P3.ResetForRound()
}

func (game *Game) getTileProcess(pMain *Player, tileID Tile) {
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

// dealTile deal a tile
func (game *Game) dealTile(dealRinshan bool) Tile {
	return game.Tiles.DealTile(dealRinshan)
}

func (game *Game) discardTileProcess(pMain *Player, tileID Tile) {
	if !game.Tiles.allTiles[tileID].discardable {
		panic("Illegal Discard ID")
	}
	if tileID == pMain.HandTiles[len(pMain.HandTiles)-1] {
		pMain.TilesTsumoGiri = append(pMain.TilesTsumoGiri, true)
	} else {
		pMain.TilesTsumoGiri = append(pMain.TilesTsumoGiri, false)
	}
	pMain.HandTiles.Remove(tileID)
	pMain.DiscardTiles.Append(tileID)
	pMain.BoardTiles.Append(tileID)
	sort.Sort(&pMain.HandTiles)
	pMain.IppatsuStatus = false
	game.Tiles.allTiles[tileID].discardWind = pMain.Wind
	pMain.ShantenNum = pMain.GetShantenNum()
	if pMain.ShantenNum == 0 {
		pMain.TenpaiSlice = pMain.GetTenpaiSlice()
		flag := false
		for _, tile := range pMain.DiscardTiles {
			if common.SliceContain(pMain.TenpaiSlice, tile.Class()) {
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
		if common.SliceContain(game.PosPlayer[wind].TenpaiSlice, tileID.Class()) {
			game.PosPlayer[wind].FuritenStatus = true
		} else {
			game.PosPlayer[wind].FuritenStatus = false
		}
	}
}

func (game *Game) getDiscardableSlice(handTiles Tiles) Tiles {
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
	discardableSlice := game.getDiscardableSlice(pMain.HandTiles)
	for _, tile := range discardableSlice {
		call := Call{
			CallType:         Discard,
			CallTiles:        Tiles{tile, TileDummy, TileDummy, TileDummy},
			CallTilesFromWho: []Wind{pMain.Wind, WindDummy, WindDummy, WindDummy},
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
	game.discardTileProcess(pMain, riichiTile)
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
	game.PosPlayer[subWind].BoardTiles.Remove(tileID)
	pMain.Melds = append(pMain.Melds, call)

	// 食替
	tileClass := tileID.Class()
	tile1Class := call.CallTiles[0].Class()
	tile2Class := call.CallTiles[1].Class()
	posClass := make(TileClasses, 0, 4)
	if tile1Class-tile2Class == 1 || tile2Class-tile1Class == 1 {
		if !common.SliceContain(TileClasses{0, 9, 18, 8, 17, 26}, tile1Class) &&
			!common.SliceContain(TileClasses{0, 9, 18, 8, 17, 26}, tile2Class) {
			posClass = TileClasses{tile1Class - 1, tile1Class + 1, tile2Class - 1, tile2Class + 1}
			posClass.Remove(tileClass)
			posClass.Remove(tile1Class)
			posClass.Remove(tile2Class)
		}
	}
	posClass.Append(tileClass)
	for _, tile := range pMain.HandTiles {
		if common.SliceContain(posClass, tile.Class()) {
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
	game.PosPlayer[subWind].BoardTiles.Remove(tileID)
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
	game.PosPlayer[subWind].BoardTiles.Remove(tileID)

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
	if pMain.IsFuriten() || !common.SliceContain(pMain.TenpaiSlice, tileID.Class()) {
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
		CallTiles:        Tiles{tileID, TileDummy, TileDummy, TileDummy},
		CallTilesFromWho: []Wind{game.Tiles.allTiles[tileID].discardWind, WindDummy, WindDummy, WindDummy},
	}}
}

func (game *Game) judgeChanKan(pMain *Player, tileID Tile, isAnKan bool) Calls {
	if pMain.IsFuriten() || !common.SliceContain(pMain.TenpaiSlice, tileID.Class()) {
		return make(Calls, 0)
	}
	pMain.IsChankan = true
	result := game.getRonResult(pMain, tileID)
	if result == nil {
		pMain.IsChankan = false
		return make(Calls, 0)
	}
	if isAnKan &&
		common.SliceContain(result.YakuResult.Yakumans, YakumanKokushi) &&
		common.SliceContain(result.YakuResult.Yakumans, YakumanKokushi13) {
		pMain.IsChankan = false
		return make(Calls, 0)
	}
	return Calls{&Call{
		CallType:         ChanKan,
		CallTiles:        Tiles{tileID, TileDummy, TileDummy, TileDummy},
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
		CallTiles:        Tiles{winTile, TileDummy, TileDummy, TileDummy},
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
	tiles := pMain.GetPossibleTenpaiTiles()
	riichiCalls := make(Calls, 0, len(tiles))
	for _, tileID := range tiles {
		riichiCalls = append(riichiCalls, &Call{
			CallType:         Riichi,
			CallTiles:        Tiles{tileID, TileDummy, TileDummy, TileDummy},
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
	if !(common.SliceContain(handTilesClass, chiClass-1) ||
		common.SliceContain(handTilesClass, chiClass-2) ||
		common.SliceContain(handTilesClass, chiClass+1) ||
		common.SliceContain(handTilesClass, chiClass+2)) {
		return make(Calls, 0)
	}
	var posCombinations []TileClasses
	if common.SliceContain(TileClasses{0, 9, 18}, chiClass) {
		posCombinations = append(posCombinations, TileClasses{chiClass + 1, chiClass + 2})
	} else if common.SliceContain(TileClasses{1, 10, 19}, chiClass) {
		posCombinations = append(posCombinations, TileClasses{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, TileClasses{chiClass + 1, chiClass + 2})
	} else if common.SliceContain(TileClasses{7, 16, 25}, chiClass) {
		posCombinations = append(posCombinations, TileClasses{chiClass - 1, chiClass + 1})
		posCombinations = append(posCombinations, TileClasses{chiClass - 2, chiClass - 1})
	} else if common.SliceContain(TileClasses{8, 17, 26}, chiClass) {
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
		if !(common.SliceContain(handTilesClass, tile1Class) && common.SliceContain(handTilesClass, tile2Class)) {
			continue
		}
		tile1Idx1 := handTilesClass.Index(tile1Class, 0)
		tile1ID := pMain.HandTiles[tile1Idx1]
		tile2Idx1 := handTilesClass.Index(tile2Class, 0)
		tile2ID := pMain.HandTiles[tile2Idx1]
		posCall := &Call{
			CallType:         Chi,
			CallTiles:        Tiles{tile1ID, tile2ID, tileID, TileDummy},
			CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
		}
		posCalls.Append(posCall)
		if common.SliceContain(Tiles{16, 52, 88}, tile1ID) {
			tile1Idx2 := handTilesClass.Index(tile1Class, tile1Idx1+1)
			if tile1Idx2 != -1 {
				tile1ID = pMain.HandTiles[tile1Idx2]
				posCall = &Call{
					CallType:         Chi,
					CallTiles:        Tiles{tile1ID, tile2ID, tileID, TileDummy},
					CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
				}
				posCalls.Append(posCall)
			}
		} else if common.SliceContain(Tiles{16, 52, 88}, tile2ID) {
			tile2Idx2 := handTilesClass.Index(tile2Class, tile2Idx1+1)
			if tile2Idx2 != -1 {
				tile2ID = pMain.HandTiles[tile2Idx2]
				posCall = &Call{
					CallType:         Chi,
					CallTiles:        Tiles{tile1ID, tile2ID, tileID, TileDummy},
					CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
				}
				posCalls.Append(posCall)
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
			if !common.SliceContain(TileClasses{0, 9, 18, 8, 17, 26}, tile1Class) &&
				!common.SliceContain(TileClasses{0, 9, 18, 8, 17, 26}, tile2Class) {
				posClass = TileClasses{tile1Class - 1, tile1Class + 1, tile2Class - 1, tile2Class + 1}
				posClass.Remove(tileClass)
				posClass.Remove(tile1Class)
				posClass.Remove(tile2Class)
			}
		}
		posClass.Append(tileClass)
		flag := true
		for _, handTilesID := range handTilesCopy {
			if !common.SliceContain(posClass, handTilesID.Class()) {
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
	posCall := &Call{
		CallType:         Pon,
		CallTiles:        Tiles{tile1ID, tile2ID, tileID, -1},
		CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
	}
	posCalls.Append(posCall)
	if tileCount == 3 {
		tile3Idx := tilesClass.Index(ponClass, tile2Idx+1)
		if tile3Idx == -1 {
			panic("no tile3")
		}
		tile3ID := pMain.HandTiles[tile3Idx]
		if common.SliceContain(Tiles{16, 52, 88}, tile1ID) {
			posCall = &Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile2ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
			}
			posCalls.Append(posCall)
		} else if common.SliceContain(Tiles{16, 52, 88}, tile2ID) {
			posCall = &Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile1ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
			}
			posCalls.Append(posCall)
		} else if common.SliceContain(Tiles{16, 52, 88}, tile3ID) {
			posCall = &Call{
				CallType:         Pon,
				CallTiles:        Tiles{tile1ID, tile3ID, tileID, -1},
				CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, discardWind, WindDummy},
			}
			posCalls.Append(posCall)
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
	posKanTiles := kanClass.To4Tiles()
	posKanTiles.Remove(tileID)
	tile0 := posKanTiles[0]
	tile1 := posKanTiles[1]
	tile2 := posKanTiles[2]
	if !common.SliceContain(pMain.HandTiles, tile0) ||
		!common.SliceContain(pMain.HandTiles, tile1) ||
		!common.SliceContain(pMain.HandTiles, tile2) {
		return make(Calls, 0)
	}
	var posCalls Calls
	posCall := &Call{
		CallType:         DaiMinKan,
		CallTiles:        Tiles{tile0, tile1, tile2, tileID},
		CallTilesFromWho: []Wind{pMain.Wind, pMain.Wind, pMain.Wind, discardWind},
	}
	posCalls.Append(posCall)
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
		if common.SliceContain(posClass, tileClass) {
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
			posCall := &Call{
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
				// if a riichi player's Tenpai changed after ankan, then this ankan is not valid
				TenpaiSlice := pMain.TenpaiSlice
				tmpHandTiles := common.RemoveIndex(pMain.HandTiles, a, b, c, d)
				melds := pMain.Melds.Copy()
				melds.Append(posCall)
				TenpaiSliceAfterKan := GetTenpaiSlice(tmpHandTiles, melds)
				if !common.SliceEqual(TenpaiSlice, TenpaiSliceAfterKan) {
					continue
				}
			}
			posCalls.Append(posCall)
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
		if common.SliceContain(YaoKyuTileClasses, tileID.Class()) {
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
		Rules:       game.Rule.YakuRule(),
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
	scoreResult := score.GetScoreByResult(game.Rule.ScoreRule(), yakuResult, score.Honba(game.NumHonba))
	return GenerateRonResult(yakuResult, &scoreResult)
}

func (game *Game) judgeSuuFonRenDa() bool {
	var tileClass = game.PosPlayer[East].BoardTiles[len(game.PosPlayer[South].BoardTiles)-1].Class()
	if tileClass != 27 && tileClass != 28 && tileClass != 29 && tileClass != 30 {
		return false
	}
	for _, player := range game.PosPlayer {
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
	for _, player := range game.PosPlayer {
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
	for _, player := range game.PosPlayer {
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
	for wind, player := range game.PosPlayer {
		if wind == game.Position {
			continue
		}
		player.IppatsuStatus = false
	}
}

func (game *Game) breakRyuukyoku() {
	for _, player := range game.PosPlayer {
		player.RyuukyokuStatus = false
	}
}

func (game *Game) getCurrentRiichiNum() int {
	num := 0
	for _, player := range game.PosPlayer {
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
	winds := common.SortMapByKey(game.PosPlayer)
	for _, wind := range winds {
		player := game.PosPlayer[wind]
		if player.judgeNagashiMangan() {
			retSlice = append(retSlice, wind)
		}
	}
	return retSlice
}

func (game *Game) CheckGameEnd() bool {
	for _, player := range game.PosPlayer {
		if player.Points < 0 {
			// player is bankrupt
			return true
		}
	}
	if game.WindRound < WindRound(game.Rule.GameLength) {
		return false
	} else if game.WindRound > WindRound(game.Rule.GameLength+4) || game.WindRound > WindRoundNorth4 {
		return true
	} else {
		for _, player := range game.PosPlayer {
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
			game.PosPlayer[East].Points += 12000 + game.NumHonba*300 + game.NumRiichi*1000
			for _, w := range otherWinds {
				game.PosPlayer[w].Points -= 4000 + game.NumHonba*100
			}
		} else {
			game.PosPlayer[wind].Points += 8000 + game.NumHonba*300 + game.NumRiichi*1000
			for _, w := range otherWinds {
				if w == East {
					game.PosPlayer[w].Points -= 4000 + game.NumHonba*100
				} else {
					game.PosPlayer[w].Points -= 2000 + game.NumHonba*100
				}
			}
		}
	} else if len(winds) == 2 {
		// two players nagashi mangan
		wind0 := winds[0]
		wind1 := winds[1]
		if wind0 == 0 {
			game.PosPlayer[wind0].Points += 12000 - 4000 + game.NumHonba*300 + game.NumRiichi*1000
			game.PosPlayer[wind1].Points += 8000 - 4000 + game.NumHonba*300
			for _, w := range otherWinds {
				game.PosPlayer[w].Points -= 4000 + 2000 + game.NumHonba*300
			}
		} else {
			game.PosPlayer[wind0].Points += 8000 - 2000 + game.NumHonba*300 + game.NumRiichi*1000
			game.PosPlayer[wind1].Points += 8000 - 2000 + game.NumHonba*300
			for _, w := range otherWinds {
				if w == East {
					game.PosPlayer[w].Points -= 4000 + 4000 + game.NumHonba*300
				} else {
					game.PosPlayer[w].Points -= 2000 + 2000 + game.NumHonba*300
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
			game.PosPlayer[wind0].Points += 12000 - 4000 - 4000 + game.NumHonba*300 + game.NumRiichi*1000
			game.PosPlayer[wind1].Points += 8000 - 4000 - 2000 + game.NumHonba*300
			game.PosPlayer[wind2].Points += 8000 - 4000 - 2000 + game.NumHonba*300
			game.PosPlayer[otherWind].Points -= 4000 + 2000 + 2000 + game.NumHonba*900
		} else {
			game.PosPlayer[wind0].Points += 8000 - 2000 - 2000 + game.NumHonba*300 + game.NumRiichi*1000
			game.PosPlayer[wind1].Points += 8000 - 2000 - 2000 + game.NumHonba*300
			game.PosPlayer[wind2].Points += 8000 - 2000 - 2000 + game.NumHonba*300
			game.PosPlayer[otherWind].Points -= 4000 + 4000 + 4000 + game.NumHonba*900
		}
	}
}

// judgeTenpaiWinds returns the winds of players who have ten hai in the ryuukyoku situation.
func (game *Game) judgeTenpaiWinds() []Wind {
	var retSlice = make([]Wind, 0, 4)
	winds := common.SortMapByKey(game.PosPlayer)
	for _, wind := range winds {
		player := game.PosPlayer[wind]
		if len(player.TenpaiSlice) > 0 {
			retSlice = append(retSlice, wind)
		}
	}

	return retSlice
}

func (game *Game) processNormalRyuuKyoku(winds []Wind) {
	if len(winds) == 0 {
		// no player Tenpai
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
		// one player Tenpai
		wind := winds[0]
		game.PosPlayer[wind].Points += 3000
		for _, w := range otherWinds {
			game.PosPlayer[w].Points -= 1000
		}
	} else if len(winds) == 2 {
		// two players Tenpai
		wind0 := winds[0]
		wind1 := winds[1]
		game.PosPlayer[wind0].Points += 1500
		game.PosPlayer[wind1].Points += 1500
		for _, w := range otherWinds {
			game.PosPlayer[w].Points -= 1500
		}
	} else {
		// three players Tenpai
		wind0 := winds[0]
		wind1 := winds[1]
		wind2 := winds[2]
		otherWind := otherWinds[0]
		game.PosPlayer[wind0].Points += 1000
		game.PosPlayer[wind1].Points += 1000
		game.PosPlayer[wind2].Points += 1000
		game.PosPlayer[otherWind].Points -= 3000
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
		game.PosPlayer[wind].Points += pointsChange.TotalWin()
		totalPoints += pointsChange.TotalPayed()
	}
	game.PosPlayer[game.Position].Points -= totalPoints
	game.NumRiichi = 0 // Clear Riichi Sticks
}

func (game *Game) addRonEvents(results map[Wind]*Result) {
	if len(results) > 1 {
		fmt.Println(1)
	}
	var posEvent = make(map[Wind]Event)
	for wind, result := range results {
		if result.RonCall.CallType == Ron {
			for w := range game.PosPlayer {
				posEvent[w] = &EventRon{
					Who:       wind,
					FromWho:   result.RonCall.CallTilesFromWho[0],
					HandTiles: game.PosPlayer[wind].HandTiles,
					WinTile:   result.RonCall.CallTiles[0],
					Result:    result,
				}
			}
		} else {
			// chan kan
			for w := range game.PosPlayer {
				posEvent[w] = &EventChanKan{
					Who:       wind,
					FromWho:   result.RonCall.CallTilesFromWho[0],
					HandTiles: game.PosPlayer[wind].HandTiles,
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
	game.PosPlayer[wind].Points += int(pointsChange.TotalWin())
	if wind == East {
		// dealer tsumo
		for _, w := range otherWinds {
			game.PosPlayer[w].Points -= int(result.ScoreResult.PayTsumoDealer)
		}
	} else {
		for _, w := range otherWinds {
			if w == East {
				game.PosPlayer[w].Points -= int(result.ScoreResult.PayTsumoDealer)
			} else {
				game.PosPlayer[w].Points -= int(result.ScoreResult.PayTsumo)
			}
		}
	}
	game.NumRiichi = 0 // Clear Riichi Sticks
}

func (game *Game) GetGlobalEvents() Events {
	var events = make(Events, 0)
	// add all tiles event
	events = append(events, &EventGlobalInit{
		AllTiles:  game.Tiles.tiles,
		WindRound: game.WindRound,
		Seed:      game.Seed,
		NumGame:   game.NumGame,
		NumHonba:  game.NumHonba,
		NumRiichi: game.NumRiichi,
		Rule:      game.Rule,
	})
	// add all players event
	windIndex := map[Wind]int{
		East:  0,
		South: 0,
		West:  0,
		North: 0,
	}
	winds := []Wind{East, South, West, North}
	esEast := game.posEvents[East]
	esSouth := game.posEvents[South]
	esWest := game.posEvents[West]
	esNorth := game.posEvents[North]

	for windIndex[East] < len(esEast) &&
		windIndex[South] < len(esSouth) &&
		windIndex[West] < len(esWest) &&
		windIndex[North] < len(esNorth) {

		eastIndex := windIndex[East]
		southIndex := windIndex[South]
		westIndex := windIndex[West]
		northIndex := windIndex[North]

		eastEvent := esEast[eastIndex]
		southEvent := esSouth[southIndex]
		westEvent := esWest[westIndex]
		northEvent := esNorth[northIndex]

		var uniqueWind = WindDummy

		if eastEvent.GetType() == southEvent.GetType() {
			if eastEvent.GetType() == westEvent.GetType() {
				if eastEvent.GetType() != northEvent.GetType() {
					uniqueWind = North
				}
			} else {
				uniqueWind = West
			}
		} else {
			if eastEvent.GetType() == westEvent.GetType() {
				uniqueWind = South
			} else {
				uniqueWind = East
			}
		}

		if esEast[eastIndex].GetType() == EventTypeStart {
			windIndex[East]++
			windIndex[South]++
			windIndex[West]++
			windIndex[North]++
			continue
		}

		if uniqueWind != WindDummy {
			windIndex[uniqueWind]++
		} else {
			if eastEvent.GetType() == EventTypeGet {
				for _, e := range []Event{eastEvent, southEvent, westEvent, northEvent} {
					if e.(*EventGet).Tile != TileDummy {
						events = append(events, e)
						break
					}
				}
			} else {
				events = append(events, eastEvent)
			}
			windIndex[East]++
			windIndex[South]++
			windIndex[West]++
			windIndex[North]++
		}
	}

	var initPoints = make(map[Wind]int)
	if events[len(events)-1].GetType() == EventTypeEnd {
		for _, wind := range winds {
			initPoints[wind] = game.PosPlayer[wind].Points - events[len(events)-1].(*EventEnd).PointsChange[wind]
		}
	} else {
		for _, wind := range winds {
			initPoints[wind] = game.PosPlayer[wind].Points
		}
	}
	events[0].(*EventGlobalInit).InitPoints = initPoints
	return events
}
