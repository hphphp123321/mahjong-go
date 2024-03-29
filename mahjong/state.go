package mahjong

import (
	"errors"
	"github.com/hphphp123321/go-common"
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
	s.g.newGameRound()

	initTiles := s.g.Tiles.Setup(s.tiles)

	// generate event
	var posEvent = make(map[Wind]Event)
	for wind, player := range s.g.PosPlayer {
		player.HandTiles = append(player.HandTiles, initTiles[wind]...)
		sort.Sort(&player.HandTiles)
		player.Wind = wind
		player.ShantenNum = player.GetShantenNum()
		posEvent[wind] = &EventStart{
			WindRound:         s.g.WindRound,
			InitWind:          wind,
			Seed:              s.g.Seed,
			NumGame:           s.g.NumGame,
			NumHonba:          s.g.NumHonba,
			NumRiichi:         s.g.NumRiichi,
			InitDoraIndicator: s.g.Tiles.GetCurrentIndicator(),
			InitTiles:         initTiles[wind],
			Rule:              s.g.Rule,
			PlayersPoints: func() map[Wind]int {
				m := make(map[Wind]int)
				for w, p := range s.g.PosPlayer {
					m[w] = p.Points
				}
				return m
			}(),
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
		for wind := range s.g.PosPlayer {
			posEvent[wind] = &EventNewIndicator{
				Tile: indicatorTileID,
			}
		}
		s.g.addPosEvent(posEvent)
	}

	pMain := s.g.PosPlayer[s.g.Position]
	s.g.getTileProcess(pMain, tile)
	validActions := s.g.JudgeSelfCalls(pMain)

	// generate event
	var posEvent = make(map[Wind]Event)
	for wind := range s.g.PosPlayer {
		var t = TileDummy
		var tenpaiInfos TenpaiInfos
		if wind == pMain.Wind {
			t = tile
			tenpaiInfos = GetTenpaiInfos(s.g, pMain)
		}
		posEvent[wind] = &EventGet{
			Who:         pMain.Wind,
			Tile:        t,
			TenpaiInfos: tenpaiInfos,
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
	pMain := s.g.PosPlayer[s.g.Position]
	switch call.CallType {
	case Discard:
		furitenBef := pMain.IsFuriten()
		s.g.discardTileProcess(pMain, call.CallTiles[0])
		tsumoGiri := pMain.TilesTsumoGiri[len(pMain.TilesTsumoGiri)-1]
		s.g.State = &DiscardState{
			g:                s.g,
			tileID:           call.CallTiles[0],
			tsumoGiri:        tsumoGiri,
			isFuritenChanged: furitenBef != pMain.IsFuriten(),
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
		furitenBef := pMain.IsFuriten()
		s.g.processRiichiStep1(pMain, call)
		s.g.breakIppatsu()
		// generate riichi event
		var posEvent = make(map[Wind]Event)
		for wind := range s.g.PosPlayer {
			posEvent[wind] = &EventRiichi{
				Who:  pMain.Wind,
				Step: 1,
			}
		}
		s.g.addPosEvent(posEvent)
		s.g.breakIppatsu()
		tsumoGiri := pMain.TilesTsumoGiri[len(pMain.TilesTsumoGiri)-1]
		s.g.State = &DiscardState{
			g:                s.g,
			tileID:           call.CallTiles[0],
			tsumoGiri:        tsumoGiri,
			isFuritenChanged: furitenBef != pMain.IsFuriten(),
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
	g                *Game
	tileID           Tile
	tsumoGiri        bool
	isFuritenChanged bool
}

// step after one player discard a tile
func (s *DiscardState) step() map[Wind]Calls {
	pMain := s.g.PosPlayer[s.g.Position]

	// judge other calls
	var validCalls = make(map[Wind]Calls)
	for wind, player := range s.g.PosPlayer {
		if wind == s.g.Position {
			continue
		}
		if calls := s.g.JudgeOtherCalls(player, s.tileID); len(calls) > 0 {
			validCalls[wind] = calls
		}
	}

	// generate discard event
	var event Event
	var posEvent = make(map[Wind]Event)
	var tenpaiInfo = GetTenpaiInfo(s.g, pMain)
	for wind := range s.g.PosPlayer {
		if wind == s.g.Position {
			if s.tsumoGiri {
				event = &EventTsumoGiri{
					Who:        s.g.Position,
					Tile:       s.tileID,
					TenpaiInfo: tenpaiInfo,
				}
			} else {
				event = &EventDiscard{
					Who:        s.g.Position,
					Tile:       s.tileID,
					TenpaiInfo: tenpaiInfo,
				}
			}
		} else {
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
		}
		posEvent[wind] = event
	}
	s.g.addPosEvent(posEvent)

	// self player furiten check
	posEvent = make(map[Wind]Event)
	if s.isFuritenChanged {
		if pMain.IsFuriten() {
			event = &EventFuriten{
				Who:           pMain.Wind,
				FuritenReason: FuritenDiscard,
			}
		} else {
			event = &EventFuriten{
				Who:           pMain.Wind,
				FuritenReason: FuritenNone,
			}
		}
		posEvent[s.g.Position] = event
		s.g.addPosEvent(posEvent)
	}
	return validCalls
}

func (s *DiscardState) next(posCalls map[Wind]*Call) error {
	// other players furiten check
	for _, wind := range s.g.getOtherWinds() {
		player := s.g.PosPlayer[wind]
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

	pMain := s.g.PosPlayer[s.g.Position]
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
			for wind := range s.g.PosPlayer {
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
			player := s.g.PosPlayer[wind]
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
			for w := range s.g.PosPlayer {
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
		player := s.g.PosPlayer[wind]
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
	validCalls[s.g.Position] = s.g.JudgeDiscardCall(s.g.PosPlayer[s.g.Position])
	pMain := s.g.PosPlayer[s.g.Position]

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
	for wind := range s.g.PosPlayer {
		if wind == s.g.Position {
			switch event.GetType() {
			case EventTypeChi:
				event = &EventChi{
					Who:         s.g.Position,
					Call:        s.call,
					TenpaiInfos: GetTenpaiInfos(s.g, pMain),
				}
			case EventTypePon:
				event = &EventPon{
					Who:         s.g.Position,
					Call:        s.call,
					TenpaiInfos: GetTenpaiInfos(s.g, pMain),
				}
			}
		}
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
	player := s.g.PosPlayer[s.g.Position]
	tileID := posCalls[s.g.Position].CallTiles[0]
	s.g.discardTileProcess(player, tileID)
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
	for wind, player := range s.g.PosPlayer {
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
	for wind := range s.g.PosPlayer {
		posEvent[wind] = kanEvent
	}
	s.g.addPosEvent(posEvent)
	return validCalls
}

func (s *KanState) next(posCalls map[Wind]*Call) error {
	// if there is no chan kan call
	if len(posCalls) == 0 {
		s.g.PosPlayer[s.g.Position].KanNum++
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
			player := s.g.PosPlayer[wind]
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
		East:  s.g.PosPlayer[East].Points,
		South: s.g.PosPlayer[South].Points,
		West:  s.g.PosPlayer[West].Points,
		North: s.g.PosPlayer[North].Points,
	}
	var pointsChanges = make(map[Wind]int)
	var posEvent = make(map[Wind]Event)

	switch len(s.posResults) {
	case 4:
		// ryuu kyoku
		if s.posResults[East].RyuuKyokuReason == RyuuKyokuNormal {
			TenpaiWinds := s.g.judgeTenpaiWinds()
			if s.g.Rule.IsNagashiMangan {
				// judge ryuu kyoku mangan
				bSlice := s.g.judgeNagashiMangan()
				if len(bSlice) != 0 {
					s.g.processNagashiMangan(bSlice)
					// generate nagashi mangan events
					for _, wind := range bSlice {
						for w := range s.g.PosPlayer {
							posEvent[w] = &EventNagashiMangan{
								Who: wind,
							}
						}
						s.g.addPosEvent(posEvent)
						posEvent = make(map[Wind]Event)
					}

				} else {
					s.g.processNormalRyuuKyoku(TenpaiWinds)
				}
			} else {
				// process normal ryuu kyoku
				s.g.processNormalRyuuKyoku(TenpaiWinds)
			}
			for w := range s.g.PosPlayer {
				for _, wind := range TenpaiWinds {
					player := s.g.PosPlayer[wind]
					posEvent[w] = &EventTenpaiEnd{
						Who:         wind,
						HandTiles:   player.HandTiles,
						TenpaiSlice: player.TenpaiSlice,
					}
				}
				s.g.addPosEvent(posEvent)
				posEvent = make(map[Wind]Event)
			}

			if !common.SliceContain(TenpaiWinds, East) {
				// east not ten hai
				s.g.nextRound = true
			}
			s.g.honbaPlus = true

		} else {
			// process special ryuu kyoku
			s.g.honbaPlus = true
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
		if !s.g.Rule.IsSanChaHou {
			// san cha ron not allowed
			// generate ryuu kyoku events
			for wind := range s.g.PosPlayer {
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
			}
			s.g.honbaPlus = true
		} else {
			s.g.processRonResult(s.posResults)
			// generate san cha ron events
			s.g.addRonEvents(s.posResults)

			_, ok := s.posResults[East]
			if ok {
				s.g.honbaPlus = true
			} else {
				s.g.nextRound = true
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
			for w := range s.g.PosPlayer {
				posEvent[w] = &EventRyuuKyoku{
					Who:       wind,
					HandTiles: s.g.PosPlayer[wind].HandTiles,
					Reason:    RyuuKyokuKyuuShuKyuuHai,
				}
			}
			s.g.addPosEvent(posEvent)
			posEvent = make(map[Wind]Event)

			if wind != East {
				s.g.nextRound = true
			}
			s.g.honbaPlus = true
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
				for w := range s.g.PosPlayer {
					posEvent[w] = &EventTsumo{
						Who:       wind,
						HandTiles: s.g.PosPlayer[wind].HandTiles,
						WinTile:   result.RonCall.CallTiles[0],
						Result:    result,
					}
				}
				s.g.addPosEvent(posEvent)
				posEvent = make(map[Wind]Event)
			}

			if wind == East {
				s.g.honbaPlus = true
			} else {
				s.g.nextRound = true
			}
		}
	}

	// calculate scores changes
	for wind, points := range prePoints {
		pointsChanges[wind] = s.g.PosPlayer[wind].Points - points
	}

	// generate end events
	for w := range s.g.PosPlayer {
		posEvent[w] = &EventEnd{
			PointsChange: pointsChanges,
		}
	}
	s.g.addPosEvent(posEvent)

	if s.g.CheckGameEnd() {
		return make(map[Wind]Calls)
	}
	return map[Wind]Calls{
		East:  {NextCall},
		South: {NextCall},
		West:  {NextCall},
		North: {NextCall},
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
