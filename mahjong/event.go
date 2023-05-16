package mahjong

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Events []Event

type Event interface {
	GetType() EventType
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

var eventTypes = map[string]reflect.Type{
	EventTypeGet.String():           reflect.TypeOf(EventGet{}),
	EventTypeTsumoGiri.String():     reflect.TypeOf(EventTsumoGiri{}),
	EventTypeRiichi.String():        reflect.TypeOf(EventRiichi{}),
	EventTypeDiscard.String():       reflect.TypeOf(EventDiscard{}),
	EventTypeChi.String():           reflect.TypeOf(EventChi{}),
	EventTypePon.String():           reflect.TypeOf(EventPon{}),
	EventTypeAnKan.String():         reflect.TypeOf(EventAnKan{}),
	EventTypeShouMinKan.String():    reflect.TypeOf(EventShouMinKan{}),
	EventTypeDaiMinKan.String():     reflect.TypeOf(EventDaiMinKan{}),
	EventTypeRon.String():           reflect.TypeOf(EventRon{}),
	EventTypeTsumo.String():         reflect.TypeOf(EventTsumo{}),
	EventTypeChanKan.String():       reflect.TypeOf(EventChanKan{}),
	EventTypeRiichi.String():        reflect.TypeOf(EventRiichi{}),
	EventTypeNewIndicator.String():  reflect.TypeOf(EventNewIndicator{}),
	EventTypeRyuuKyoku.String():     reflect.TypeOf(EventRyuuKyoku{}),
	EventTypeStart.String():         reflect.TypeOf(EventStart{}),
	EventTypeEnd.String():           reflect.TypeOf(EventEnd{}),
	EventTypeFuriten.String():       reflect.TypeOf(EventFuriten{}),
	EventTypeNagashiMangan.String(): reflect.TypeOf(EventNagashiMangan{}),
	EventTypeTenpaiEnd.String():     reflect.TypeOf(EventTenpaiEnd{}),
	EventTypeGlobalInit.String():    reflect.TypeOf(EventGlobalInit{}),
	// ... 添加其他事件类型
}

func (events *Events) MarshalJSON() ([]byte, error) {
	var tmp []*struct {
		Type  string `json:"type"`
		Event Event  `json:"event"`
	}
	for _, e := range *events {
		tmp = append(tmp, &struct {
			Type  string `json:"type"`
			Event Event  `json:"event"`
		}{
			Type:  e.GetType().String(),
			Event: e,
		})
	}
	return json.Marshal(tmp)
}

func (events *Events) UnmarshalJSON(data []byte) error {
	var tmp []*struct {
		Type  string          `json:"type"`
		Event json.RawMessage `json:"event"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	for _, e := range tmp {

		eventType, ok := eventTypes[e.Type]
		if !ok {
			return fmt.Errorf("unknown event type: %s", e.Type)
		}

		// create a new event
		eventValue := reflect.New(eventType).Elem()
		eventInterface := eventValue.Addr().Interface().(Event)

		// unmarshal the event
		err = eventInterface.UnmarshalJSON(e.Event)
		if err != nil {
			return err
		}

		*events = append(*events, eventInterface)
	}
	return err
}

type EventGet struct {
	Who         Wind        `json:"who"`
	Tile        Tile        `json:"tile,omitempty"`
	TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
}

func (event *EventGet) GetType() EventType {
	return EventTypeGet
}

func (event *EventGet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who         string      `json:"who"`
		Tile        string      `json:"tile,omitempty"`
		TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
	}{
		Who:         event.Who.String(),
		Tile:        event.Tile.String(),
		TenpaiInfos: event.TenpaiInfos,
	})
}

func (event *EventGet) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who         string      `json:"who"`
		Tile        string      `json:"tile,omitempty"`
		TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = MapStringToTile[tmp.Tile]
	event.TenpaiInfos = tmp.TenpaiInfos
	return nil
}

type EventTsumoGiri struct {
	Who        Wind        `json:"who"`
	Tile       Tile        `json:"tile"`
	TenpaiInfo *TenpaiInfo `json:"tenpai_info,omitempty"`
}

func (event *EventTsumoGiri) GetType() EventType {
	return EventTypeTsumoGiri
}

func (event *EventTsumoGiri) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who        string      `json:"who"`
		Tile       string      `json:"tile"`
		TenpaiInfo *TenpaiInfo `json:"tenpai_info,omitempty"`
	}{
		Who:        event.Who.String(),
		Tile:       event.Tile.String(),
		TenpaiInfo: event.TenpaiInfo,
	})
}

func (event *EventTsumoGiri) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who        string      `json:"who"`
		Tile       string      `json:"tile"`
		TenpaiInfo *TenpaiInfo `json:"tenpai_info,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = MapStringToTile[tmp.Tile]
	event.TenpaiInfo = tmp.TenpaiInfo
	return nil
}

type EventDiscard struct {
	Who        Wind        `json:"who"`
	Tile       Tile        `json:"tile"`
	TenpaiInfo *TenpaiInfo `json:"tenpai_info,omitempty"`
}

func (event *EventDiscard) GetType() EventType {
	return EventTypeDiscard
}

func (event *EventDiscard) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who        string      `json:"who"`
		Tile       string      `json:"tile"`
		TenpaiInfo *TenpaiInfo `json:"tenpai_info,omitempty"`
	}{
		Who:        event.Who.String(),
		Tile:       event.Tile.String(),
		TenpaiInfo: event.TenpaiInfo,
	})
}

func (event *EventDiscard) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who        string      `json:"who"`
		Tile       string      `json:"tile"`
		TenpaiInfo *TenpaiInfo `json:"tenpai_info,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = MapStringToTile[tmp.Tile]
	event.TenpaiInfo = tmp.TenpaiInfo
	return nil
}

type EventChi struct {
	Who         Wind        `json:"who"`
	Call        *Call       `json:"call"`
	TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
}

func (event *EventChi) GetType() EventType {
	return EventTypeChi
}

func (event *EventChi) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who         string      `json:"who"`
		Call        *Call       `json:"call"`
		TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
	}{
		Who:         event.Who.String(),
		Call:        event.Call,
		TenpaiInfos: event.TenpaiInfos,
	})
}

func (event *EventChi) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who         string      `json:"who"`
		Call        *Call       `json:"call"`
		TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Call = tmp.Call
	event.TenpaiInfos = tmp.TenpaiInfos
	return nil
}

type EventPon struct {
	Who         Wind        `json:"who"`
	Call        *Call       `json:"call"`
	TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
}

func (event *EventPon) GetType() EventType {
	return EventTypePon
}

func (event *EventPon) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who         string      `json:"who"`
		Call        *Call       `json:"call"`
		TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
	}{
		Who:         event.Who.String(),
		Call:        event.Call,
		TenpaiInfos: event.TenpaiInfos,
	})
}

func (event *EventPon) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who         string      `json:"who"`
		Call        *Call       `json:"call"`
		TenpaiInfos TenpaiInfos `json:"tenpai_infos,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Call = tmp.Call
	event.TenpaiInfos = tmp.TenpaiInfos
	return nil
}

type EventDaiMinKan struct {
	Who  Wind  `json:"who"`
	Call *Call `json:"call"`
}

func (event *EventDaiMinKan) GetType() EventType {
	return EventTypeDaiMinKan
}

func (event *EventDaiMinKan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}{
		Who:  event.Who.String(),
		Call: event.Call,
	})
}

func (event *EventDaiMinKan) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Call = tmp.Call
	return nil
}

type EventShouMinKan struct {
	Who  Wind  `json:"who"`
	Call *Call `json:"call"`
}

func (event *EventShouMinKan) GetType() EventType {
	return EventTypeShouMinKan
}

func (event *EventShouMinKan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}{
		Who:  event.Who.String(),
		Call: event.Call,
	})
}

func (event *EventShouMinKan) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Call = tmp.Call
	return nil
}

type EventAnKan struct {
	Who  Wind  `json:"who"`
	Call *Call `json:"call"`
}

func (event *EventAnKan) GetType() EventType {
	return EventTypeAnKan
}

func (event *EventAnKan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}{
		Who:  event.Who.String(),
		Call: event.Call,
	})
}

func (event *EventAnKan) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Call = tmp.Call
	return nil
}

type EventRiichi struct {
	Who  Wind `json:"who"`
	Step int  `json:"step"`
}

func (event *EventRiichi) GetType() EventType {
	return EventTypeRiichi
}

func (event *EventRiichi) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Step int    `json:"step"`
	}{
		Who:  event.Who.String(),
		Step: event.Step,
	})
}

func (event *EventRiichi) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Step int    `json:"step"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Step = tmp.Step
	return nil
}

type EventRon struct {
	Who       Wind    `json:"who"`
	FromWho   Wind    `json:"from_who"`
	HandTiles Tiles   `json:"hand_tiles"`
	WinTile   Tile    `json:"win_tile"`
	Result    *Result `json:"result"`
}

func (event *EventRon) GetType() EventType {
	return EventTypeRon
}

func (event *EventRon) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who       string  `json:"who"`
		FromWho   string  `json:"from_who"`
		HandTiles Tiles   `json:"hand_tiles"`
		WinTile   string  `json:"win_tile"`
		Result    *Result `json:"result"`
	}{
		Who:       event.Who.String(),
		FromWho:   event.FromWho.String(),
		HandTiles: event.HandTiles,
		WinTile:   event.WinTile.String(),
		Result:    event.Result,
	})
}

func (event *EventRon) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who       string `json:"who"`
		FromWho   string `json:"from_who"`
		HandTiles Tiles  `json:"hand_tiles"`
		WinTile   string `json:"win_tile"`
		Result    *Result
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.FromWho = MapStringToWind[tmp.FromWho]
	event.HandTiles = tmp.HandTiles
	event.WinTile = MapStringToTile[tmp.WinTile]
	event.Result = tmp.Result
	return nil
}

type EventTsumo struct {
	Who       Wind    `json:"who"`
	HandTiles Tiles   `json:"hand_tiles"`
	WinTile   Tile    `json:"win_tile"`
	Result    *Result `json:"result"`
}

func (event *EventTsumo) GetType() EventType {
	return EventTypeTsumo
}

func (event *EventTsumo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who       string `json:"who"`
		HandTiles Tiles  `json:"hand_tiles"`
		WinTile   string `json:"win_tile"`
		Result    *Result
	}{
		Who:       event.Who.String(),
		HandTiles: event.HandTiles,
		WinTile:   event.WinTile.String(),
		Result:    event.Result,
	})
}

func (event *EventTsumo) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who       string  `json:"who"`
		HandTiles Tiles   `json:"hand_tiles"`
		WinTile   string  `json:"win_tile"`
		Result    *Result `json:"result"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.HandTiles = tmp.HandTiles
	event.WinTile = MapStringToTile[tmp.WinTile]
	event.Result = tmp.Result
	return nil
}

type EventNewIndicator struct {
	Tile Tile `json:"tile"`
}

func (event *EventNewIndicator) GetType() EventType {
	return EventTypeNewIndicator
}

func (event *EventNewIndicator) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Tile string `json:"tile"`
	}{
		Tile: event.Tile.String(),
	})
}

func (event *EventNewIndicator) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Tile string `json:"tile"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Tile = MapStringToTile[tmp.Tile]
	return nil
}

type EventChanKan struct {
	Who       Wind    `json:"who"`
	FromWho   Wind    `json:"from_who"`
	HandTiles Tiles   `json:"hand_tiles"`
	WinTile   Tile    `json:"win_tile"`
	Result    *Result `json:"result"`
}

func (event *EventChanKan) GetType() EventType {
	return EventTypeChanKan
}

func (event *EventChanKan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who       string  `json:"who"`
		FromWho   string  `json:"from_who"`
		HandTiles Tiles   `json:"hand_tiles"`
		WinTile   string  `json:"win_tile"`
		Result    *Result `json:"result"`
	}{
		Who:       event.Who.String(),
		FromWho:   event.FromWho.String(),
		HandTiles: event.HandTiles,
		WinTile:   event.WinTile.String(),
		Result:    event.Result,
	})
}

func (event *EventChanKan) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who       string  `json:"who"`
		FromWho   string  `json:"from_who"`
		HandTiles Tiles   `json:"hand_tiles"`
		WinTile   string  `json:"win_tile"`
		Result    *Result `json:"result"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.FromWho = MapStringToWind[tmp.FromWho]
	event.WinTile = MapStringToTile[tmp.WinTile]
	event.HandTiles = tmp.HandTiles
	event.Result = tmp.Result
	return nil
}

type EventRyuuKyoku struct {
	Who       Wind            `json:"who,omitempty"`
	HandTiles Tiles           `json:"hand_tiles,omitempty"`
	Reason    RyuuKyokuReason `json:"reason,int"`
}

func (event *EventRyuuKyoku) GetType() EventType {
	return EventTypeRyuuKyoku
}

func (event *EventRyuuKyoku) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who       string `json:"who,omitempty"`
		HandTiles Tiles  `json:"hand_tiles,omitempty"`
		Reason    string `json:"reason"`
	}{
		Who:       event.Who.String(),
		HandTiles: event.HandTiles,
		Reason:    event.Reason.String(),
	})
}

func (event *EventRyuuKyoku) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who       string `json:"who,omitempty"`
		HandTiles Tiles  `json:"hand_tiles,omitempty"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	who := MapStringToWind[tmp.Who]
	reason := MapStringToRyuuKyokuReason[tmp.Reason]
	event.Who = who
	event.Reason = reason
	event.HandTiles = tmp.HandTiles
	return nil
}

type EventGlobalInit struct {
	AllTiles   Tiles        `json:"all_tiles"`
	WindRound  WindRound    `json:"wind_round"`
	Seed       int64        `json:"seed"`
	NumGame    int          `json:"num_game"`
	NumHonba   int          `json:"num_honba"`
	NumRiichi  int          `json:"num_riichi"`
	Rule       *Rule        `json:"rule"`
	InitPoints map[Wind]int `json:"init_points"`
}

func (event *EventGlobalInit) GetType() EventType {
	return EventTypeGlobalInit
}

func (event *EventGlobalInit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		AllTiles   Tiles          `json:"all_tiles"`
		WindRound  string         `json:"wind_round"`
		Seed       int64          `json:"seed"`
		NumGame    int            `json:"num_game"`
		NumHonba   int            `json:"num_honba"`
		NumRiichi  int            `json:"num_riichi"`
		Rule       *Rule          `json:"rule"`
		InitPoints map[string]int `json:"init_points"`
	}{
		AllTiles:  event.AllTiles,
		WindRound: event.WindRound.String(),
		Seed:      event.Seed,
		NumGame:   event.NumGame,
		NumHonba:  event.NumHonba,
		NumRiichi: event.NumRiichi,
		Rule:      event.Rule,
		InitPoints: func() map[string]int {
			m := make(map[string]int)
			for k, v := range event.InitPoints {
				m[k.String()] = v
			}
			return m
		}(),
	})
}

func (event *EventGlobalInit) UnmarshalJSON(data []byte) error {
	var tmp struct {
		AllTiles   Tiles          `json:"all_tiles"`
		WindRound  string         `json:"wind_round"`
		Seed       int64          `json:"seed"`
		NumGame    int            `json:"num_game"`
		NumHonba   int            `json:"num_honba"`
		NumRiichi  int            `json:"num_riichi"`
		Rule       *Rule          `json:"rule"`
		InitPoints map[string]int `json:"init_points"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.AllTiles = tmp.AllTiles
	event.WindRound = MapStringToWindRound[tmp.WindRound]
	event.Seed = tmp.Seed
	event.NumGame = tmp.NumGame
	event.NumHonba = tmp.NumHonba
	event.NumRiichi = tmp.NumRiichi
	event.Rule = tmp.Rule
	event.InitPoints = func() map[Wind]int {
		m := make(map[Wind]int)
		for k, v := range tmp.InitPoints {
			m[MapStringToWind[k]] = v
		}
		return m
	}()
	return nil
}

type EventStart struct {
	WindRound WindRound `json:"wind_round"`
	InitWind  Wind      `json:"init_wind"`

	Seed      int64 `json:"seed"`
	NumGame   int   `json:"num_game"`
	NumHonba  int   `json:"num_honba"`
	NumRiichi int   `json:"num_riichi"`

	InitDoraIndicator Tile  `json:"init_dora_indicator"`
	InitTiles         Tiles `json:"init_tiles"`

	PlayersPoints map[Wind]int `json:"players_points"`

	Rule *Rule `json:"rule"`
}

func (event *EventStart) GetType() EventType {
	return EventTypeStart
}

func (event *EventStart) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		WindRound         string `json:"wind_round"`
		InitWind          string `json:"init_wind"`
		Seed              int64  `json:"seed"`
		NumGame           int    `json:"num_game"`
		NumHonba          int    `json:"num_honba"`
		NumRiichi         int    `json:"num_riichi"`
		InitDoraIndicator string `json:"init_dora_indicator"`
		InitTiles         Tiles  `json:"init_tiles"`
		Rule              *Rule  `json:"rule"`
		PlayersPoints     map[string]int
	}{
		WindRound:         event.WindRound.String(),
		InitWind:          event.InitWind.String(),
		Seed:              event.Seed,
		NumGame:           event.NumGame,
		NumHonba:          event.NumHonba,
		NumRiichi:         event.NumRiichi,
		InitDoraIndicator: event.InitDoraIndicator.String(),
		InitTiles:         event.InitTiles,
		Rule:              event.Rule,
		PlayersPoints: func() map[string]int {
			m := make(map[string]int)
			for k, v := range event.PlayersPoints {
				m[k.String()] = v
			}
			return m
		}(),
	})
}

func (event *EventStart) UnmarshalJSON(data []byte) error {
	var tmp struct {
		WindRound         string         `json:"wind_round"`
		InitWind          string         `json:"init_wind"`
		Seed              int64          `json:"seed"`
		NumGame           int            `json:"num_game"`
		NumHonba          int            `json:"num_honba"`
		NumRiichi         int            `json:"num_riichi"`
		InitDoraIndicator string         `json:"init_dora_indicator"`
		InitTiles         Tiles          `json:"init_tiles"`
		Rule              *Rule          `json:"rule"`
		PlayersPoints     map[string]int `json:"players_points"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.WindRound = MapStringToWindRound[tmp.WindRound]
	event.InitWind = MapStringToWind[tmp.InitWind]
	event.Seed = tmp.Seed
	event.NumGame = tmp.NumGame
	event.NumHonba = tmp.NumHonba
	event.NumRiichi = tmp.NumRiichi
	event.InitDoraIndicator = MapStringToTile[tmp.InitDoraIndicator]
	event.InitTiles = tmp.InitTiles
	event.Rule = tmp.Rule
	event.PlayersPoints = func() map[Wind]int {
		m := make(map[Wind]int)
		for k, v := range tmp.PlayersPoints {
			m[MapStringToWind[k]] = v
		}
		return m
	}()
	return nil
}

type EventEnd struct {
	PointsChange map[Wind]int `json:"points_change"`
}

func (event *EventEnd) GetType() EventType {
	return EventTypeEnd
}

func (event *EventEnd) MarshalJSON() ([]byte, error) {
	pointsChange := make(map[string]int)
	for k, v := range event.PointsChange {
		pointsChange[k.String()] = v
	}
	return json.Marshal(&struct {
		PointsChange map[string]int `json:"points_change"`
	}{
		PointsChange: pointsChange,
	})
}

func (event *EventEnd) UnmarshalJSON(data []byte) error {
	var tmp struct {
		PointsChange map[string]int `json:"points_change"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	pointsChange := make(map[Wind]int)
	for k, v := range tmp.PointsChange {
		pointsChange[MapStringToWind[k]] = v
	}
	event.PointsChange = pointsChange
	return nil
}

type EventFuriten struct {
	Who           Wind          `json:"who"`
	FuritenReason FuritenReason `json:"furiten_reason"`
}

func (event *EventFuriten) GetType() EventType {
	return EventTypeFuriten
}

func (event *EventFuriten) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who           string `json:"who"`
		FuritenReason string `json:"furiten_reason"`
	}{
		Who:           event.Who.String(),
		FuritenReason: event.FuritenReason.String(),
	})
}

func (event *EventFuriten) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who           string `json:"who"`
		FuritenReason string `json:"furiten_reason"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.FuritenReason = MapStringToFuritenReason[tmp.FuritenReason]
	return nil
}

type EventNagashiMangan struct {
	Who Wind `json:"who"`
}

func (event *EventNagashiMangan) GetType() EventType {
	return EventTypeNagashiMangan
}

func (event *EventNagashiMangan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who string `json:"who"`
	}{
		Who: event.Who.String(),
	})
}

func (event *EventNagashiMangan) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who string `json:"who"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	return nil
}

type EventTenpaiEnd struct {
	Who         Wind        `json:"who"`
	HandTiles   Tiles       `json:"hand_tiles"`
	TenpaiSlice TileClasses `json:"Tenpai_slice"`
}

func (event *EventTenpaiEnd) GetType() EventType {
	return EventTypeTenpaiEnd
}

func (event *EventTenpaiEnd) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who         string      `json:"who"`
		HandTiles   Tiles       `json:"hand_tiles"`
		TenpaiSlice TileClasses `json:"Tenpai_slice"`
	}{
		Who:         event.Who.String(),
		HandTiles:   event.HandTiles,
		TenpaiSlice: event.TenpaiSlice,
	})
}

func (event *EventTenpaiEnd) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who         string      `json:"who"`
		HandTiles   Tiles       `json:"hand_tiles"`
		TenpaiSlice TileClasses `json:"Tenpai_slice"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.HandTiles = tmp.HandTiles
	event.TenpaiSlice = tmp.TenpaiSlice
	return nil
}
