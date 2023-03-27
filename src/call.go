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
	call.CallTiles = tmp.CallTiles
	call.CallTilesFromWho = callTilesFromWho
	return nil
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
	panic("call not in calls!")
}
