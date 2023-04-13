package mahjong

import "encoding/json"

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
	NumRemainTiles int         `json:"remain_tiles"`
	PlayerEast     PlayerState `json:"player_east"`
	PlayerSouth    PlayerState `json:"player_south"`
	PlayerWest     PlayerState `json:"player_west"`
	PlayerNorth    PlayerState `json:"player_north"`
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
		DoraIndicators: nil,
		PlayerWind:     -1,
		Position:       -1,
		HandTiles:      nil,
		ValidActions:   nil,
		//RealActionIdx:  -1,
		NumRemainTiles: -1,
		PlayerEast:     PlayerState{},
		PlayerSouth:    PlayerState{},
		PlayerWest:     PlayerState{},
		PlayerNorth:    PlayerState{},
	}
}

func BoardStateCopy(boardState *BoardState) *BoardState {
	return &BoardState{
		WindRound:      boardState.WindRound,
		NumHonba:       boardState.NumHonba,
		NumRiichi:      boardState.NumRiichi,
		DoraIndicators: boardState.DoraIndicators,
		PlayerWind:     boardState.PlayerWind,
		Position:       boardState.Position,
		HandTiles:      boardState.HandTiles[:],
		ValidActions:   boardState.ValidActions[:],
		//RealActionIdx:  boardState.RealActionIdx,
		NumRemainTiles: boardState.NumRemainTiles,
		PlayerEast:     boardState.PlayerEast,
		PlayerSouth:    boardState.PlayerSouth,
		PlayerWest:     boardState.PlayerWest,
		PlayerNorth:    boardState.PlayerNorth,
	}
}

func (b *BoardState) MarshalJson() ([]byte, error) {
	return json.Marshal(
		struct {
			WindRound      string `json:"wind_round"`
			NumHonba       int    `json:"num_honba"`
			NumRiichi      int    `json:"num_riichi"`
			DoraIndicators Tiles  `json:"dora_indicators"`
			PlayerWind     Wind   `json:"player_wind"`
			Position       Wind   `json:"position"`
			HandTiles      Tiles  `json:"hand_tiles"`
			ValidActions   Calls  `json:"valid_actions,omitempty"`
			//RealActionIdx  int         `json:"action_idx"`
			NumRemainTiles int         `json:"remain_tiles"`
			PlayerEast     PlayerState `json:"player_east"`
			PlayerSouth    PlayerState `json:"player_south"`
			PlayerWest     PlayerState `json:"player_west"`
			PlayerNorth    PlayerState `json:"player_north"`
		}{
			WindRound:      b.WindRound.String(),
			NumHonba:       b.NumHonba,
			NumRiichi:      b.NumRiichi,
			DoraIndicators: b.DoraIndicators,
			PlayerWind:     b.PlayerWind,
			Position:       b.Position,
			HandTiles:      b.HandTiles,
			ValidActions:   b.ValidActions,
			//RealActionIdx:  b.RealActionIdx,
			NumRemainTiles: b.NumRemainTiles,
			PlayerEast:     b.PlayerEast,
			PlayerSouth:    b.PlayerSouth,
			PlayerWest:     b.PlayerWest,
			PlayerNorth:    b.PlayerNorth,
		},
	)
}

func (b *BoardState) UnmarshalJson(data []byte) error {
	var tmp struct {
		WindRound      string `json:"wind_round"`
		NumHonba       int    `json:"num_honba"`
		NumRiichi      int    `json:"num_riichi"`
		DoraIndicators Tiles  `json:"dora_indicators"`
		PlayerWind     Wind   `json:"player_wind"`
		Position       Wind   `json:"position"`
		HandTiles      Tiles  `json:"hand_tiles"`
		ValidActions   Calls  `json:"valid_actions,omitempty"`
		//RealActionIdx  int         `json:"action_idx"`
		NumRemainTiles int         `json:"remain_tiles"`
		PlayerEast     PlayerState `json:"player_east"`
		PlayerSouth    PlayerState `json:"player_south"`
		PlayerWest     PlayerState `json:"player_west"`
		PlayerNorth    PlayerState `json:"player_north"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	b.WindRound = MapStringToWindRound[tmp.WindRound]
	b.NumHonba = tmp.NumHonba
	b.NumRiichi = tmp.NumRiichi
	b.DoraIndicators = tmp.DoraIndicators
	b.PlayerWind = tmp.PlayerWind
	b.Position = tmp.Position
	b.HandTiles = tmp.HandTiles
	b.ValidActions = tmp.ValidActions
	//b.RealActionIdx = tmp.RealActionIdx
	b.NumRemainTiles = tmp.NumRemainTiles
	b.PlayerEast = tmp.PlayerEast
	b.PlayerSouth = tmp.PlayerSouth
	b.PlayerWest = tmp.PlayerWest
	b.PlayerNorth = tmp.PlayerNorth
	return nil
}
