package mahjong

import (
	"encoding/json"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
)

type Rule struct {
	GameLength int `json:"game_length"` // game length 1 for 1 round only, 4 for tonpuusen, 8 for hanchan, etc.

	// Yaku Rule
	IsOpenTanyao         bool  `json:"is_open_tanyao"`           // true for open tanyao, false for closed tanyao
	HasAkaDora           bool  `json:"has_aka_dora"`             // true for aka dora, false for no aka dora
	RenhouLimit          Limit `json:"renhou_limit"`             // Limit None for no Limit, Mangan for mangan limit, Haneman for haneman limit, etc.
	IsHaiteiFromLiveOnly bool  `json:"is_haitei_from_live_only"` // true for haitei from live only(means no rinshan and haitei both), false for can rinshan and haitei both
	IsUra                bool  `json:"is_ura"`                   // true for has ura dora, false for no ura dora
	IsIpatsu             bool  `json:"is_ipatsu"`                // true for has ipatsu, false for no ipatsu
	IsGreenRequired      bool  `json:"is_green_required"`        // true for Ryuuiisou(Yakuman) need green(hatsu) required, false for no green required
	IsRinshanFu          bool  `json:"is_rinshan_fu"`            // true for has rinshan fu, false for no rinshan fu

	// Score Rule
	IsManganRound     bool `json:"is_mangan_round"`     // true for points round up to mangan, false for no round
	IsKazoeYakuman    bool `json:"is_kazoe_yakuman"`    // true for has kazoe yakuman(13 han), false for no kazoe yakuman
	HasDoubleYakumans bool `json:"has_double_yakumans"` // true for has double yakuman, false for no double yakuman
	IsYakumanSum      bool `json:"is_yakuman_sum"`      // true for  yakuman, false for no sum yakuman
	HonbaValue        int  `json:"honba_value"`         // HonbaValue represent the value of one honba in score calculation, default is 100

	// Other Rule
	IsSanChaHou     bool `json:"is_san_cha_hou"`    // can san chan ron, true for can, false for can't -> ryuu kyoku
	IsNagashiMangan bool `json:"is_nagashi_mangan"` // can nagashi/ryuukyoku mangan, true for can, false for can't
}

func (r *Rule) YakuRule() *yaku.RulesStruct {
	var akaDora []tile.Instance
	if r.HasAkaDora {
		akaDora = []tile.Instance{
			tile.Man5.Instance(0),
			tile.Pin5.Instance(0),
			tile.Sou5.Instance(0),
		}
	}
	return &yaku.RulesStruct{
		IsOpenTanyao:         r.IsOpenTanyao,
		AkaDoras:             akaDora,
		RenhouLimit:          yaku.Limit(r.RenhouLimit),
		IsHaiteiFromLiveOnly: r.IsHaiteiFromLiveOnly,
		IsUra:                r.IsUra,
		IsIpatsu:             r.IsIpatsu,
		IsGreenRequired:      r.IsGreenRequired,
		IsRinshanFu:          r.IsRinshanFu,
	}
}

func (r *Rule) ScoreRule() *score.RulesStruct {
	var doubleYakuMans map[yaku.Yakuman]bool
	if r.HasDoubleYakumans {
		doubleYakuMans = DefaultDoubleYakumans()
	}
	return &score.RulesStruct{
		IsManganRound:  r.IsManganRound,
		IsKazoeYakuman: r.IsKazoeYakuman,
		DoubleYakumans: doubleYakuMans,
		IsYakumanSum:   r.IsYakumanSum,
		HonbaValue:     score.Money(r.HonbaValue),
	}
}

func GetDefaultRule() *Rule {
	// TenhouRed
	return &Rule{
		GameLength: 8,

		IsOpenTanyao:         true,
		HasAkaDora:           true,
		RenhouLimit:          LimitNone,
		IsHaiteiFromLiveOnly: true,
		IsUra:                true,
		IsIpatsu:             true,
		IsGreenRequired:      false,
		IsRinshanFu:          true,

		IsManganRound:  false,
		IsKazoeYakuman: true,
		IsYakumanSum:   true,
		HonbaValue:     100,

		IsSanChaHou:     false,
		IsNagashiMangan: true,
	}
}

func DefaultDoubleYakumans() map[yaku.Yakuman]bool {
	return map[yaku.Yakuman]bool{
		yaku.YakumanChuurenpooto9: true,
		yaku.YakumanKokushi13:     true,
		yaku.YakumanSuuankouTanki: true,
		yaku.YakumanDaisuushi:     true,
	}
}

func (r *Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		GameLength           int    `json:"game_length"`
		IsOpenTanyao         bool   `json:"is_open_tanyao"`
		HasAkaDora           bool   `json:"has_aka_dora"`
		RenhouLimit          string `json:"renhou_limit"`
		IsHaiteiFromLiveOnly bool   `json:"is_haitei_from_live_only"`
		IsUra                bool   `json:"is_ura"`
		IsIpatsu             bool   `json:"is_ipatsu"`
		IsGreenRequired      bool   `json:"is_green_required"`
		IsRinshanFu          bool   `json:"is_rinshan_fu"`
		IsManganRound        bool   `json:"is_mangan_round"`
		IsKazoeYakuman       bool   `json:"is_kazoe_yakuman"`
		HasDoubleYakumans    bool   `json:"has_double_yakumans"`
		IsYakumanSum         bool   `json:"is_yakuman_sum"`
		HonbaValue           int    `json:"honba_value"`
		IsSanChaHou          bool   `json:"is_san_cha_hou"`
		IsNagashiMangan      bool   `json:"is_nagashi_mangan"`
	}{
		GameLength:           r.GameLength,
		IsOpenTanyao:         r.IsOpenTanyao,
		HasAkaDora:           r.HasAkaDora,
		RenhouLimit:          r.RenhouLimit.String(),
		IsHaiteiFromLiveOnly: r.IsHaiteiFromLiveOnly,
		IsUra:                r.IsUra,
		IsIpatsu:             r.IsIpatsu,
		IsGreenRequired:      r.IsGreenRequired,
		IsRinshanFu:          r.IsRinshanFu,
		IsManganRound:        r.IsManganRound,
		IsKazoeYakuman:       r.IsKazoeYakuman,
		HasDoubleYakumans:    r.HasDoubleYakumans,
		IsYakumanSum:         r.IsYakumanSum,
		HonbaValue:           r.HonbaValue,
		IsSanChaHou:          r.IsSanChaHou,
		IsNagashiMangan:      r.IsNagashiMangan,
	})
}

func (r *Rule) UnmarshalJSON(data []byte) error {
	var s struct {
		GameLength           int    `json:"game_length"`
		IsOpenTanyao         bool   `json:"is_open_tanyao"`
		HasAkaDora           bool   `json:"has_aka_dora"`
		RenhouLimit          string `json:"renhou_limit"`
		IsHaiteiFromLiveOnly bool   `json:"is_haitei_from_live_only"`
		IsUra                bool   `json:"is_ura"`
		IsIpatsu             bool   `json:"is_ipatsu"`
		IsGreenRequired      bool   `json:"is_green_required"`
		IsRinshanFu          bool   `json:"is_rinshan_fu"`
		IsManganRound        bool   `json:"is_mangan_round"`
		IsKazoeYakuman       bool   `json:"is_kazoe_yakuman"`
		HasDoubleYakumans    bool   `json:"has_double_yakumans"`
		IsYakumanSum         bool   `json:"is_yakuman_sum"`
		HonbaValue           int    `json:"honba_value"`
		IsSanChaHou          bool   `json:"is_san_cha_hou"`
		IsNagashiMangan      bool   `json:"is_nagashi_mangan"`
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	r.GameLength = s.GameLength
	r.IsOpenTanyao = s.IsOpenTanyao
	r.HasAkaDora = s.HasAkaDora
	r.RenhouLimit = MapStringToLimit[s.RenhouLimit]
	r.IsHaiteiFromLiveOnly = s.IsHaiteiFromLiveOnly
	r.IsUra = s.IsUra
	r.IsIpatsu = s.IsIpatsu
	r.IsGreenRequired = s.IsGreenRequired
	r.IsRinshanFu = s.IsRinshanFu
	r.IsManganRound = s.IsManganRound
	r.IsKazoeYakuman = s.IsKazoeYakuman
	r.HasDoubleYakumans = s.HasDoubleYakumans
	r.IsYakumanSum = s.IsYakumanSum
	r.HonbaValue = s.HonbaValue
	r.IsSanChaHou = s.IsSanChaHou
	r.IsNagashiMangan = s.IsNagashiMangan
	return nil
}
