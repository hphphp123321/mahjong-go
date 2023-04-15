package mahjong

import (
	"errors"
)

func ReConstructGame(playerSlice []*Player, globalEvents Events) *Game {
	game := &Game{}

	e := globalEvents[0]
	if e.GetType() != EventTypeGlobalInit {
		panic(errors.New("first event must be EventTypeGlobalInit"))
	}
	et := e.(*EventGlobalInit)
	game.Tiles = NewMahjongTiles(nil)
	game.Reset(playerSlice, et.AllTiles)
	game.WindRound = et.WindRound
	game.seed = et.Seed
	game.NumGame = et.NumGame
	game.NumHonba = et.NumHonba
	game.NumRiichi = et.NumRiichi
	game.rule = et.Rule
	game.posPlayer[Wind((16-game.WindRound)%4)].Points = et.InitPoints[Wind((16-game.WindRound)%4)]
	game.posPlayer[Wind((17-game.WindRound)%4)].Points = et.InitPoints[Wind((17-game.WindRound)%4)]
	game.posPlayer[Wind((18-game.WindRound)%4)].Points = et.InitPoints[Wind((18-game.WindRound)%4)]
	game.posPlayer[Wind((19-game.WindRound)%4)].Points = et.InitPoints[Wind((19-game.WindRound)%4)]

	index := 1
	var success = true
	var posCalls = make(map[Wind]Calls)
	for index < len(globalEvents) {
		if success {
			posCalls = game.Step()
		}
		var posCall = make(map[Wind]*Call)
		if len(posCalls) == 4 {
			return game
		} else if len(posCalls) == 0 {
			game.Next(posCall)
			index++
			continue
		}

		event := globalEvents[index]
		switch event.GetType() {
		case EventTypeGet:
			who := event.(*EventGet).Who
			if len(game.posPlayer[who].HandTiles) >= 14 {
				panic(errors.New("player hand tiles more than 14"))
			}
			for wind := range posCalls {
				if game.Position != wind {
					posCall[wind] = SkipCall
				}
			}
		case EventTypeDiscard:
			who := event.(*EventDiscard).Who
			if len(game.posPlayer[who].HandTiles) > 14 {
				panic(errors.New("player hand tiles more than 14"))
			}
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallTiles[0] == event.(*EventDiscard).Tile {
					posCall[who] = call
					break
				}
			}
		case EventTypeTsumoGiri:
			who := event.(*EventTsumoGiri).Who
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallTiles[0] == event.(*EventTsumoGiri).Tile {
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
			if step == 2 {
				continue
			}
			tile := globalEvents[index+1].(*EventDiscard).Tile
			calls := posCalls[who]
			for _, call := range calls {
				if call.CallType == Riichi && call.CallTiles[0] == tile {
					posCall[who] = call
					break
				}
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
		success = true
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
		game.Next(posCall)
		index++
	}

	return game
}
