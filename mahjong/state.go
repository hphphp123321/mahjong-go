package mahjong

import (
	"errors"
	"github.com/hphphp123321/mahjong-go/common"
	"sort"
)

type gameState interface {
	step() map[Wind]Calls
	next(posCalls map[Wind]*Call) error
	String() string
}

type InitState struct {
	g     *Game
	tiles Tiles
}

func (s *InitState) step() map[Wind]Calls {
	s.g.NewGameRound()

	initTiles := s.g.Tiles.Setup(s.tiles)

	// generate event
	var posEvent = make(map[Wind]Event)
	for wind, player := range s.g.posPlayer {
		player.HandTiles = append(player.HandTiles, initTiles[wind]...)
		sort.Sort(&player.HandTiles)
		player.Wind = wind
		posEvent[wind] = &EventStart{
			WindRound: s.g.WindRound,
			InitWind:  wind,
			Seed:      s.g.seed,
			NumGame:   s.g.NumGame,
			NumHonba:  s.g.NumHonba,
			NumRiichi: s.g.NumRiichi,
			InitTiles: initTiles[wind],
			Rule:      s.g.rule,
		}
	}
	s.g.addPosEvent(posEvent)
	return make(map[Wind]Calls)
}

func (s *InitState) next(posCalls map[Wind]*Call) error {
	if len(posCalls) != 0 {
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
		return make(map[Wind]Calls)
	}
	tile := s.g.Tiles.DealTile(s.dealRinshan)
	if s.dealRinshan {
		// generate new indicator event
		indicatorTileID := s.g.Tiles.GetCurrentIndicator()
		var posEvent = make(map[Wind]Event)
		for wind := range s.g.posPlayer {
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
	for wind := range s.g.posPlayer {
		var t = TileDummy
		if wind == pMain.Wind {
			t = tile
		}
		posEvent[wind] = &EventGet{
			Who:  pMain.Wind,
			Tile: t,
		}
	}
	s.g.addPosEvent(posEvent)

	if len(validActions) == 0 {
		panic("no valid action")
	}
	return map[Wind]Calls{
		pMain.Wind: validActions,
	}
}

func (s *DealState) next(posCalls map[Wind]*Call) error {
	if len(posCalls) == 0 {
		// normal ryuu kyoku
		if s.g.GetNumRemainTiles() != 0 {
			panic("remain tiles not 0")
		}
		s.g.State = &EndState{
			g: s.g,
			posResults: map[Wind]*Result{
				East:  {RyuuKyokuReason: RyuuKyokuNormal},
				South: {RyuuKyokuReason: RyuuKyokuNormal},
				West:  {RyuuKyokuReason: RyuuKyokuNormal},
				North: {RyuuKyokuReason: RyuuKyokuNormal},
			},
		}
		return nil
	}
	if len(posCalls) > 1 {
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
		s.g.breakIppatsu()
		// generate riichi event
		var posEvent = make(map[Wind]Event)
		for wind := range s.g.posPlayer {
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
	case KyuuShuKyuuHai:
		s.g.processKyuuShuKyuuHai()
		s.g.State = &EndState{
			g:          s.g,
			posResults: map[Wind]*Result{pMain.Wind: {RyuuKyokuReason: RyuuKyokuKyuuShuKyuuHai}},
		}
	default:
		return errors.New("unknown call type")
	}
	return nil
}

func (s *DealState) String() string {
	return "After Deal"
}

type DiscardState struct {
	g         *Game
	tileID    Tile
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
	for wind := range s.g.posPlayer {
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
		if _, ok := posCalls[wind]; !ok {
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

	pMain := s.g.posPlayer[s.g.Position]
	// if there is no call after discard, then next player deal
	if len(posCalls) == 0 {
		if s.g.judgeSuuChaRiichi() {
			// if suu cha riichi
			s.g.State = &EndState{
				g: s.g,
				posResults: map[Wind]*Result{
					East:  {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
					South: {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
					West:  {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
					North: {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
				},
			}
			return nil
		}

		if pMain.RiichiStep == 1 {
			var posEvent map[Wind]Event
			s.g.processRiichiStep2(pMain)
			// generate riichi step 2 event
			posEvent = make(map[Wind]Event)
			for wind := range s.g.posPlayer {
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
		return nil
	}

	// if there are calls, then process it
	var maxCallType = Skip
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
		// if max call is not ron
		if s.g.judgeSuuChaRiichi() {
			// if suu cha riichi
			s.g.State = &EndState{
				g: s.g,
				posResults: map[Wind]*Result{
					East:  {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
					South: {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
					West:  {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
					North: {RyuuKyokuReason: RyuuKyokuSuuChaRiichi},
				},
			}
			return nil
		}

		var wind Wind
		var call *Call
		for w, c := range posCalls {
			if c.CallType == maxCallType {
				wind = w
				call = c
				break
			}
		}
		if pMain.RiichiStep == 1 {
			var posEvent map[Wind]Event
			s.g.processRiichiStep2(pMain)
			// generate riichi step 2 event
			posEvent = make(map[Wind]Event)
			for w := range s.g.posPlayer {
				posEvent[w] = &EventRiichi{
					Who:  pMain.Wind,
					Step: 2,
				}
			}
			s.g.addPosEvent(posEvent)
		}
		// if max call is skip, then next player deal
		if maxCallType == Skip {
			s.g.Position = (s.g.Position + 1) % 4
			s.g.State = &DealState{
				g:           s.g,
				dealRinshan: false,
			}
			return nil
		}
		player := s.g.posPlayer[wind]
		s.g.Position = wind
		switch call.CallType {
		case Chi, Pon:
			if call.CallType == Chi {
				s.g.processChi(player, call)
				s.g.breakIppatsu()
				s.g.breakRyuukyoku()
			} else {
				s.g.processPon(player, call)
				s.g.breakIppatsu()
				s.g.breakRyuukyoku()
			}
			s.g.State = &ChiPonState{
				g:    s.g,
				call: call,
			}
		case DaiMinKan:
			s.g.processDaiMinKan(player, call)
			s.g.breakIppatsu()
			s.g.breakRyuukyoku()
			s.g.State = &KanState{
				g:    s.g,
				call: call,
			}
		default:
			return errors.New("unknown call type")
		}
	}

	return nil
}

func (s *DiscardState) String() string {
	return "After Discard"
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
	for wind := range s.g.posPlayer {
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
	for wind := range s.g.posPlayer {
		posEvent[wind] = kanEvent
	}
	s.g.addPosEvent(posEvent)
	return validCalls
}

func (s *KanState) next(posCalls map[Wind]*Call) error {
	// if there is no chan kan call
	if len(posCalls) == 0 {
		s.g.posPlayer[s.g.Position].KanNum++
		if s.g.judgeSuuKaiKan() {
			s.g.State = &EndState{
				g: s.g,
				posResults: map[Wind]*Result{
					East:  {RyuuKyokuReason: RyuuKyokuSuuKaiKan},
					South: {RyuuKyokuReason: RyuuKyokuSuuKaiKan},
					West:  {RyuuKyokuReason: RyuuKyokuSuuKaiKan},
					North: {RyuuKyokuReason: RyuuKyokuSuuKaiKan},
				},
			}
			return nil
		}
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
	return "After Kan"
}

type EndState struct {
	g          *Game
	posResults map[Wind]*Result
}

func (s *EndState) step() map[Wind]Calls {
	var prePoints = map[Wind]int{
		East:  s.g.posPlayer[East].Points,
		South: s.g.posPlayer[South].Points,
		West:  s.g.posPlayer[West].Points,
		North: s.g.posPlayer[North].Points,
	}
	var pointsChanges = make(map[Wind]int)
	var posEvent = make(map[Wind]Event)

	switch len(s.posResults) {
	case 4:
		// ryuu kyoku
		if s.posResults[East].RyuuKyokuReason == RyuuKyokuNormal {
			tenhaiWinds := s.g.judgeTenHaiWinds()
			if s.g.rule.IsNagashiMangan {
				// judge ryuu kyoku mangan
				bSlice := s.g.judgeNagashiMangan()
				if len(bSlice) != 0 {
					s.g.processNagashiMangan(bSlice)
					// generate nagashi mangan events
					for _, wind := range bSlice {
						for w := range s.g.posPlayer {
							posEvent[w] = &EventNagashiMangan{
								Who: wind,
							}
						}
						s.g.addPosEvent(posEvent)
						posEvent = make(map[Wind]Event)
					}

				} else {
					s.g.processNormalRyuuKyoku(tenhaiWinds)
				}
			} else {
				// process normal ryuu kyoku
				s.g.processNormalRyuuKyoku(tenhaiWinds)
			}
			for w := range s.g.posPlayer {
				for _, wind := range tenhaiWinds {
					player := s.g.posPlayer[wind]
					posEvent[w] = &EventTenhaiEnd{
						Who:         wind,
						HandTiles:   player.HandTiles,
						TenhaiSlice: player.TenhaiSlice,
					}
				}
				s.g.addPosEvent(posEvent)
				posEvent = make(map[Wind]Event)
			}

			if !common.Contain(tenhaiWinds, East) {
				// east not ten hai
				s.g.nextRound = true
				//s.g.WindRound++
			}
			s.g.honbaPlus = true
			//s.g.NumHonba++

		} else {
			// process special ryuu kyoku
			s.g.honbaPlus = true
			//s.g.NumHonba++
		}
		for wind, result := range s.posResults {
			posEvent[wind] = &EventRyuuKyoku{
				Who:    wind,
				Reason: result.RyuuKyokuReason,
			}
		}
		s.g.addPosEvent(posEvent)
		posEvent = make(map[Wind]Event)

	case 3:
		// san cha ron
		if !s.g.rule.IsSanChaHou {
			// san cha ron not allowed
			// generate ryuu kyoku events
			for wind := range s.g.posPlayer {
				posEvent[wind] = &EventRyuuKyoku{
					Who:    wind,
					Reason: RyuuKyokuSanChaHou,
				}
			}
			s.g.addPosEvent(posEvent)
			posEvent = make(map[Wind]Event)

			_, ok := s.posResults[East]
			if !ok {
				s.g.nextRound = true
				//s.g.WindRound++
			}
			s.g.honbaPlus = true
			//s.g.NumHonba++
		} else {
			s.g.processRonResult(s.posResults)
			// generate san cha ron events
			s.g.addRonEvents(s.posResults)

			_, ok := s.posResults[East]
			if ok {
				s.g.honbaPlus = true
				//s.g.NumHonba++
			} else {
				s.g.nextRound = true
				//s.g.WindRound++
			}
		}

	case 2:
		// double ron
		s.g.processRonResult(s.posResults)
		// generate double ron events
		s.g.addRonEvents(s.posResults)
	default:
		// normal ron, tsumo, chankan, kyuushukyuuhai
		var wind Wind
		var result *Result
		for w, r := range s.posResults {
			wind = w
			result = r
		}
		if result.RyuuKyokuReason == RyuuKyokuKyuuShuKyuuHai {
			// generate ryuu kyoku events
			for w := range s.g.posPlayer {
				posEvent[w] = &EventRyuuKyoku{
					Who:       wind,
					HandTiles: s.g.posPlayer[wind].HandTiles,
					Reason:    RyuuKyokuKyuuShuKyuuHai,
				}
			}
			s.g.addPosEvent(posEvent)
			posEvent = make(map[Wind]Event)

			if wind != East {
				s.g.nextRound = true
				//s.g.WindRound++
			}
			s.g.honbaPlus = true
			//s.g.NumHonba++
		} else {
			if result.RonCall.CallType != Tsumo {
				s.g.processRonResult(s.posResults)
				// generate ron events
				s.g.addRonEvents(s.posResults)
			} else {
				if result.RonCall.CallType != Tsumo {
					panic("result.RonCall.CallType != Tsumo")
				}
				s.g.processTsumoResult(wind, result)
				// generate tsumo events
				for w := range s.g.posPlayer {
					posEvent[w] = &EventTsumo{
						Who:       wind,
						HandTiles: s.g.posPlayer[wind].HandTiles,
						WinTile:   result.RonCall.CallTiles[0],
						Result:    result,
					}
				}
				s.g.addPosEvent(posEvent)
				posEvent = make(map[Wind]Event)

				if wind == East {
					s.g.honbaPlus = true
					//s.g.NumHonba++
				} else {
					s.g.nextRound = true
					//s.g.WindRound++
				}
			}
		}
	}

	// calculate scores changes
	for wind, points := range prePoints {
		pointsChanges[wind] = s.g.posPlayer[wind].Points - points
	}

	// generate end events
	for w := range s.g.posPlayer {
		posEvent[w] = &EventEnd{
			PointsChange: pointsChanges,
		}
	}
	s.g.addPosEvent(posEvent)

	if s.g.CheckGameEnd() {
		return make(map[Wind]Calls)
	}
	return map[Wind]Calls{
		East:  {NewCall(Next, nil, nil)},
		South: {NewCall(Next, nil, nil)},
		West:  {NewCall(Next, nil, nil)},
		North: {NewCall(Next, nil, nil)},
	}
}

func (s *EndState) next(posCalls map[Wind]*Call) error {
	if len(posCalls) == 0 {
		return ErrGameEnd
	}
	s.g.State = &InitState{
		g: s.g,
	}
	return nil
}

func (s *EndState) String() string {
	return "End"
}
