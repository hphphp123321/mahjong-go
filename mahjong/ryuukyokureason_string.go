// Code generated by "stringer -type=RyuuKyokuReason"; DO NOT EDIT.

package mahjong

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NoRyuuKyoku-0]
	_ = x[RyuuKyokuNormal-1]
	_ = x[RyuuKyokuKyuShuKyuHai-2]
	_ = x[RyuuKyokuSuuChaRiichi-3]
	_ = x[RyuuKyokuSuuKaiKan-4]
	_ = x[RyuuKyokuSuufonRenda-5]
	_ = x[RyuuKyokuSanChaHou-6]
}

const _RyuuKyokuReason_name = "NoRyuuKyokuRyuuKyokuNormalRyuuKyokuKyuShuKyuHaiRyuuKyokuSuuChaRiichiRyuuKyokuSuuKaiKanRyuuKyokuSuufonRendaRyuuKyokuSanChaHou"

var _RyuuKyokuReason_index = [...]uint8{0, 11, 26, 47, 68, 86, 106, 124}

func (i RyuuKyokuReason) String() string {
	if i < 0 || i >= RyuuKyokuReason(len(_RyuuKyokuReason_index)-1) {
		return "RyuuKyokuReason(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RyuuKyokuReason_name[_RyuuKyokuReason_index[i]:_RyuuKyokuReason_index[i+1]]
}
