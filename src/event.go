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
	EventTypeGet.String():          reflect.TypeOf(EventGet{}),
	EventTypeTsumoGiri.String():    reflect.TypeOf(EventTsumoGiri{}),
	EventTypeRiichi.String():       reflect.TypeOf(EventRiichi{}),
	EventTypeDiscard.String():      reflect.TypeOf(EventDiscard{}),
	EventTypeChi.String():          reflect.TypeOf(EventChi{}),
	EventTypePon.String():          reflect.TypeOf(EventPon{}),
	EventTypeAnKan.String():        reflect.TypeOf(EventAnKan{}),
	EventTypeShouMinKan.String():   reflect.TypeOf(EventShouMinKan{}),
	EventTypeDaiMinKan.String():    reflect.TypeOf(EventDaiMinKan{}),
	EventTypeRon.String():          reflect.TypeOf(EventRon{}),
	EventTypeTsumo.String():        reflect.TypeOf(EventTsumo{}),
	EventTypeChanKan.String():      reflect.TypeOf(EventChanKan{}),
	EventTypeRiichi.String():       reflect.TypeOf(EventRiichi{}),
	EventTypeNewIndicator.String(): reflect.TypeOf(EventNewIndicator{}),
	EventTypeRyuuKyoku.String():    reflect.TypeOf(EventRyuuKyoku{}),
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

		// 创建事件实例
		eventValue := reflect.New(eventType).Elem()
		eventInterface := eventValue.Addr().Interface().(Event)

		// 反序列化事件
		err = eventInterface.UnmarshalJSON(e.Event)
		if err != nil {
			return err
		}

		*events = append(*events, eventInterface)

		//switch e.Type {
		//case EventTypeGet.String():
		//	var event EventGet
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeTsumoGiri.String():
		//	var event EventTsumoGiri
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeDiscard.String():
		//	var event EventDiscard
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeChi.String():
		//	var event EventChi
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypePon.String():
		//	var event EventPon
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeAnKan.String():
		//	var event EventAnKan
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeShouMinKan.String():
		//	var event EventShouMinKan
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeDaiMinKan.String():
		//	var event EventDaiMinKan
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeRon.String():
		//	var event EventRon
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeTsumo.String():
		//	var event EventTsumo
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeChanKan.String():
		//	var event EventChanKan
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeRiichi.String():
		//	var event EventRiichi
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeNewIndicator.String():
		//	var event EventNewIndicator
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//case EventTypeRyuuKyoku.String():
		//	var event EventRyuuKyoku
		//	err = event.UnmarshalJSON(e.Event)
		//	*events = append(*events, &event)
		//}
	}
	return err
}

type EventGet struct {
	Who  Wind `json:"who"`
	Tile int  `json:"tile,omitempty"`
}

func (event *EventGet) GetType() EventType {
	return EventTypeGet
}

func (event *EventGet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Tile int    `json:"tile,omitempty"`
	}{
		Who:  event.Who.String(),
		Tile: event.Tile,
	})
}

func (event *EventGet) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Tile int    `json:"tile,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = tmp.Tile
	return nil
}

type EventTsumoGiri struct {
	Who  Wind `json:"who"`
	Tile int  `json:"tile: int,omitempty"`
}

func (event *EventTsumoGiri) GetType() EventType {
	return EventTypeTsumoGiri
}

func (event *EventTsumoGiri) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Tile int    `json:"tile,omitempty"`
	}{
		Who:  event.Who.String(),
		Tile: event.Tile,
	})
}

func (event *EventTsumoGiri) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Tile int    `json:"tile,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = tmp.Tile
	return nil
}

type EventDiscard struct {
	Who  Wind `json:"who"`
	Tile int  `json:"tile: int,omitempty"`
}

func (event *EventDiscard) GetType() EventType {
	return EventTypeDiscard
}

func (event *EventDiscard) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Tile int    `json:"tile,omitempty"`
	}{
		Who:  event.Who.String(),
		Tile: event.Tile,
	})
}

func (event *EventDiscard) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Tile int    `json:"tile,omitempty"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = tmp.Tile
	return nil
}

type EventChi struct {
	Who  Wind  `json:"who"`
	Call *Call `json:"call"`
}

func (event *EventChi) GetType() EventType {
	return EventTypeChi
}

func (event *EventChi) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}{
		Who:  event.Who.String(),
		Call: event.Call,
	})
}

func (event *EventChi) UnmarshalJSON(data []byte) error {
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

type EventPon struct {
	Who  Wind  `json:"who"`
	Call *Call `json:"call"`
}

func (event *EventPon) GetType() EventType {
	return EventTypePon
}

func (event *EventPon) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Call *Call  `json:"call"`
	}{
		Who:  event.Who.String(),
		Call: event.Call,
	})
}

func (event *EventPon) UnmarshalJSON(data []byte) error {
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
	Who Wind `json:"who"`
}

func (event *EventRiichi) GetType() EventType {
	return EventTypeRiichi
}

func (event *EventRiichi) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who string `json:"who"`
	}{
		Who: event.Who.String(),
	})
}

func (event *EventRiichi) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who string `json:"who"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	return nil
}

type EventRon struct {
	Who  Wind `json:"who"`
	Tile int  `json:"tile"`
}

func (event *EventRon) GetType() EventType {
	return EventTypeRon
}

func (event *EventRon) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Tile int    `json:"tile"`
	}{
		Who:  event.Who.String(),
		Tile: event.Tile,
	})
}

func (event *EventRon) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Tile int    `json:"tile"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = tmp.Tile
	return nil
}

type EventTsumo struct {
	Who  Wind `json:"who"`
	Tile int  `json:"tile"`
}

func (event *EventTsumo) GetType() EventType {
	return EventTypeTsumo
}

func (event *EventTsumo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Tile int    `json:"tile"`
	}{
		Who:  event.Who.String(),
		Tile: event.Tile,
	})
}

func (event *EventTsumo) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Tile int    `json:"tile"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = tmp.Tile
	return nil
}

type EventNewIndicator struct {
	Tile int `json:"tile"`
}

func (event *EventNewIndicator) GetType() EventType {
	return EventTypeNewIndicator
}

func (event *EventNewIndicator) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Tile int `json:"tile"`
	}{
		Tile: event.Tile,
	})
}

func (event *EventNewIndicator) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Tile int `json:"tile"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Tile = tmp.Tile
	return nil
}

type EventChanKan struct {
	Who  Wind `json:"who"`
	Tile int  `json:"tile"`
}

func (event *EventChanKan) GetType() EventType {
	return EventTypeChanKan
}

func (event *EventChanKan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who  string `json:"who"`
		Tile int    `json:"tile"`
	}{
		Who:  event.Who.String(),
		Tile: event.Tile,
	})
}

func (event *EventChanKan) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who  string `json:"who"`
		Tile int    `json:"tile"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	event.Who = MapStringToWind[tmp.Who]
	event.Tile = tmp.Tile
	return nil
}

type EventRyuuKyoku struct {
	Who    *Wind           `json:"who,omitempty"`
	Reason RyuuKyokuReason `json:"reason,int"`
}

func (event *EventRyuuKyoku) GetType() EventType {
	return EventTypeRyuuKyoku
}

func (event *EventRyuuKyoku) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Who    string `json:"who,omitempty"`
		Reason string `json:"reason"`
	}{
		Who:    event.Who.String(),
		Reason: event.Reason.String(),
	})
}

func (event *EventRyuuKyoku) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Who    string `json:"who,omitempty"`
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	who := MapStringToWind[tmp.Who]
	reason := MapStringToRyuuKyokuReason[tmp.Reason]
	event.Who = &who
	event.Reason = reason
	return nil
}
