package mahjong

import (
	"encoding/json"
	"fmt"
	"github.com/hphphp123321/go-common"
	"sort"
	"strings"
)

type BoardState struct {
	WindRound      WindRound             `json:"wind_round"`
	NumHonba       int                   `json:"num_honba"`
	NumRiichi      int                   `json:"num_riichi"`
	DoraIndicators Tiles                 `json:"dora_indicators"`
	PlayerWind     Wind                  `json:"player_wind"`
	Position       Wind                  `json:"position"`
	HandTiles      Tiles                 `json:"hand_tiles"`
	ValidActions   Calls                 `json:"valid_actions,omitempty"`
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
		NumRemainTiles: -1,
		PlayerStates: map[Wind]*PlayerState{
			East: {
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
			South: {
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
			West: {
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
			North: {
				Points:         25000,
				Melds:          make(Calls, 0, 4),
				DiscardTiles:   make(Tiles, 0, 25),
				TilesTsumoGiri: make([]bool, 0, 25),
				IsRiichi:       false,
			},
		},
	}
}

func (b *BoardState) Reset() {
	b.WindRound = -1
	b.NumHonba = -1
	b.NumRiichi = -1
	b.DoraIndicators = make(Tiles, 0, 5)
	b.PlayerWind = -1
	b.Position = -1
	b.HandTiles = make(Tiles, 0, 14)
	b.ValidActions = nil
	b.NumRemainTiles = -1
	b.PlayerStates = map[Wind]*PlayerState{
		East: {
			Points:         25000,
			Melds:          make(Calls, 0, 4),
			DiscardTiles:   make(Tiles, 0, 25),
			TilesTsumoGiri: make([]bool, 0, 25),
			IsRiichi:       false,
		},
		South: {
			Points:         25000,
			Melds:          make(Calls, 0, 4),
			DiscardTiles:   make(Tiles, 0, 25),
			TilesTsumoGiri: make([]bool, 0, 25),
			IsRiichi:       false,
		},
		West: {
			Points:         25000,
			Melds:          make(Calls, 0, 4),
			DiscardTiles:   make(Tiles, 0, 25),
			TilesTsumoGiri: make([]bool, 0, 25),
			IsRiichi:       false,
		},
		North: {
			Points:         25000,
			Melds:          make(Calls, 0, 4),
			DiscardTiles:   make(Tiles, 0, 25),
			TilesTsumoGiri: make([]bool, 0, 25),
			IsRiichi:       false,
		},
	}
}

func (b *BoardState) UTF8() string {
	var s strings.Builder

	// 顶部边框
	s.WriteString(" ______________________________________________________________________________________________________________________________________________________________________________\n")

	// 顶部玩家信息
	topPlayer := b.PlayerStates[(b.Position+2)%4]
	s.WriteString(fmt.Sprintf("|                                                     Discards: %-72s \n", topPlayer.DiscardTiles.UTF8()))
	s.WriteString(fmt.Sprintf("|                                                     Melds: %-72s \n", topPlayer.Melds.UTF8()))
	s.WriteString(fmt.Sprintf("|                                                     Points: %-72d \n", topPlayer.Points))
	s.WriteString(fmt.Sprintf("|                                                     Riichi: %-72t \n", topPlayer.IsRiichi))

	s.WriteString("|______________________________________________________________________________________________________________________________________________________________________________\n")

	// 左侧玩家信息、中间信息和右侧玩家信息
	leftPlayer := b.PlayerStates[(b.Position+3)%4]
	rightPlayer := b.PlayerStates[(b.Position+1)%4]
	leftPlayerDiscards := divideIntoLines(leftPlayer.DiscardTiles.UTF8(), 4)
	rightPlayerDiscards := divideIntoLines(rightPlayer.DiscardTiles.UTF8(), 4)

	middleLines := []string{
		fmt.Sprintf("Wind Round: %s  Num Honba: %d  Num Riichi: %d", b.WindRound, b.NumHonba, b.NumRiichi),
		fmt.Sprintf("Player Wind: %s  Position: %s", b.PlayerWind, b.Position),
		fmt.Sprintf("Dora Indicators: %s", b.DoraIndicators.UTF8()),
		fmt.Sprintf("Num Remaining Tiles: %d", b.NumRemainTiles),
	}

	for i := 0; i < 4; i++ {
		s.WriteString(fmt.Sprintf("| L Discards: %-30s | %-81s | R Discards: %-23s \n", leftPlayerDiscards[i], middleLines[i], rightPlayerDiscards[i]))
	}

	s.WriteString(fmt.Sprintf("| L Melds: %-30s | %-81s | R Melds: %-27s \n", leftPlayer.Melds.UTF8(), " ", rightPlayer.Melds.UTF8()))
	s.WriteString(fmt.Sprintf("| L Points: %-30d | %-81s | R Points: %-26d \n", leftPlayer.Points, " ", rightPlayer.Points))
	s.WriteString(fmt.Sprintf("| L Riichi: %-32t | %-81s | R Riichi: %-26t \n", leftPlayer.IsRiichi, " ", rightPlayer.IsRiichi))

	s.WriteString("|______________________________________________________________________________________________________________________________________________________________________________\n")
	// 底部玩家信息
	bottomPlayer := b.PlayerStates[b.Position]
	s.WriteString(fmt.Sprintf("|                                                     Discards: %-72s \n", bottomPlayer.DiscardTiles.UTF8()))
	s.WriteString(fmt.Sprintf("|                                                     Melds: %-75s \n", bottomPlayer.Melds.UTF8()))
	s.WriteString(fmt.Sprintf("|                                                     Points: %-73d \n", bottomPlayer.Points))
	s.WriteString(fmt.Sprintf("|                                                     Riichi: %-73t \n", bottomPlayer.IsRiichi))

	// 手牌信息
	s.WriteString(fmt.Sprintf("|                                                     Hand Tiles: %-68s \n", b.HandTiles.UTF8()))

	// 底部边框
	s.WriteString(" ______________________________________________________________________________________________________________________________________________________________________________\n")

	return s.String()
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
	b.Reset()
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
