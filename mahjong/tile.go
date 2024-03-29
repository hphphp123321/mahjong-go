package mahjong

import (
	"encoding/json"
	"errors"
	"github.com/hphphp123321/go-common"
	"math/rand"
	"sort"
)

const NumTiles int = 136

func (tile Tile) Class() TileClass {
	if tile == TileDummy {
		return TileClassDummy
	}
	return TileClass(tile / 4)
}

func (tile Tile) UTF8() string {
	return TileClassUTF[TileClassMap[tile]]
}

type Tiles []Tile

func (tiles *Tiles) UTF8() string {
	var s string
	for i, tile := range *tiles {
		if i != 0 {
			s += " "
		}
		s += TileClassUTF[TileClassMap[tile]]
	}
	return s
}

func (tiles *Tiles) String() string {
	var s string
	for i, tile := range *tiles {
		if tile == TileDummy {
			continue
		}
		if i != 0 {
			s += " "
		}
		s += tile.String()
	}
	return s
}

func (tiles *Tiles) Classes() *TileClasses {
	var tileClasses = TileClasses(common.MapSlice(*tiles, func(t Tile) TileClass { return t.Class() }))
	return &tileClasses
}

type TileT struct {
	tileId      Tile
	discardable bool
	isRinshan   bool
	isLast      bool
	discardWind Wind
}

func newTile(tileId Tile) *TileT {
	return &TileT{
		tileId:      tileId,
		discardable: true,
		isRinshan:   false,
		isLast:      false,
		discardWind: -1}
}

type MahjongTiles struct {
	randP *rand.Rand

	allTiles       map[Tile]*TileT
	tiles          Tiles
	kanNum         int
	NumRemainTiles int

	tilePointer    int
	rinshanPointer int
}

func NewMahjongTiles(randP *rand.Rand) *MahjongTiles {
	if randP == nil {
		randP = rand.New(rand.NewSource(1))
	}
	mahjongTiles := MahjongTiles{
		allTiles: make(map[Tile]*TileT, NumTiles),
		tiles:    make(Tiles, NumTiles),
		randP:    randP,
	}
	for i := 0; i < NumTiles; i++ {
		mahjongTiles.allTiles[Tile(i)] = newTile(Tile(i))
		mahjongTiles.tiles[i] = Tile(i)
	}
	return &mahjongTiles
}

func (tiles *MahjongTiles) Reset() {

	for i := 0; i < NumTiles; i++ {
		tiles.allTiles[Tile(i)] = newTile(Tile(i))
		tiles.tiles[i] = Tile(i)
	}
	tiles.randP.Shuffle(NumTiles, func(i, j int) {
		tiles.tiles[i], tiles.tiles[j] = tiles.tiles[j], tiles.tiles[i]
	})

	tiles.kanNum = 0
	tiles.NumRemainTiles = 70
	tiles.tilePointer = 52
	tiles.rinshanPointer = 135
}

// Setup
//
//	@Description: setup tiles for each player
//	@receiver tiles
//	@param ts: prepared tiles, len must be 136, nil for default random tiles
//	@return map[Wind]Tiles
func (tiles *MahjongTiles) Setup(ts Tiles) map[Wind]Tiles {
	if ts != nil {
		if len(ts) != NumTiles {
			panic(errors.New("len of prepared tiles must be 136"))
		}
		tiles.tiles = ts
	}
	t := tiles.tiles[0:13]
	tonTiles := t.Copy()
	t = tiles.tiles[13:26]
	nanTiles := t.Copy()
	t = tiles.tiles[26:39]
	shaaTiles := t.Copy()
	t = tiles.tiles[39:52]
	peiTiles := t.Copy()
	return map[Wind]Tiles{
		East:  tonTiles,
		South: nanTiles,
		West:  shaaTiles,
		North: peiTiles,
	}
}

func (tiles *MahjongTiles) DealTile(isRinshan bool) Tile {
	if tiles.NumRemainTiles <= 0 {
		panic(errors.New("no more tiles"))
	}
	if isRinshan && tiles.kanNum == 4 {
		panic(errors.New("no more rinshan tiles"))
	}
	tiles.NumRemainTiles--
	var tile Tile
	if isRinshan {
		tile = tiles.tiles[tiles.rinshanPointer]
		tiles.rinshanPointer--
		tiles.kanNum++
		tiles.allTiles[tile].isRinshan = true

	} else {
		tile = tiles.tiles[tiles.tilePointer]
		tiles.tilePointer++
		if tiles.NumRemainTiles == 0 {
			tiles.allTiles[tile].isLast = true
		}
	}
	return tile
}

func (tiles *MahjongTiles) DoraIndicators() Tiles {
	t := make(Tiles, 0, 5)
	for i := 0; i < tiles.kanNum+1; i++ {
		t.Append(tiles.tiles[130-2*i])
	}
	return t
}

func (tiles *MahjongTiles) UraDoraIndicators() Tiles {
	t := make(Tiles, 0, 5)
	for i := 0; i < tiles.kanNum+1; i++ {
		t.Append(tiles.tiles[131-2*i])
	}
	return t
}

func (tiles *MahjongTiles) GetCurrentIndicator() Tile {
	return tiles.DoraIndicators()[tiles.kanNum]
}

func (tiles *Tiles) Remove(tile Tile) {
	for i := 0; i < len(*tiles); i++ {
		if (*tiles)[i] == tile {
			*tiles = append((*tiles)[:i], (*tiles)[i+1:]...)
			return
		}
	}
	panic("tile: " + tile.String() + "not in tiles: " + tiles.String())
}

func (tiles *Tiles) Append(tile Tile) {
	*tiles = append(*tiles, tile)
}

func TilesEqual(tiles1 Tiles, tiles2 Tiles) bool {
	newArray1 := append(Tiles{}, tiles1...)
	newArray2 := append(Tiles{}, tiles2...)
	sort.Sort(&newArray1)
	sort.Sort(&newArray2)
	for i, tile := range newArray1 {
		if TileClassMap[newArray2[i]] != TileClassMap[tile] {
			return false
		}
	}
	return true
}

func (tiles *Tiles) MarshalJSON() ([]byte, error) {
	str := make([]string, 0)
	for _, t := range *tiles {
		if t == TileDummy {
			continue
		}
		str = append(str, t.String())
	}
	return json.Marshal(str)
}
func (tiles *Tiles) UnmarshalJSON(data []byte) error {
	var str []string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*tiles = make(Tiles, 0)
	for _, t := range str {
		*tiles = append(*tiles, MapStringToTile[t])
	}
	return nil
}

// Len sort.Interface
func (tiles *Tiles) Len() int {
	return len(*tiles)
}

// Less sort.Interface
func (tiles *Tiles) Less(i, j int) bool {
	return (*tiles)[i] < (*tiles)[j]
}

// Swap sort.Interface
func (tiles *Tiles) Swap(i, j int) {
	(*tiles)[i], (*tiles)[j] = (*tiles)[j], (*tiles)[i]
}

func (tiles *Tiles) Copy() Tiles {
	tilesCopy := make(Tiles, len(*tiles), cap(*tiles))
	copy(tilesCopy, *tiles)
	return tilesCopy
}

func (tiles *Tiles) Index(tileId Tile, startIdx int) int {
	for i := startIdx; i < len(*tiles); i++ {
		if (*tiles)[i] == tileId {
			return i
		}
	}
	return -1
}

func (tiles *Tiles) Count(tileId Tile) int {
	count := 0
	for _, tile := range *tiles {
		if tile == tileId {
			count++
		}
	}
	return count
}

type TileClasses []TileClass

func (tileClasses *TileClasses) Append(tileClass TileClass) {
	*tileClasses = append(*tileClasses, tileClass)
}

func (tileClasses *TileClasses) Remove(tileClass TileClass) {
	for i := 0; i < len(*tileClasses); i++ {
		if (*tileClasses)[i] == tileClass {
			*tileClasses = append((*tileClasses)[:i], (*tileClasses)[i+1:]...)
			return
		}
	}
	panic("tileClass: " + tileClass.String() + "not in tileClasses: " + tileClasses.String())
}

func (tileClasses *TileClasses) Index(tileClass TileClass, startIdx int) int {
	for i := startIdx; i < len(*tileClasses); i++ {
		if (*tileClasses)[i] == tileClass {
			return i
		}
	}
	return -1
}

func (tileClasses *TileClasses) Count(tileClass TileClass) int {
	count := 0
	for _, tileC := range *tileClasses {
		if tileC == tileClass {
			count++
		}
	}
	return count
}

func (tileClasses *TileClasses) Copy() TileClasses {
	tileClassesCopy := make(TileClasses, len(*tileClasses), cap(*tileClasses))
	copy(tileClassesCopy, *tileClasses)
	return tileClassesCopy
}

func (tileClasses *TileClasses) String() string {
	var s string
	for i, tileClass := range *tileClasses {
		if tileClass == TileClassDummy {
			continue
		}
		if i != 0 {
			s += " "
		}
		s += tileClass.String()
	}
	return s
}

func (tileClasses *TileClasses) MarshalJSON() ([]byte, error) {
	str := make([]string, len(*tileClasses))
	for i, t := range *tileClasses {
		str[i] = t.String()
	}
	return json.Marshal(str)
}

func (tileClasses *TileClasses) UnmarshalJSON(data []byte) error {
	var str []string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*tileClasses = make(TileClasses, len(str))
	for i, t := range str {
		(*tileClasses)[i] = MapStringToTileClass[t]
	}
	return nil
}

type TenpaiInfos map[Tile]*TenpaiInfo

func (tenpaiInfos *TenpaiInfos) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TenpaiInfos map[string]*TenpaiInfo `json:"tenpai_infos"`
	}{
		TenpaiInfos: func() map[string]*TenpaiInfo {
			tenpaiInfosMap := make(map[string]*TenpaiInfo)
			for tile, tenpaiInfo := range *tenpaiInfos {
				tenpaiInfosMap[tile.String()] = tenpaiInfo
			}
			return tenpaiInfosMap
		}(),
	})
}

func (tenpaiInfos *TenpaiInfos) UnmarshalJSON(data []byte) error {
	var str map[string]*TenpaiInfo
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*tenpaiInfos = make(TenpaiInfos)
	for tile, tenpaiInfo := range str {
		(*tenpaiInfos)[MapStringToTile[tile]] = tenpaiInfo
	}
	return nil
}

type TenpaiInfo struct {
	TileClassesTenpaiResult map[TileClass]*TenpaiResult `json:"tile_classes_tenpai_result"`
	Furiten                 bool                        `json:"furiten"`
}

func NewTenpaiInfo() *TenpaiInfo {
	return &TenpaiInfo{
		TileClassesTenpaiResult: make(map[TileClass]*TenpaiResult),
		Furiten:                 false,
	}
}

type TenpaiResult struct {
	RemainNum int     `json:"remain_num"`
	Result    *Result `json:"result"`
}

func NewTenpaiResult(remainNum int, result *Result) *TenpaiResult {
	return &TenpaiResult{
		RemainNum: remainNum,
		Result:    result,
	}
}

func (tenpaiInfo *TenpaiInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TileClassesTenpaiResult map[string]*TenpaiResult `json:"tile_classes_tenpai_result"`
		Furiten                 bool                     `json:"furiten"`
	}{
		TileClassesTenpaiResult: func() map[string]*TenpaiResult {
			tileClassesRemainNum := make(map[string]*TenpaiResult)
			for tileClass, result := range tenpaiInfo.TileClassesTenpaiResult {
				tileClassesRemainNum[tileClass.String()] = result
			}
			return tileClassesRemainNum
		}(),
		Furiten: tenpaiInfo.Furiten,
	})
}

func (tenpaiInfo *TenpaiInfo) UnmarshalJSON(data []byte) error {
	var tmp struct {
		TileClassesTenpaiResult map[string]*TenpaiResult `json:"tile_classes_tenpai_result"`
		Furiten                 bool                     `json:"furiten"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tenpaiInfo.TileClassesTenpaiResult = make(map[TileClass]*TenpaiResult)
	for tileClass, result := range tmp.TileClassesTenpaiResult {
		tenpaiInfo.TileClassesTenpaiResult[MapStringToTileClass[tileClass]] = result
	}
	tenpaiInfo.Furiten = tmp.Furiten
	return nil
}
