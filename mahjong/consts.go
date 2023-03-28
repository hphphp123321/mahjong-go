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

type WindRound int

//go:generate stringer -type=WindRound
const (
	WindRoundEast1 WindRound = iota
	WindRoundEast2
	WindRoundEast3
	WindRoundEast4
	WindRoundSouth1
	WindRoundSouth2
	WindRoundSouth3
	WindRoundSouth4
	WindRoundWest1
	WindRoundWest2
	WindRoundWest3
	WindRoundWest4
	WindRoundNorth1
	WindRoundNorth2
	WindRoundNorth3
	WindRoundNorth4
)

var MapStringToWindRound = func() map[string]WindRound {
	m := make(map[string]WindRound)
	for i := WindRoundEast1; i <= WindRoundNorth4; i++ {
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
	EventTypeStart
	EventTypeEnd
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
