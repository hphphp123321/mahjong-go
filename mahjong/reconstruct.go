package mahjong

import (
	"errors"
	"github.com/hphphp123321/go-common"
)

// ReConstructGame
//
//	@Description: reconstruct game from globalEvents
//	@param playerSlice: player slice
//	@param globalEvents: global events
//	@return *Game
func ReConstructGame(playerSlice []*Player, globalEvents Events) *Game {
	game := &Game{}
	var posCalls = make(map[Wind]Calls)

	e := globalEvents[0]
	if e.GetType() != EventTypeGlobalInit {
		panic(errors.New("first event must be EventTypeGlobalInit"))
	}
	et := e.(*EventGlobalInit)
	game.Tiles = NewMahjongTiles(nil)
	posCalls = game.Reset(playerSlice, et.AllTiles)
	game.WindRound = et.WindRound
	game.Seed = et.Seed
	game.NumGame = et.NumGame
	game.NumHonba = et.NumHonba
	game.NumRiichi = et.NumRiichi
	game.Rule = et.Rule
	game.PosPlayer[Wind((16-game.WindRound)%4)].Points = et.InitPoints[Wind((16-game.WindRound)%4)]
	game.PosPlayer[Wind((17-game.WindRound)%4)].Points = et.InitPoints[Wind((17-game.WindRound)%4)]
	game.PosPlayer[Wind((18-game.WindRound)%4)].Points = et.InitPoints[Wind((18-game.WindRound)%4)]
	game.PosPlayer[Wind((19-game.WindRound)%4)].Points = et.InitPoints[Wind((19-game.WindRound)%4)]

	index := 1
	for index < len(globalEvents) {
		var posCall = make(map[Wind]*Call)
		if len(posCalls) == 0 {
			posCalls, _ = game.Step(posCall)
			index++
			continue
		}

		event := globalEvents[index]
		switch event.GetType() {
		case EventTypeGlobalInit:
			return ReConstructGame(playerSlice, globalEvents[index:]) // reconstruct from this event
		case EventTypeGet:
			for wind := range posCalls {
				if game.Position != wind || common.SliceContain(posCalls[wind], SkipCall) {
					posCall[wind] = SkipCall
				}
			}
		case EventTypeDiscard:
			who := event.(*EventDiscard).Who
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallType == Discard && call.CallTiles[0] == event.(*EventDiscard).Tile {
					posCall[who] = call
					break
				}
			}
		case EventTypeTsumoGiri:
			who := event.(*EventTsumoGiri).Who
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallType == Discard && call.CallTiles[0] == event.(*EventTsumoGiri).Tile {
					posCall[who] = call
					break
				}
			}
		case EventTypeAnKan:
			who := event.(*EventAnKan).Who
			calls := posCalls[who]
			for _, call := range calls {
				if CallEqual(call, event.(*EventAnKan).Call) {
					posCall[who] = call
					break
				}
			}
		case EventTypeShouMinKan:
			who := event.(*EventShouMinKan).Who
			calls := posCalls[who]
			for _, call := range calls {
				if CallEqual(call, event.(*EventShouMinKan).Call) {
					posCall[who] = call
					break
				}
			}
		case EventTypeRiichi:
			who := event.(*EventRiichi).Who
			step := event.(*EventRiichi).Step
			var tile Tile
			if step == 2 {
				index++
				continue
			}
			nextEvent := globalEvents[index+1]
			switch nextEvent.GetType() {
			case EventTypeDiscard:
				tile = nextEvent.(*EventDiscard).Tile
			case EventTypeTsumoGiri:
				tile = nextEvent.(*EventTsumoGiri).Tile
			default:
				panic(errors.New("next event must be EventTypeDiscard or EventTypeTsumoGiri"))
			}
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallType == Riichi && call.CallTiles[0] == tile {
					posCall[who] = call
					break
				}
			}
			if _, ok := posCall[who]; !ok {
				panic(errors.New("can't find riichi call"))
			}
			index++
		case EventTypeTsumo:
			who := event.(*EventTsumo).Who
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallType == Tsumo {
					posCall[who] = call
					break
				}
			}
		case EventTypeChi:
			who := event.(*EventChi).Who
			for wind, calls := range posCalls {
				if wind == who {
					for _, call := range calls {
						if CallEqual(call, event.(*EventChi).Call) {
							posCall[who] = call
							break
						}
					}
				} else {
					posCall[wind] = SkipCall
				}
			}
		case EventTypePon:
			who := event.(*EventPon).Who
			for wind, calls := range posCalls {
				if wind == who {
					for _, call := range calls {
						if CallEqual(call, event.(*EventPon).Call) {
							posCall[who] = call
							break
						}
					}
				} else {
					posCall[wind] = SkipCall
				}
			}
		case EventTypeDaiMinKan:
			who := event.(*EventDaiMinKan).Who
			for wind, calls := range posCalls {
				if wind == who {
					for _, call := range calls {
						if CallEqual(call, event.(*EventDaiMinKan).Call) {
							posCall[who] = call
							break
						}
					}
				} else {
					posCall[wind] = SkipCall
				}
			}
		case EventTypeRyuuKyoku:
			reason := event.(*EventRyuuKyoku).Reason
			if reason == RyuuKyokuKyuuShuKyuuHai {
				who := event.(*EventRyuuKyoku).Who
				calls := posCalls[who]
				for _, call := range calls {
					if call.CallType == KyuuShuKyuuHai {
						posCall[who] = call
						break
					}
				}
			}
		case EventTypeRon:
			who := event.(*EventRon).Who
			for wind, calls := range posCalls {
				if wind == who {
					for _, call := range calls {
						if call.CallType != Ron {
							continue
						}
						if call.CallTiles[0] == event.(*EventRon).WinTile {
							posCall[who] = call
							break
						}
					}
				} else {
					var isRon = false
					for _, call := range calls {
						// multiple ron
						if call.CallType == Ron {
							isRon = true
							posCall[wind] = nil
						}
					}
					if !isRon {
						posCall[wind] = SkipCall
					}
				}
			}
		case EventTypeChanKan:
			who := event.(*EventChanKan).Who
			for wind, calls := range posCalls {
				if wind == who {
					for _, call := range calls {
						if call.CallType != ChanKan {
							continue
						}
						if call.CallTiles[0] == event.(*EventChanKan).WinTile {
							posCall[who] = call
							break
						}
					}
				} else {
					var isRon = false
					for _, call := range calls {
						// multiple chankan
						if call.CallType == ChanKan {
							isRon = true
							posCall[wind] = nil
						}
					}
					if !isRon {
						posCall[wind] = SkipCall
					}
				}
			}
		}
		success := true
		for wind, calls := range posCalls {
			if _, ok := posCall[wind]; !ok {
				success = false
				break
			}
			if posCall[wind] == nil {
				success = false
				break
			}
			if calls.Index(posCall[wind]) == -1 {
				panic("posCall not in posCalls")
			}
		}
		if !success {
			index++
			continue
		}
		posCalls, _ = game.Step(posCall)
		index++
	}

	return game
}
