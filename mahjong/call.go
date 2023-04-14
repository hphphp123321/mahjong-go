package mahjong

import "encoding/json"

type Calls []*Call

type Call struct {
	CallType         CallType `json:"type"`
	CallTiles        Tiles    `json:"tiles"`
	CallTilesFromWho []Wind   `json:"who"`
}

func (call *Call) MarshalJSON() ([]byte, error) {
	var callTilesFromWho []string
	for _, w := range call.CallTilesFromWho {
		if w == WindDummy {
			continue
		}
		callTilesFromWho = append(callTilesFromWho, w.String())
	}
	return json.Marshal(&struct {
		CallType         string   `json:"type"`
		CallTiles        Tiles    `json:"tiles"`
		CallTilesFromWho []string `json:"who"`
	}{
		CallType:         call.CallType.String(),
		CallTiles:        call.CallTiles,
		CallTilesFromWho: callTilesFromWho,
	})
}

func (call *Call) UnmarshalJSON(data []byte) error {
	var callTilesFromWho []Wind
	var tmp struct {
		CallType         string   `json:"type"`
		CallTiles        Tiles    `json:"tiles"`
		CallTilesFromWho []string `json:"who"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	call.CallType = MapStringToCallType[tmp.CallType]
	for _, w := range tmp.CallTilesFromWho {
		callTilesFromWho = append(callTilesFromWho, MapStringToWind[w])
	}
	for len(tmp.CallTiles) < 4 {
		tmp.CallTiles = append(tmp.CallTiles, TileDummy)
	}
	for len(callTilesFromWho) < 4 {
		callTilesFromWho = append(callTilesFromWho, WindDummy)
	}
	call.CallTiles = tmp.CallTiles
	call.CallTilesFromWho = callTilesFromWho
	return nil
}

func (call *Call) Copy() *Call {
	tilesFromWho := make([]Wind, len(call.CallTilesFromWho))
	copy(tilesFromWho, call.CallTilesFromWho)
	return &Call{
		CallType:         call.CallType,
		CallTiles:        call.CallTiles.Copy(),
		CallTilesFromWho: tilesFromWho,
	}
}

func NewCall(meldType CallType, CallTiles Tiles, CallTilesFromWho []Wind) *Call {
	return &Call{
		CallType:         meldType,
		CallTiles:        CallTiles,
		CallTilesFromWho: CallTilesFromWho,
	}
}

func CallEqual(call1 *Call, call2 *Call) bool {
	if call1.CallType != call2.CallType {
		return false
	}
	if TilesEqual(call1.CallTiles, call2.CallTiles) {
		return true
	}
	return false
}

func (calls *Calls) Index(call *Call) int {
	for idx, c := range *calls {
		if CallEqual(c, call) {
			return idx
		}
	}
	return -1
}

func (calls *Calls) Copy() Calls {
	callsCopy := make(Calls, len(*calls), cap(*calls))
	copy(callsCopy, *calls)
	return callsCopy
}

func (calls *Calls) Append(call *Call) {
	*calls = append(*calls, call)
}

func (calls *Calls) Remove(call *Call) {
	idx := calls.Index(call)
	*calls = append((*calls)[:idx], (*calls)[idx+1:]...)
}

var SkipCall = &Call{
	CallType:         Skip,
	CallTiles:        Tiles{TileDummy, TileDummy, TileDummy, TileDummy},
	CallTilesFromWho: []Wind{WindDummy, WindDummy, WindDummy, WindDummy},
}
