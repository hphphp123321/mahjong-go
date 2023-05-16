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
	WindRoundDummy WindRound = iota
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
	EventTypeTenpaiEnd
	EventTypeGlobalInit
)

var MapStringToEventType = func() map[string]EventType {
	m := make(map[string]EventType)
	for i := EventTypeGet; i <= EventTypeGlobalInit; i++ {
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
	Man1T1
	Man1T2
	Man1T3
	Man1T4
	Man2T1
	Man2T2
	Man2T3
	Man2T4
	Man3T1
	Man3T2
	Man3T3
	Man3T4
	Man4T1
	Man4T2
	Man4T3
	Man4T4
	Man5T1
	Man5T2
	Man5T3
	Man5T4
	Man6T1
	Man6T2
	Man6T3
	Man6T4
	Man7T1
	Man7T2
	Man7T3
	Man7T4
	Man8T1
	Man8T2
	Man8T3
	Man8T4
	Man9T1
	Man9T2
	Man9T3
	Man9T4
	Pin1T1
	Pin1T2
	Pin1T3
	Pin1T4
	Pin2T1
	Pin2T2
	Pin2T3
	Pin2T4
	Pin3T1
	Pin3T2
	Pin3T3
	Pin3T4
	Pin4T1
	Pin4T2
	Pin4T3
	Pin4T4
	Pin5T1
	Pin5T2
	Pin5T3
	Pin5T4
	Pin6T1
	Pin6T2
	Pin6T3
	Pin6T4
	Pin7T1
	Pin7T2
	Pin7T3
	Pin7T4
	Pin8T1
	Pin8T2
	Pin8T3
	Pin8T4
	Pin9T1
	Pin9T2
	Pin9T3
	Pin9T4
	Sou1T1
	Sou1T2
	Sou1T3
	Sou1T4
	Sou2T1
	Sou2T2
	Sou2T3
	Sou2T4
	Sou3T1
	Sou3T2
	Sou3T3
	Sou3T4
	Sou4T1
	Sou4T2
	Sou4T3
	Sou4T4
	Sou5T1
	Sou5T2
	Sou5T3
	Sou5T4
	Sou6T1
	Sou6T2
	Sou6T3
	Sou6T4
	Sou7T1
	Sou7T2
	Sou7T3
	Sou7T4
	Sou8T1
	Sou8T2
	Sou8T3
	Sou8T4
	Sou9T1
	Sou9T2
	Sou9T3
	Sou9T4
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
	for i := Man1T1; i <= Chun4; i++ {
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

func (ts TileClass) To4Tiles() Tiles {
	return Tiles{Tile(ts * 4), Tile(ts*4 + 1), Tile(ts*4 + 2), Tile(ts*4 + 3)}
}

var MapStringToTileClass = func() map[string]TileClass {
	m := make(map[string]TileClass)
	for i := Man1; i <= RedSou5; i++ {
		m[i.String()] = i
	}
	return m
}()

var YaoKyuTileClasses = TileClasses{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

var TileClassMap = map[Tile]TileClass{0: 0, 1: 0, 2: 0, 3: 0, 4: 1, 5: 1, 6: 1, 7: 1, 8: 2, 9: 2, 10: 2, 11: 2, 12: 3,
	13: 3, 14: 3, 15: 3, 16: 34, 17: 4, 18: 4, 19: 4, 20: 5, 21: 5, 22: 5, 23: 5, 24: 6, 25: 6, 26: 6, 27: 6, 28: 7,
	29: 7, 30: 7, 31: 7, 32: 8, 33: 8, 34: 8, 35: 8, 36: 9, 37: 9, 38: 9, 39: 9, 40: 10, 41: 10, 42: 10, 43: 10,
	44: 11, 45: 11, 46: 11, 47: 11, 48: 12, 49: 12, 50: 12, 51: 12, 52: 35, 53: 13, 54: 13, 55: 13, 56: 14, 57: 14,
	58: 14, 59: 14, 60: 15, 61: 15, 62: 15, 63: 15, 64: 16, 65: 16, 66: 16, 67: 16, 68: 17, 69: 17, 70: 17, 71: 17,
	72: 18, 73: 18, 74: 18, 75: 18, 76: 19, 77: 19, 78: 19, 79: 19, 80: 20, 81: 20, 82: 20, 83: 20, 84: 21, 85: 21,
	86: 21, 87: 21, 88: 36, 89: 22, 90: 22, 91: 22, 92: 23, 93: 23, 94: 23, 95: 23, 96: 24, 97: 24, 98: 24, 99: 24,
	100: 25, 101: 25, 102: 25, 103: 25, 104: 26, 105: 26, 106: 26, 107: 26, 108: 27, 109: 27, 110: 27, 111: 27,
	112: 28, 113: 28, 114: 28, 115: 28, 116: 29, 117: 29, 118: 29, 119: 29, 120: 30, 121: 30, 122: 30, 123: 30,
	124: 31, 125: 31, 126: 31, 127: 31, 128: 32, 129: 32, 130: 32, 131: 32, 132: 33, 133: 33, 134: 33, 135: 33, -1: -1}

var TileClassUTF = map[TileClass]string{
	0: "🀇", 1: "🀈", 2: "🀉", 3: "🀊", 4: "🀋", 5: "🀌", 6: "🀍", 7: "🀎", 8: "🀏",
	9: "🀙", 10: "🀚", 11: "🀛", 12: "🀜", 13: "🀝", 14: "🀞", 15: "🀟", 16: "🀠", 17: "🀡",
	18: "🀐", 19: "🀑", 20: "🀒", 21: "🀓", 22: "🀔", 23: "🀕", 24: "🀖", 25: "🀗", 26: "🀘",
	27: "🀀", 28: "🀁", 29: "🀂", 30: "🀃", 31: "🀆", 32: "🀅", 33: "🀄",
	34: "[🀋]", 35: "[🀝]", 36: "[🀔]",
}

var ErrGameEnd = errors.New("game is end")

// Limit numbers are now fixed and should not be changed
type Limit int

//go:generate stringer -type=Limit
const (
	LimitNone Limit = iota
	LimitMangan
	LimitHaneman
	LimitBaiman
	LimitSanbaiman
	LimitYakuman
)

var MapStringToLimit = func() map[string]Limit {
	m := make(map[string]Limit)
	for i := LimitNone; i <= LimitYakuman; i++ {
		m[i.String()] = i
	}
	return m
}()

type EndType int

const (
	EndTypeNone EndType = iota
	EndTypeRound
	EndTypeGame
)
