package mahjong

import "errors"

type gameState interface {
	step() map[Wind]Calls
	next(posCalls map[Wind]Call) error
	String() string
}

type InitState struct {
	g *Game
}

func (s *InitState) step() map[Wind]Calls {
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
	return nil
}

func (s *InitState) next(posCalls map[Wind]Call) error {
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

func (s *DealState) next(posCalls map[Wind]Call) error {
	if len(posCalls) != 1 {
		return errors.New("invalid call nums")
	}
	call := posCalls[s.g.Position]
	pMain := s.g.posPlayer[s.g.Position]
	switch call.CallType {
	case Discard:
		s.g.DiscardTileProcess(pMain, call.CallTiles[0])
		s.g.State = &DiscardState{
			g: s.g,
		}
	case ShouMinKan:
		s.g.processShouMinKan(pMain, &call)
		s.g.breakIppatsu()
		s.g.breakRyuukyoku()
		s.g.State = &KanState{
			g: s.g,
		}
	case AnKan:
		s.g.processAnKan(pMain, &call)
		s.g.breakIppatsu()
		s.g.breakRyuukyoku()
		s.g.State = &KanState{
			g: s.g,
		}
	case Riichi:
		s.g.processRiichi(pMain, &call)
		s.g.breakIppatsu()
		s.g.State = &RiichiState{
			g: s.g,
		}
	case Tsumo:
		result := s.g.processTsumo(pMain, &call)
		s.g.State = &EndState{
			g:      s.g,
			result: result,
		}
	case KyuShuKyuHai:
		s.g.processKyuShuKyuHai(pMain, &call)
	default:
		return errors.New("unknown call type")
	}
	return nil
}

func (s *DealState) String() string {
	return "Deal"
}

type DiscardState struct {
	g *Game
}

func (s *DiscardState) step() map[Wind]Calls {
	return nil
}

func (s *DiscardState) next(posCalls map[Wind]Call) error {
	return nil
}

func (s *DiscardState) String() string {
	return "Discard"
}

type KanState struct {
	g *Game
}

func (s *KanState) next(posCalls map[Wind]Call) error {
	return nil
}

func (s *KanState) step() map[Wind]Calls {
	return nil
}

func (s *KanState) String() string {
	return "Call"
}

type RiichiState struct {
	g *Game
}

func (s *RiichiState) step() map[Wind]Calls {
	return nil
}

func (s *RiichiState) next(posCalls map[Wind]Call) error {
	return nil
}

func (s *RiichiState) String() string {
	return "Riichi"
}

type EndState struct {
	g      *Game
	result *Result
}

func (s *EndState) step() map[Wind]Calls {
	return nil
}

func (s *EndState) next(posCalls map[Wind]Call) error {
	return nil
}

func (s *EndState) String() string {
	return "End"
}
