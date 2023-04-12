package mahjong

import "errors"

type Wind int

//go:generate stringer -type=Wind
const (
	WindDummy Wind = -1 + iota
	East
	South
	West
	North
)

var MapStringToWind = func() map[string]Wind {
	m := make(map[string]Wind)
	for i := WindDummy; i <= North; i++ {
		m[i.String()] = i
	}
	return m
}()

type WindRound int

//go:generate stringer -type=WindRound
const (
	WindRoundDummy WindRound = -1 + iota
	WindRoundEast1
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
	KyuuShuKyuuHai
	ChanKan
	Next
)

var MapStringToCallType = func() map[string]CallType {
	m := make(map[string]CallType)
	for i := Get; i <= Next; i++ {
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
	EventTypeFuriten
	EventTypeNagashiMangan
	EventTypeTenhaiEnd
)

var MapStringToEventType = func() map[string]EventType {
	m := make(map[string]EventType)
	for i := EventTypeGet; i <= EventTypeTenhaiEnd; i++ {
		m[i.String()] = i
	}
	return m
}()

type RyuuKyokuReason int

//go:generate stringer -type=RyuuKyokuReason
const (
	NoRyuuKyoku RyuuKyokuReason = iota
	RyuuKyokuNormal
	RyuuKyokuKyuuShuKyuuHai
	RyuuKyokuSuuChaRiichi
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

type FuritenReason int

//go:generate stringer -type=FuritenReason
const (
	FuritenNone FuritenReason = iota
	FuritenJun
	FuritenDiscard
	FuritenRiichi
)

var MapStringToFuritenReason = func() map[string]FuritenReason {
	m := make(map[string]FuritenReason)
	for i := FuritenJun; i <= FuritenRiichi; i++ {
		m[i.String()] = i
	}
	return m
}()

type Tile int

//go:generate stringer -type=Tile
const (
	TileDummy Tile = -1 + iota
	Man11
	Man12
	Man13
	Man14
	Man21
	Man22
	Man23
	Man24
	Man31
	Man32
	Man33
	Man34
	Man41
	Man42
	Man43
	Man44
	Man51
	Man52
	Man53
	Man54
	Man61
	Man62
	Man63
	Man64
	Man71
	Man72
	Man73
	Man74
	Man81
	Man82
	Man83
	Man84
	Man91
	Man92
	Man93
	Man94
	Pin11
	Pin12
	Pin13
	Pin14
	Pin21
	Pin22
	Pin23
	Pin24
	Pin31
	Pin32
	Pin33
	Pin34
	Pin41
	Pin42
	Pin43
	Pin44
	Pin51
	Pin52
	Pin53
	Pin54
	Pin61
	Pin62
	Pin63
	Pin64
	Pin71
	Pin72
	Pin73
	Pin74
	Pin81
	Pin82
	Pin83
	Pin84
	Pin91
	Pin92
	Pin93
	Pin94
	Sou11
	Sou12
	Sou13
	Sou14
	Sou21
	Sou22
	Sou23
	Sou24
	Sou31
	Sou32
	Sou33
	Sou34
	Sou41
	Sou42
	Sou43
	Sou44
	Sou51
	Sou52
	Sou53
	Sou54
	Sou61
	Sou62
	Sou63
	Sou64
	Sou71
	Sou72
	Sou73
	Sou74
	Sou81
	Sou82
	Sou83
	Sou84
	Sou91
	Sou92
	Sou93
	Sou94
	Ton1
	Ton2
	Ton3
	Ton4
	Nan1
	Nan2
	Nan3
	Nan4
	Shaa1
	Shaa2
	Shaa3
	Shaa4
	Pei1
	Pei2
	Pei3
	Pei4
	Haku1
	Haku2
	Haku3
	Haku4
	Hatsu1
	Hatsu2
	Hatsu3
	Hatsu4
	Chun1
	Chun2
	Chun3
	Chun4
)

var MapStringToTile = func() map[string]Tile {
	m := make(map[string]Tile)
	for i := Man11; i <= Chun4; i++ {
		m[i.String()] = i
	}
	return m
}()

type TileClass int

//go:generate stringer -type=TileClass
const (
	TileClassDummy TileClass = -1 + iota
	Man1
	Man2
	Man3
	Man4
	Man5
	Man6
	Man7
	Man8
	Man9
	Pin1
	Pin2
	Pin3
	Pin4
	Pin5
	Pin6
	Pin7
	Pin8
	Pin9
	Sou1
	Sou2
	Sou3
	Sou4
	Sou5
	Sou6
	Sou7
	Sou8
	Sou9
	Ton
	Nan
	Shaa
	Pei
	Haku
	Hatsu
	Chun
	RedMan5
	RedPin5
	RedSou5
)

var MapStringToTileClass = func() map[string]TileClass {
	m := make(map[string]TileClass)
	for i := Man1; i <= RedSou5; i++ {
		m[i.String()] = i
	}
	return m
}()

var ErrGameEnd = errors.New("game is end")
