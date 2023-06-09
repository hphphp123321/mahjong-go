// Code generated by "stringer -type=WindRound -trimprefix WindRound"; DO NOT EDIT.

package mahjong

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WindRoundDummy-0]
	_ = x[WindRoundEast1-1]
	_ = x[WindRoundEast2-2]
	_ = x[WindRoundEast3-3]
	_ = x[WindRoundEast4-4]
	_ = x[WindRoundSouth1-5]
	_ = x[WindRoundSouth2-6]
	_ = x[WindRoundSouth3-7]
	_ = x[WindRoundSouth4-8]
	_ = x[WindRoundWest1-9]
	_ = x[WindRoundWest2-10]
	_ = x[WindRoundWest3-11]
	_ = x[WindRoundWest4-12]
	_ = x[WindRoundNorth1-13]
	_ = x[WindRoundNorth2-14]
	_ = x[WindRoundNorth3-15]
	_ = x[WindRoundNorth4-16]
}

const _WindRound_name = "DummyEast1East2East3East4South1South2South3South4West1West2West3West4North1North2North3North4"

var _WindRound_index = [...]uint8{0, 5, 10, 15, 20, 25, 31, 37, 43, 49, 54, 59, 64, 69, 75, 81, 87, 93}

func (i WindRound) String() string {
	if i < 0 || i >= WindRound(len(_WindRound_index)-1) {
		return "WindRound(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _WindRound_name[_WindRound_index[i]:_WindRound_index[i+1]]
}
