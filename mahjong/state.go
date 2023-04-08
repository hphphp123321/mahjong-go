package mahjong

import "errors"

type gameState interface {
	step() map[Wind]Calls
	next(posCalls map[Wind]*Call) error
	String() string
}

type InitState struct {
	g *Game
}

func (s *InitState) step() map[Wind]Calls {
	var validCalls = make(map[Wind]Calls)
	for wind, _ := range s.g.posPlayer {
		validCalls[wind] = Calls{
			NewCall(Next, nil, nil),
		}
	}
	s.g.WindRound += 1
	s.g.NewGameRound()

	initTiles := s.g.Tiles.Setup()

	// generate event
	var posEvent = make(map[Wind]Event)
	for wind, player := range s.g.posPlayer {
		player.HandTiles = initTiles[wind]
		player.Wind = wind
		posEvent[wind] = &EventStart{
			WindRound: s.g.WindRound,
			InitWind:  wind,
			Seed:      s.g.seed,
			NumGame:   s.g.NumGame,
			NumHonba:  s.g.NumHonba,
			NumRiichi: s.g.NumRiichi,
			InitTiles: initTiles[wind],
		}
	}
	s.g.addPosEvent(posEvent)
	return validCalls
}

func (s *InitState) next(posCalls map[Wind]*Call) error {
	for _, call := range posCalls {
		if call.CallType != Next {
			return errors.New("invalid call")
		}
	}
	if len(posCalls) != 4 {
		return errors.New("invalid call nums")
	}
	s.g.State = &DealState{
		s.g,
		false,
	}
	return nil
}

func (s *InitState) String() string {
	return "Init"
}

type DealState struct {
	g           *Game
	dealRinshan bool
}

func (s *DealState) step() map[Wind]Calls {
	if s.g.GetNumRemainTiles() == 0 {
		return nil
	}
	tile := s.g.Tiles.DealTile(s.dealRinshan)
	if s.dealRinshan {
		// generate new indicator event
		indicatorTileID := s.g.Tiles.GetCurrentIndicator()
		var posEvent = make(map[Wind]Event)
		for wind, _ := range s.g.posPlayer {
			posEvent[wind] = &EventNewIndicator{
				Tile: indicatorTileID,
			}
		}
		s.g.addPosEvent(posEvent)
	}

	pMain := s.g.posPlayer[s.g.Position]
	s.g.GetTileProcess(pMain, tile)
	validActions := s.g.JudgeSelfCalls(pMain)

	// generate event
	var posEvent = make(map[Wind]Event)
	for wind, _ := range s.g.posPlayer {
		var t = -1
		if wind == pMain.Wind {
			t = tile
		}
		posEvent[wind] = &EventGet{
			Who:  pMain.Wind,
			Tile: t,
		}
	}
	s.g.addPosEvent(posEvent)

	return map[Wind]Calls{
		pMain.Wind: validActions,
	}
}

func (s *DealState) next(posCalls map[Wind]*Call) error {
	if len(posCalls) != 1 {
		return errors.New("invalid call nums")
	}
	call := posCalls[s.g.Position]
	pMain := s.g.posPlayer[s.g.Position]
	switch call.CallType {
	case Discard:
		s.g.DiscardTileProcess(pMain, call.CallTiles[0])
		tsumoGiri := pMain.TilesTsumoGiri[len(pMain.TilesTsumoGiri)-1]
		s.g.State = &DiscardState{
			g:         s.g,
			tileID:    call.CallTiles[0],
			tsumoGiri: tsumoGiri,
		}
	case ShouMinKan:
		s.g.processShouMinKan(pMain, call)
		s.g.breakIppatsu()
		s.g.breakRyuukyoku()
		s.g.State = &KanState{
			g:    s.g,
			call: call,
		}
	case AnKan:
		s.g.processAnKan(pMain, call)
		s.g.breakIppatsu()
		s.g.breakRyuukyoku()
		s.g.State = &KanState{
			g:    s.g,
			call: call,
		}
	case Riichi:
		s.g.processRiichiStep1(pMain, call)
		// generate riichi event
		var posEvent = make(map[Wind]Event)
		for wind, _ := range s.g.posPlayer {
			posEvent[wind] = &EventRiichi{
				Who:  pMain.Wind,
				Step: 1,
			}
		}
		s.g.addPosEvent(posEvent)
		s.g.breakIppatsu()
		tsumoGiri := pMain.TilesTsumoGiri[len(pMain.TilesTsumoGiri)-1]
		s.g.State = &DiscardState{
			g:         s.g,
			tileID:    call.CallTiles[0],
			tsumoGiri: tsumoGiri,
		}
	case Tsumo:
		result := s.g.processTsumo(pMain, call)
		s.g.State = &EndState{
			g:          s.g,
			posResults: map[Wind]*Result{pMain.Wind: result},
		}
	case KyuShuKyuHai:
		s.g.processKyuShuKyuHai(pMain, call)
	default:
		return errors.New("unknown call type")
	}
	return nil
}

func (s *DealState) String() string {
	return "Deal"
}

type DiscardState struct {
	g         *Game
	tileID    int
	tsumoGiri bool
}

// step after one player discard a tile
func (s *DiscardState) step() map[Wind]Calls {
	var validCalls = make(map[Wind]Calls)
	for wind, player := range s.g.posPlayer {
		if wind == s.g.Position {
			continue
		}
		if calls := s.g.JudgeOtherCalls(player, s.tileID); len(calls) > 0 {
			validCalls[wind] = calls
		}
	}
	// generate discard event
	var posEvent = make(map[Wind]Event)
	var event Event
	if s.tsumoGiri {
		event = &EventTsumoGiri{
			Who:  s.g.Position,
			Tile: s.tileID,
		}
	} else {
		event = &EventDiscard{
			Who:  s.g.Position,
			Tile: s.tileID,
		}
	}
	for wind, _ := range s.g.posPlayer {
		posEvent[wind] = event
	}
	s.g.addPosEvent(posEvent)

	return validCalls
}

func (s *DiscardState) next(posCalls map[Wind]*Call) error {
	// furiten check
	for _, wind := range s.g.getOtherWinds() {
		player := s.g.posPlayer[wind]
		if !player.FuritenStatus {
			continue
		}
		if posCalls[wind].CallType == Ron {
			continue
		}
		var furitenReason FuritenReason
		if player.IsRiichi {
			player.RiichiFuriten = true
			furitenReason = FuritenRiichi
		} else {
			player.JunFuriten = true
			furitenReason = FuritenJun
		}
		furitenEvent := &EventFuriten{
			Who:           wind,
			FuritenReason: furitenReason,
		}
		s.g.addPosEvent(map[Wind]Event{
			wind: furitenEvent,
		})
	}

	// if there is no call after discard, then next player deal
	if posCalls == nil {

		pMain := s.g.posPlayer[s.g.Position]
		if pMain.RiichiStep == 1 {
			var posEvent map[Wind]Event
			s.g.processRiichiStep2(pMain)
			// generate riichi step 2 event
			posEvent = make(map[Wind]Event)
			for wind, _ := range s.g.posPlayer {
				posEvent[wind] = &EventRiichi{
					Who:  pMain.Wind,
					Step: 2,
				}
			}
			s.g.addPosEvent(posEvent)
		}
		s.g.Position = (s.g.Position + 1) % 4
		s.g.State = &DealState{
			s.g,
			false,
		}
	}

	// if there are calls, then process it
	var maxCallType CallType = -1
	for _, call := range posCalls {
		if call.CallType > maxCallType {
			maxCallType = call.CallType
		}
	}
	// if max call is ron, process all ron
	if maxCallType == Ron {
		var results = make(map[Wind]*Result)
		for wind, call := range posCalls {
			if call.CallType != Ron {
				continue
			}
			player := s.g.posPlayer[wind]
			result := s.g.processRon(player, call)
			results[wind] = result
		}
		s.g.State = &EndState{
			g:          s.g,
			posResults: results,
		}
		return nil
	} else {
		for wind, call := range posCalls {
			player := s.g.posPlayer[wind]
			if call.CallType != maxCallType {
				continue
			}
			switch call.CallType {
			case Chi, Pon:
				s.g.processChi(player, call)
				s.g.State = &ChiPonState{
					g:    s.g,
					call: call,
				}
			case DaiMinKan:
				s.g.processDaiMinKan(player, call)
				s.g.State = &KanState{
					g:    s.g,
					call: call,
				}
			case KyuShuKyuHai:
				result := s.g.processKyuShuKyuHai(player, call)
				s.g.State = &EndState{
					g:          s.g,
					posResults: map[Wind]*Result{wind: result},
				}
			default:
				return errors.New("unknown call type")
			}
		}
	}
	return nil
}

func (s *DiscardState) String() string {
	return "Discard"
}

type ChiPonState struct {
	g    *Game
	call *Call
}

func (s *ChiPonState) String() string {
	return "After Chi Pon"
}

// step after one player chi or pon
func (s *ChiPonState) step() map[Wind]Calls {
	var validCalls = make(map[Wind]Calls)
	validCalls[s.g.Position] = s.g.JudgeDiscardCall(s.g.posPlayer[s.g.Position])

	// generate event
	var posEvent = make(map[Wind]Event)
	var event Event
	switch s.call.CallType {
	case Chi:
		event = &EventChi{
			Who:  s.g.Position,
			Call: s.call,
		}
	case Pon:
		event = &EventPon{
			Who:  s.g.Position,
			Call: s.call,
		}
	}
	for wind, _ := range s.g.posPlayer {
		posEvent[wind] = event
	}
	s.g.addPosEvent(posEvent)
	return validCalls
}

func (s *ChiPonState) next(posCalls map[Wind]*Call) error {
	// after chi pon you only can discard
	if len(posCalls) != 1 || posCalls[s.g.Position].CallType != Discard {
		return errors.New("invalid call after chi pon")
	}
	player := s.g.posPlayer[s.g.Position]
	tileID := posCalls[s.g.Position].CallTiles[0]
	s.g.DiscardTileProcess(player, tileID)
	tsumoGiri := player.TilesTsumoGiri[len(player.TilesTsumoGiri)-1]
	s.g.State = &DiscardState{
		g:         s.g,
		tileID:    tileID,
		tsumoGiri: tsumoGiri,
	}
	return nil
}

type KanState struct {
	g    *Game
	call *Call
}

// after one player kan
func (s *KanState) step() map[Wind]Calls {
	// judge chan kan
	var validCalls = make(map[Wind]Calls)
	for wind, player := range s.g.posPlayer {
		if wind == s.g.Position {
			continue
		}
		if calls := s.g.judgeChanKan(player, s.call.CallTiles[3], s.call.CallType == AnKan); len(calls) > 0 {
			validCalls[wind] = calls
		}
	}

	// generate event
	var posEvent = make(map[Wind]Event)
	var kanEvent Event
	switch s.call.CallType {
	case AnKan:
		kanEvent = &EventAnKan{
			Who:  s.g.Position,
			Call: s.call,
		}
	case ShouMinKan:
		kanEvent = &EventShouMinKan{
			Who:  s.g.Position,
			Call: s.call,
		}
	case DaiMinKan:
		kanEvent = &EventDaiMinKan{
			Who:  s.g.Position,
			Call: s.call,
		}

	}
	for wind, _ := range s.g.posPlayer {
		posEvent[wind] = kanEvent
	}
	s.g.addPosEvent(posEvent)
	return validCalls
}

func (s *KanState) next(posCalls map[Wind]*Call) error {
	// if there is no chan kan call
	if len(posCalls) == 0 {
		s.g.posPlayer[s.g.Position].KanNum++
		s.g.State = &DealState{
			g:           s.g,
			dealRinshan: true,
		}
	} else {
		var posResults = make(map[Wind]*Result)
		// if there are chan kan call, process it
		for wind, call := range posCalls {
			player := s.g.posPlayer[wind]
			posResults[wind] = s.g.processChanKan(player, call)
		}
		s.g.State = &EndState{
			g:          s.g,
			posResults: posResults,
		}
	}
	return nil
}

func (s *KanState) String() string {
	return "Call"
}

type RiichiState struct {
	g          *Game
	riichiStep int
	tileID     int
}

// after one player claim riichi
func (s *RiichiState) step() map[Wind]Calls {

	return nil
}

func (s *RiichiState) next(posCalls map[Wind]*Call) error {
	return nil
}

func (s *RiichiState) String() string {
	return "Riichi"
}

type EndState struct {
	g          *Game
	posResults map[Wind]*Result
}

func (s *EndState) step() map[Wind]Calls {
	return nil
}

func (s *EndState) next(posCalls map[Wind]*Call) error {
	return nil
}

func (s *EndState) String() string {
	return "End"
}
