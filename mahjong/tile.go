package mahjong

import (
	"errors"
	"math/rand"
	"sort"
)

var TileClassMap = map[int]int{0: 0, 1: 0, 2: 0, 3: 0, 4: 1, 5: 1, 6: 1, 7: 1, 8: 2, 9: 2, 10: 2, 11: 2, 12: 3, 13: 3, 14: 3, 15: 3, 16: 35, 17: 4, 18: 4, 19: 4, 20: 5, 21: 5, 22: 5, 23: 5, 24: 6, 25: 6, 26: 6, 27: 6, 28: 7, 29: 7, 30: 7, 31: 7, 32: 8, 33: 8, 34: 8, 35: 8, 36: 9, 37: 9, 38: 9, 39: 9, 40: 10, 41: 10, 42: 10, 43: 10, 44: 11, 45: 11, 46: 11, 47: 11, 48: 12, 49: 12, 50: 12, 51: 12, 52: 36, 53: 13, 54: 13, 55: 13, 56: 14, 57: 14, 58: 14, 59: 14, 60: 15, 61: 15, 62: 15, 63: 15, 64: 16, 65: 16, 66: 16, 67: 16, 68: 17, 69: 17, 70: 17, 71: 17, 72: 18, 73: 18, 74: 18, 75: 18, 76: 19, 77: 19, 78: 19, 79: 19, 80: 20, 81: 20, 82: 20, 83: 20, 84: 21, 85: 21, 86: 21, 87: 21, 88: 37, 89: 22, 90: 22, 91: 22, 92: 23, 93: 23, 94: 23, 95: 23, 96: 24, 97: 24, 98: 24, 99: 24, 100: 25, 101: 25, 102: 25, 103: 25, 104: 26, 105: 26, 106: 26, 107: 26, 108: 27, 109: 27, 110: 27, 111: 27, 112: 28, 113: 28, 114: 28, 115: 28, 116: 29, 117: 29, 118: 29, 119: 29, 120: 30, 121: 30, 122: 30, 123: 30, 124: 31, 125: 31, 126: 31, 127: 31, 128: 32, 129: 32, 130: 32, 131: 32, 132: 33, 133: 33, 134: 33, 135: 33, -1: -1}

const NumTiles int = 136

type Tiles []int

func (tiles Tiles) String() string {
	return ""
}

type Tile struct {
	tileId      int
	tileClass   int
	discardable bool
	isRinshan   bool
	isLast      bool
	discardWind Wind
}

func newTile(tileId int) *Tile {
	return &Tile{
		tileId:      tileId,
		tileClass:   tileId / 4,
		discardable: true,
		isRinshan:   false,
		isLast:      false,
		discardWind: -1}
}

type MahjongTiles struct {
	allTiles       map[int]*Tile
	tiles          [NumTiles]int
	kanNum         int
	NumRemainTiles int

	tilePointer    int
	rinshanPointer int
}

func NewMahjongTiles() *MahjongTiles {
	mahjongTiles := MahjongTiles{
		allTiles: make(map[int]*Tile, NumTiles),
	}
	mahjongTiles.Reset()
	return &mahjongTiles
}

func (tiles *MahjongTiles) Reset() {
	for i := 0; i < NumTiles; i++ {
		tiles.allTiles[i] = newTile(i)
		tiles.tiles[i] = i
	}

	rand.Shuffle(NumTiles, func(i, j int) {
		tiles.tiles[i], tiles.tiles[j] = tiles.tiles[j], tiles.tiles[i]
	})

	tiles.kanNum = 0
	tiles.NumRemainTiles = 70
	tiles.tilePointer = 52
	tiles.rinshanPointer = 135
}

func (tiles *MahjongTiles) Setup() (tonTiles Tiles, nanTiles Tiles, shaaTiles Tiles, peiTiles Tiles) {
	tonTiles = tiles.tiles[0:13]
	nanTiles = tiles.tiles[13:26]
	shaaTiles = tiles.tiles[26:39]
	peiTiles = tiles.tiles[39:52]
	return
}

func (tiles *MahjongTiles) DealTile(isRinshan bool) int {
	if tiles.NumRemainTiles <= 0 {
		panic(errors.New("no more tiles"))
	}
	if isRinshan && tiles.kanNum == 4 {
		panic(errors.New("no more rinshan tiles"))
	}
	tiles.NumRemainTiles--
	var tile int
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
		t = t.Append(tiles.tiles[130-2*i])
	}
	return t
}

func (tiles *MahjongTiles) UraDoraIndicators() Tiles {
	t := make(Tiles, 0, 5)
	for i := 0; i < tiles.kanNum+1; i++ {
		t = t.Append(tiles.tiles[131-2*i])
	}
	return t
}

func (tiles Tiles) Remove(tile int) Tiles {
	for i := 0; i < len(tiles); i++ {
		if tiles[i] == tile {
			tiles = append(tiles[:i], tiles[i+1:]...)
			return tiles
		}
	}
	panic("tile" + string(rune(tile)) + "not in tiles")
}

func (tiles Tiles) Append(tile int) Tiles {
	return append(tiles, tile)
}

func TilesEqual(tiles1 Tiles, tiles2 Tiles) bool {
	newArray1 := append([]int{}, tiles1...)
	newArray2 := append([]int{}, tiles2...)
	sort.Ints(newArray1)
	sort.Ints(newArray2)
	for i, tile := range newArray1 {
		if TileClassMap[newArray2[i]] != TileClassMap[tile] {
			return false
		}
	}
	return true
}

func (tiles Tiles) Copy() Tiles {
	tilesCopy := make(Tiles, len(tiles), cap(tiles))
	copy(tilesCopy, tiles)
	return tilesCopy
}

func (tiles Tiles) Index(tileId int, startIdx int) int {
	for i := startIdx; i < len(tiles); i++ {
		if tiles[i] == tileId {
			return i
		}
	}
	return -1
}

func (tiles Tiles) Count(tileId int) int {
	count := 0
	for _, tile := range tiles {
		if tile == tileId {
			count++
		}
	}
	return count
}
