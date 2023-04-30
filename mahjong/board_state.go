package mahjong

import (
	"encoding/json"
	"github.com/hphphp123321/go-common"
	"sort"
)

type BoardState struct {
	WindRound      WindRound `json:"wind_round"`
	NumHonba       int       `json:"num_honba"`
	NumRiichi      int       `json:"num_riichi"`
	DoraIndicators Tiles     `json:"dora_indicators"`
	PlayerWind     Wind      `json:"player_wind"`
	Position       Wind      `json:"position"`
	HandTiles      Tiles     `json:"hand_tiles"`
	ValidActions   Calls     `json:"valid_actions,omitempty"`
	//RealActionIdx  int         `json:"action_idx"`
	NumRemainTiles int                   `json:"remain_tiles"`
	PlayerStates   map[Wind]*PlayerState `json:"player_states"`
}

type PlayerState struct {
	Points         int    `json:"points"`
	Melds          Calls  `json:"melds"`
	DiscardTiles   Tiles  `json:"discards"`
	TilesTsumoGiri []bool `json:"tsumo_giri"`
	IsRiichi       bool   `json:"riichi"`
}

func NewBoardState() *BoardState {
	return &BoardState{
		WindRound:      -1,
		NumHonba:       -1,
		NumRiichi:      -1,
		DoraIndicators: make(Tiles, 0, 5),
		PlayerWind:     -1,
		Position:       -1,
		HandTiles:      make(Tiles, 0, 14),
		ValidActions:   nil,
		//RealActionIdx:  -1,
		NumRemainTiles: -1,
		PlayerStates: map[Wind]*PlayerState{
			East: &PlayerState{
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
			South: &PlayerState{
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
			West: &PlayerState{
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
			North: &PlayerState{
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
		},
	}
}

func BoardStateCopy(boardState *BoardState) *BoardState {
	return &BoardState{
		WindRound:      boardState.WindRound,
		NumHonba:       boardState.NumHonba,
		NumRiichi:      boardState.NumRiichi,
		DoraIndicators: boardState.DoraIndicators.Copy(),
		PlayerWind:     boardState.PlayerWind,
		Position:       boardState.Position,
		HandTiles:      boardState.HandTiles.Copy(),
		ValidActions:   boardState.ValidActions.Copy(),
		//RealActionIdx:  boardState.RealActionIdx,
		NumRemainTiles: boardState.NumRemainTiles,
		PlayerStates:   boardState.PlayerStates,
	}
}

func (b *BoardState) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			WindRound      string `json:"wind_round"`
			NumHonba       int    `json:"num_honba"`
			NumRiichi      int    `json:"num_riichi"`
			DoraIndicators Tiles  `json:"dora_indicators"`
			PlayerWind     string `json:"player_wind"`
			Position       string `json:"position"`
			HandTiles      Tiles  `json:"hand_tiles"`
			ValidActions   Calls  `json:"valid_actions,omitempty"`
			//RealActionIdx  int         `json:"action_idx"`
			NumRemainTiles int                   `json:"remain_tiles"`
			PlayerStates   map[Wind]*PlayerState `json:"player_states"`
		}{
			WindRound:      b.WindRound.String(),
			NumHonba:       b.NumHonba,
			NumRiichi:      b.NumRiichi,
			DoraIndicators: b.DoraIndicators,
			PlayerWind:     b.PlayerWind.String(),
			Position:       b.Position.String(),
			HandTiles:      b.HandTiles,
			ValidActions:   b.ValidActions,
			//RealActionIdx:  b.RealActionIdx,
			NumRemainTiles: b.NumRemainTiles,
			PlayerStates:   b.PlayerStates,
		},
	)
}

func (b *BoardState) UnmarshalJSON(data []byte) error {
	var tmp struct {
		WindRound      string `json:"wind_round"`
		NumHonba       int    `json:"num_honba"`
		NumRiichi      int    `json:"num_riichi"`
		DoraIndicators Tiles  `json:"dora_indicators"`
		PlayerWind     string `json:"player_wind"`
		Position       string `json:"position"`
		HandTiles      Tiles  `json:"hand_tiles"`
		ValidActions   Calls  `json:"valid_actions,omitempty"`
		//RealActionIdx  int         `json:"action_idx"`
		NumRemainTiles int                   `json:"remain_tiles"`
		PlayerStates   map[Wind]*PlayerState `json:"player_states"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	b.WindRound = MapStringToWindRound[tmp.WindRound]
	b.NumHonba = tmp.NumHonba
	b.NumRiichi = tmp.NumRiichi
	b.DoraIndicators = tmp.DoraIndicators
	b.PlayerWind = MapStringToWind[tmp.PlayerWind]
	b.Position = MapStringToWind[tmp.Position]
	b.HandTiles = tmp.HandTiles
	b.ValidActions = tmp.ValidActions
	b.NumRemainTiles = tmp.NumRemainTiles
	b.PlayerStates = tmp.PlayerStates
	return nil
}

func (b *BoardState) DecodeEvents(events Events) {
	for _, e := range events {
		switch e.GetType() {
		case EventTypeStart:
			b.handleEventStart(e)
		case EventTypeGet:
			b.handleEventGet(e)
		case EventTypeDiscard:
			b.handleEventDiscard(e)
		case EventTypeTsumoGiri:
			b.handleEventTsumoGiri(e)
		case EventTypeChi:
			b.handleEventChi(e)
		case EventTypePon:
			b.handleEventPon(e)
		case EventTypeDaiMinKan:
			b.handleEventDaiMinKan(e)
		case EventTypeShouMinKan:
			b.handleEventShouMinKan(e)
		case EventTypeAnKan:
			b.handleEventAnKan(e)
		case EventTypeRiichi:
			b.handleEventRiichi(e)
		case EventTypeNewIndicator:
			b.handleEventNewIndicator(e)
		}
	}
}

func (b *BoardState) Equal(bs *BoardState) bool {
	if b.WindRound != bs.WindRound {
		return false
	}
	if b.NumHonba != bs.NumHonba {
		return false
	}
	if b.NumRiichi != bs.NumRiichi {
		return false
	}
	if !common.SliceEqual(b.DoraIndicators, bs.DoraIndicators) {
		return false
	}
	if b.PlayerWind != bs.PlayerWind {
		return false
	}
	if b.Position != bs.Position {
		return false
	}
	if !common.SliceEqual(b.HandTiles, bs.HandTiles) {
		return false
	}
	if b.NumRemainTiles != bs.NumRemainTiles {
		return false
	}
	for wind, ps := range b.PlayerStates {
		if ps.Points != bs.PlayerStates[wind].Points {
			return false
		}
		if ps.IsRiichi != bs.PlayerStates[wind].IsRiichi {
			return false
		}
		if !common.SliceEqual(ps.DiscardTiles, bs.PlayerStates[wind].DiscardTiles) {
			return false
		}
		if !common.SliceEqual(ps.TilesTsumoGiri, bs.PlayerStates[wind].TilesTsumoGiri) {
			return false
		}
		for i, meld := range ps.Melds {
			if !CallEqual(meld, bs.PlayerStates[wind].Melds[i]) {
				return false
			}
		}
	}
	return true
}

func (b *BoardState) handleEventStart(event Event) {
	b.WindRound = event.(*EventStart).WindRound
	b.NumHonba = event.(*EventStart).NumHonba
	b.NumRiichi = event.(*EventStart).NumRiichi
	b.DoraIndicators.Append(event.(*EventStart).InitDoraIndicator)
	b.Position = East
	b.PlayerWind = event.(*EventStart).InitWind
	for _, tile := range event.(*EventStart).InitTiles {
		b.HandTiles.Append(tile)
	}
	sort.Sort(&b.HandTiles)
	for wind, points := range event.(*EventStart).PlayersPoints {
		b.PlayerStates[wind].Points = points
	}
	b.NumRemainTiles = 70
}

func (b *BoardState) handleEventGet(event Event) {
	if event.(*EventGet).Who == b.PlayerWind {
		b.HandTiles = append(b.HandTiles, event.(*EventGet).Tile)
	}
	b.NumRemainTiles--
	b.Position = event.(*EventGet).Who
}

func (b *BoardState) handleEventDiscard(event Event) {
	who := event.(*EventDiscard).Who
	if who == b.PlayerWind {
		b.HandTiles.Remove(event.(*EventDiscard).Tile)
		sort.Sort(&b.HandTiles)
	}
	b.PlayerStates[who].DiscardTiles.Append(event.(*EventDiscard).Tile)
	b.PlayerStates[who].TilesTsumoGiri = append(b.PlayerStates[who].TilesTsumoGiri, false)

}

func (b *BoardState) handleEventTsumoGiri(event Event) {
	who := event.(*EventTsumoGiri).Who
	if who == b.PlayerWind {
		b.HandTiles.Remove(event.(*EventTsumoGiri).Tile)
	}
	b.PlayerStates[who].DiscardTiles.Append(event.(*EventTsumoGiri).Tile)
	b.PlayerStates[who].TilesTsumoGiri = append(b.PlayerStates[who].TilesTsumoGiri, true)
}

func (b *BoardState) handleEventChi(event Event) {
	who := event.(*EventChi).Who
	call := event.(*EventChi).Call
	if who == b.PlayerWind {
		for i := 0; i < 2; i++ {
			b.HandTiles.Remove(call.CallTiles[i])
		}
	}
	b.PlayerStates[who].Melds.Append(call)
	b.Position = who
}

func (b *BoardState) handleEventPon(event Event) {
	who := event.(*EventPon).Who
	call := event.(*EventPon).Call

	if who == b.PlayerWind {
		for i := 0; i < 2; i++ {
			b.HandTiles.Remove(call.CallTiles[i])
		}
	}
	b.PlayerStates[who].Melds.Append(call)
	b.Position = who
}

func (b *BoardState) handleEventDaiMinKan(event Event) {
	who := event.(*EventDaiMinKan).Who
	call := event.(*EventDaiMinKan).Call

	b.PlayerStates[who].Melds.Append(call)
	if who == b.PlayerWind {
		for i := 0; i < 3; i++ {
			b.HandTiles.Remove(call.CallTiles[i])
		}
	}
	b.Position = who
}

func (b *BoardState) handleEventShouMinKan(event Event) {
	who := event.(*EventShouMinKan).Who
	call := event.(*EventShouMinKan).Call
	tc := call.CallTiles[0].Class()
	for i, meld := range b.PlayerStates[who].Melds {
		if meld.CallType == Pon && meld.CallTiles[0].Class() == tc {
			b.PlayerStates[who].Melds[i] = call.Copy()
			break
		}
	}
	if who == b.PlayerWind {
		b.HandTiles.Remove(call.CallTiles[3])
	}
	b.Position = who
}

func (b *BoardState) handleEventAnKan(event Event) {
	who := event.(*EventAnKan).Who
	call := event.(*EventAnKan).Call
	b.PlayerStates[who].Melds.Append(call)
	if who == b.PlayerWind {
		for i := 0; i < 4; i++ {
			b.HandTiles.Remove(call.CallTiles[i])
		}
	}
	b.Position = who
}

func (b *BoardState) handleEventRiichi(event Event) {
	who := event.(*EventRiichi).Who
	step := event.(*EventRiichi).Step
	if step == 2 {
		b.PlayerStates[who].IsRiichi = true
	}
}

func (b *BoardState) handleEventNewIndicator(event Event) {
	b.DoraIndicators = append(b.DoraIndicators, event.(*EventNewIndicator).Tile)
}
