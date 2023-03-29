// Code generated by "stringer -type=CallType"; DO NOT EDIT.

package mahjong

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Get - -1]
	_ = x[Skip-0]
	_ = x[Discard-1]
	_ = x[Chi-2]
	_ = x[Pon-3]
	_ = x[DaiMinKan-4]
	_ = x[ShouMinKan-5]
	_ = x[AnKan-6]
	_ = x[Riichi-7]
	_ = x[Ron-8]
	_ = x[Tsumo-9]
	_ = x[KyuShuKyuHai-10]
	_ = x[ChanKan-11]
	_ = x[Next-12]
}

const _CallType_name = "GetSkipDiscardChiPonDaiMinKanShouMinKanAnKanRiichiRonTsumoKyuShuKyuHaiChanKanNext"

var _CallType_index = [...]uint8{0, 3, 7, 14, 17, 20, 29, 39, 44, 50, 53, 58, 70, 77, 81}

func (i CallType) String() string {
	i -= -1
	if i < 0 || i >= CallType(len(_CallType_index)-1) {
		return "CallType(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _CallType_name[_CallType_index[i]:_CallType_index[i+1]]
}
