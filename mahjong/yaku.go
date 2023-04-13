package mahjong

import (
	"encoding/json"
	"errors"
)

type YakuSet map[Yaku]int
type Yakumans []Yakuman

func (yakuSet *YakuSet) MarshalJSON() ([]byte, error) {
	yakuSetMap := make(map[string]int)
	for yaku, count := range *yakuSet {
		yakuSetMap[yaku.String()] = count
	}
	return json.Marshal(yakuSetMap)
}

func (yakuSet *YakuSet) UnmarshalJSON(data []byte) error {
	yakuSetMap := make(map[string]int)
	err := json.Unmarshal(data, &yakuSetMap)
	if err != nil {
		return err
	}

	*yakuSet = make(YakuSet)
	for yakuStr, count := range yakuSetMap {
		yaku, ok := MapStringToYaku[yakuStr]
		if !ok {
			return errors.New("invalid yaku: " + yakuStr)
		}
		(*yakuSet)[yaku] = count
	}
	return nil
}

func (yakumans *Yakumans) MarshalJSON() ([]byte, error) {
	str := make([]string, len(*yakumans))
	for i, yakuman := range *yakumans {
		str[i] = yakuman.String()
	}
	return json.Marshal(str)
}

func (yakumans *Yakumans) UnmarshalJSON(data []byte) error {
	var str []string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*yakumans = make(Yakumans, len(str))
	for i, yakumanStr := range str {
		yakuman, ok := MapStringToYakuman[yakumanStr]
		if !ok {
			return errors.New("invalid yakuman: " + yakumanStr)
		}
		(*yakumans)[i] = yakuman
	}

	return nil
}

type Yakuman int

// Yakuman numbers are now fixed and should not be changed
//
//go:generate stringer -type=Yakuman
const (
	YakumanNone          Yakuman = 0
	YakumanKokushi       Yakuman = 1
	YakumanKokushi13     Yakuman = 2
	YakumanSuukantsu     Yakuman = 3
	YakumanSuuankou      Yakuman = 4
	YakumanSuuankouTanki Yakuman = 5
	YakumanDaisangen     Yakuman = 6
	YakumanShousuushi    Yakuman = 7
	YakumanDaisuushi     Yakuman = 8
	YakumanRyuuiisou     Yakuman = 9
	YakumanTsuiisou      Yakuman = 10
	YakumanChinrouto     Yakuman = 11
	YakumanChuurenpooto  Yakuman = 12
	YakumanChuurenpooto9 Yakuman = 13
	YakumanTenhou        Yakuman = 14
	YakumanChihou        Yakuman = 15
	YakumanRenhou        Yakuman = 16
)

var MapStringToYakuman = func() map[string]Yakuman {
	m := make(map[string]Yakuman)
	for i := YakumanNone; i <= YakumanRenhou; i++ {
		m[i.String()] = i
	}
	return m
}()

type Yaku int

// Yaku numbers are now fixed and should not be changed
//
//go:generate stringer -type=Yaku
const (
	YakuNone           Yaku = 0
	YakuRiichi         Yaku = 1
	YakuDaburi         Yaku = 2
	YakuIppatsu        Yaku = 3
	YakuTsumo          Yaku = 4
	YakuTanyao         Yaku = 5
	YakuChanta         Yaku = 6
	YakuJunchan        Yaku = 7
	YakuHonrouto       Yaku = 8
	YakuYakuhai        Yaku = 9
	YakuHaku           Yaku = 10
	YakuHatsu          Yaku = 11
	YakuChun           Yaku = 12
	YakuWindRound      Yaku = 13
	YakuWindSelf       Yaku = 14
	YakuTon            Yaku = 15
	YakuNan            Yaku = 16
	YakuSja            Yaku = 17
	YakuPei            Yaku = 18
	YakuTonSelf        Yaku = 19
	YakuNanSelf        Yaku = 20
	YakuSjaSelf        Yaku = 21
	YakuPeiSelf        Yaku = 22
	YakuTonRound       Yaku = 23
	YakuNanRound       Yaku = 24
	YakuSjaRound       Yaku = 25
	YakuPeiRound       Yaku = 26
	YakuChiitoi        Yaku = 27
	YakuToitoi         Yaku = 28
	YakuSanankou       Yaku = 29
	YakuSankantsu      Yaku = 30
	YakuSanshoku       Yaku = 31
	YakuShousangen     Yaku = 32
	YakuPinfu          Yaku = 33
	YakuIppeiko        Yaku = 34
	YakuRyanpeikou     Yaku = 35
	YakuItsuu          Yaku = 36
	YakuSanshokuDoukou Yaku = 37
	YakuHonitsu        Yaku = 38
	YakuChinitsu       Yaku = 39
	YakuDora           Yaku = 40
	YakuUraDora        Yaku = 41
	YakuAkaDora        Yaku = 42
	YakuRenhou         Yaku = 43
	YakuHaitei         Yaku = 44
	YakuHoutei         Yaku = 45
	YakuRinshan        Yaku = 46
	YakuChankan        Yaku = 47
)

var MapStringToYaku = func() map[string]Yaku {
	m := make(map[string]Yaku)
	for i := YakuNone; i <= YakuChankan; i++ {
		m[i.String()] = i
	}
	return m
}()

type Fu int

//go:generate stringer -type=Fu
const (
	FuNone Fu = iota
	FuBase
	FuBaseClosedRon
	FuBase7
	FuSet
	FuTsumo
	FuMeld
	FuNoOpenFu
	FuBadWait
	FuPair
)

var MapStringToFu = func() map[string]Fu {
	m := make(map[string]Fu)
	for i := FuNone; i <= FuPair; i++ {
		m[i.String()] = i
	}
	return m
}()

type Fus []*FuInfo

type FuInfo struct {
	Fu     Fu  `json:"fu"`
	Points int `json:"fu_points"`
}

func (fuInfo *FuInfo) MarshalJSON() ([]byte, error) {
	fuStr := fuInfo.Fu.String()
	fuInfoMap := map[string]interface{}{
		"fu":        fuStr,
		"fu_points": fuInfo.Points,
	}
	return json.Marshal(fuInfoMap)
}

func (fuInfo *FuInfo) UnmarshalJSON(data []byte) error {
	var fuInfoMap map[string]interface{}
	err := json.Unmarshal(data, &fuInfoMap)
	if err != nil {
		return err
	}

	fuStr, ok := fuInfoMap["fu"].(string)
	if !ok {
		return errors.New("missing or invalid fu field")
	}

	points, ok := fuInfoMap["fu_points"].(float64)
	if !ok {
		return errors.New("missing or invalid fu_points field")
	}

	fu, ok := MapStringToFu[fuStr]
	if !ok {
		return errors.New("invalid fu: " + fuStr)
	}

	fuInfo.Fu = fu
	fuInfo.Points = int(points)
	return nil
}
