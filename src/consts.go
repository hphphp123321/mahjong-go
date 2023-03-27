package mahjong

type Wind int

//go:generate stringer -type=Wind
const (
	Dummy Wind = -1 + iota
	East
	South
	West
	North
)

var MapStringToWind = func() map[string]Wind {
	m := make(map[string]Wind)
	for i := Dummy; i <= North; i++ {
		m[i.String()] = i
	}
	return m
}()

type CallType int

//go:generate stringer -type=CallType
const (
	Get CallType = -1 + iota
	Skip
	Discard
	Chi
	Pon
	DaiMinKan
	ShouMinKan
	AnKan
	Riichi
	Ron
	Tsumo
	KyuShuKyuHai
	ChanKan
)

var MapStringToCallType = func() map[string]CallType {
	m := make(map[string]CallType)
	for i := Get; i <= ChanKan; i++ {
		m[i.String()] = i
	}
	return m
}()

type EventType int

//go:generate stringer -type=EventType
const (
	EventTypeGet EventType = -1 + iota
	EventTypeTsumoGiri
	EventTypeDiscard
	EventTypeChi
	EventTypePon
	EventTypeDaiMinKan
	EventTypeShouMinKan
	EventTypeAnKan
	EventTypeRiichi
	EventTypeRon
	EventTypeTsumo
	EventTypeNewIndicator
	EventTypeChanKan
	EventTypeRyuuKyoku
)

var MapStringToEventType = func() map[string]EventType {
	m := make(map[string]EventType)
	for i := EventTypeGet; i <= EventTypeRyuuKyoku; i++ {
		m[i.String()] = i
	}
	return m
}()

type RyuuKyokuReason int

//go:generate stringer -type=RyuuKyokuReason
const (
	RyuuKyokuNormal RyuuKyokuReason = iota
	RyuuKyokuKyuShuKyuHai
	RyuuKyokuShuChaRiichi
	RyuuKyokuSuuKaiKan
	RyuuKyokuSuufonRenda
	RyuuKyokuSanChaHou
)

var MapStringToRyuuKyokuReason = func() map[string]RyuuKyokuReason {
	m := make(map[string]RyuuKyokuReason)
	for i := RyuuKyokuNormal; i <= RyuuKyokuSanChaHou; i++ {
		m[i.String()] = i
	}
	return m
}()
