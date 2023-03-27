package mahjong

type BoardState struct {
	RoundWind      int   `json:"round_wind"`
	NumHonba       int   `json:"num_honba"`
	NumRiichi      int   `json:"num_riichi"`
	DoraIndicators Tiles `json:"dora_indicators"`
	PlayerWind     int   `json:"player_wind"`
	Position       int   `json:"position"`
	HandTiles      Tiles `json:"hand_tiles"`
	//RealAction     Call        `json:"real_action"`
	ValidActions   Calls       `json:"valid_actions"`
	RealActionIdx  int         `json:"action_idx"`
	NumRemainTiles int         `json:"remain_tiles"`
	P0             PlayerState `json:"0"`
	P1             PlayerState `json:"1"`
	P2             PlayerState `json:"2"`
	P3             PlayerState `json:"3"`
}

type PlayerState struct {
	Points         int   `json:"points"`
	Melds          Calls `json:"melds"`
	DiscardTiles   Tiles `json:"discards"`
	TilesTsumoGiri []int `json:"tsumo_giri"`
	IsRiichi       bool  `json:"riichi"`
	PointsReward   int   `json:"p_reward"`
	FinalReward    int   `json:"r_reward"`
}

func NewBoardState() *BoardState {
	return &BoardState{
		RoundWind:      -1,
		NumHonba:       -1,
		NumRiichi:      -1,
		DoraIndicators: nil,
		PlayerWind:     -1,
		Position:       -1,
		HandTiles:      nil,
		//RealAction:     Call{},
		ValidActions:   nil,
		RealActionIdx:  -1,
		NumRemainTiles: -1,
		P0:             PlayerState{},
		P1:             PlayerState{},
		P2:             PlayerState{},
		P3:             PlayerState{},
	}
}

func BoardStateCopy(boardState *BoardState) *BoardState {
	return &BoardState{
		RoundWind:      boardState.RoundWind,
		NumHonba:       boardState.NumHonba,
		NumRiichi:      boardState.NumRiichi,
		DoraIndicators: boardState.DoraIndicators,
		PlayerWind:     boardState.PlayerWind,
		Position:       boardState.Position,
		HandTiles:      boardState.HandTiles[:],
		ValidActions:   boardState.ValidActions[:],
		RealActionIdx:  boardState.RealActionIdx,
		NumRemainTiles: boardState.NumRemainTiles,
		P0:             boardState.P0,
		P1:             boardState.P1,
		P2:             boardState.P2,
		P3:             boardState.P3,
	}
}
