package mahjong

type gameState interface {
	next(state gameState)
	status() string
}

type InitState struct {
	g *Game
}

func (s *InitState) next(state gameState) {
	switch state.(type) {
	case *DealState:
		s.g.State = state
	default:
		panic("Invalid state")
	}
}

func (s *InitState) status() string {
	return "Init"
}

type DealState struct {
	g *Game
}

func (s *DealState) next(state gameState) {
	s.g.ProcessSelfCall(s.g.posPlayer[s.g.Position], s.g.posCall[s.g.Position])

	s.g.State = state
}

func (s *DealState) status() string {
	return "Deal"
}

type DiscardState struct {
	g *Game
}

func (s *DiscardState) next(state gameState) {
	s.g.State = state
}

func (s *DiscardState) status() string {
	return "Discard"
}

type KanState struct {
	g *Game
}

func (s *KanState) next(state gameState) {
	s.g.State = state
}

func (s *KanState) status() string {
	return "Call"
}

type RiichiState struct {
	g *Game
}

func (s *RiichiState) next(state gameState) {
	s.g.State = state
}

func (s *RiichiState) status() string {
	return "Riichi"
}
