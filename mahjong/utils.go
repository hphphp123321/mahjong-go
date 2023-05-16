package mahjong

import (
	"github.com/dnovikoff/tempai-core/hand/shanten"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
)

func CalculateShantenNum(handTiles Tiles, melds Calls) int {
	instances, meldsOpt := TilesCallsToCalc(handTiles, melds)
	res := shanten.Calculate(instances, meldsOpt)
	return res.Total.Value
}

func GetTenpaiSlice(handTiles Tiles, melds Calls) TileClasses {
	var TenpaiSlice TileClasses

	instances, meldsOpt := TilesCallsToCalc(handTiles, melds)
	res := tempai.Calculate(instances, meldsOpt)
	tiles := tempai.GetWaits(res).Tiles()
	for _, t := range tiles {
		TenpaiSlice = append(TenpaiSlice, TileClass(int(t)-1))
	}
	return TenpaiSlice
}

func IndicatorsToDora(indicators Tiles) Tiles {
	var doraTiles Tiles
	for _, indicator := range indicators {
		doraTiles = append(doraTiles, IndicatorToDora(indicator))
	}
	return doraTiles
}

func IndicatorToDora(indicator Tile) Tile {
	switch indicator {
	case 32, 33, 34, 35, 68, 69, 70, 71, 104, 105, 106, 107:
		return indicator - 32
	case 120, 121, 122, 123:
		return indicator - 12
	case 132, 133, 134, 135:
		return indicator - 8
	default:
		return indicator + 1
	}
}

func GetYakuResult(handTiles Tiles, melds Calls, ctx *yaku.Context) *yaku.Result {
	instances, meldsOpt := TilesCallsToCalc(handTiles, melds)
	res := tempai.Calculate(instances, meldsOpt)
	if res == nil {
		return nil
	}
	yakuResult := yaku.Win(res, ctx, nil)
	return yakuResult
}

func GetScoreResult(scoreRule *score.RulesStruct, yakuResult *yaku.Result, honba int) score.Score {
	return score.GetScoreByResult(scoreRule, yakuResult, score.Honba(honba))
}

// GetTenpaiInfos
//
//	@Description: after one player get a tile, get the tenpai infos(such as how many tiles remaining to ron)
//	@param game
//	@param player
//	@return *TenpaiInfos
func GetTenpaiInfos(game *Game, player *Player) TenpaiInfos {
	tenpaiTiles := player.GetPossibleTenpaiTiles()
	if len(tenpaiTiles) == 0 {
		return nil
	}
	var tenpaiInfos = make(TenpaiInfos)
	for _, tileToDiscard := range tenpaiTiles {
		var playerCopy = player.Copy()
		playerCopy.HandTiles.Remove(tileToDiscard)
		tenpaiInfos[tileToDiscard] = GetTenpaiInfo(game, playerCopy)
	}
	return tenpaiInfos
}

// GetTenpaiInfo
//
//	@Description: get the tenpai info of a player after discard a tile
//	@param game
//	@param player
//	@return *TenpaiInfo
func GetTenpaiInfo(game *Game, player *Player) *TenpaiInfo {
	var tenpaiInfo = NewTenpaiInfo()
	tenpaiSlice := GetTenpaiSlice(player.HandTiles, player.Melds)
	if len(tenpaiSlice) == 0 {
		return nil
	}
	for _, tenpaiTileClass := range tenpaiSlice {
		tenpaiInfo.TileClassesTenpaiResult[tenpaiTileClass] = GetTenpaiResult(game, player, tenpaiTileClass)
		for _, tile := range player.DiscardTiles {
			if tile.Class() == tenpaiTileClass {
				tenpaiInfo.Furiten = true
			}
		}
		if player.IsFuriten() {
			tenpaiInfo.Furiten = true
		}
	}
	return tenpaiInfo
}

func GetTenpaiResult(game *Game, player *Player, tileClass TileClass) *TenpaiResult {
	remainNum := GetRemainTileClassNumFromPlayerPerspective(game, player, tileClass)
	var minWinTiles = tileClass.To4Tiles()
	for _, tile := range player.HandTiles {
		if tile.Class() == tileClass {
			minWinTiles.Remove(tile)
		}
	}
	var minWinTile Tile
	if len(minWinTiles) != 0 {
		minWinTile = minWinTiles[len(minWinTiles)-1]
	} else {
		minWinTile = tileClass.To4Tiles()[3]
	}

	result := game.getRonResult(player, minWinTile)
	return NewTenpaiResult(remainNum, result)
}

func GetRemainTileClassNumFromPlayerPerspective(game *Game, player *Player, tileClass TileClass) int {
	num := 4
	for _, tile := range player.HandTiles {
		if tile.Class() == tileClass {
			num--
		}
	}
	for _, p := range game.PosPlayer {
		for _, tile := range p.DiscardTiles {
			if tile.Class() == tileClass {
				num--
			}
		}
	}
	for _, tile := range game.Tiles.DoraIndicators() {
		if tile.Class() == tileClass {
			num--
		}
	}
	if num < 0 {
		panic("RemainTileClassNum < 0")
	}
	return num
}

func divideIntoLines(s string, lines int) []string {
	runes := []rune(s)
	length := len(runes)
	elementsPerLine := 12

	result := make([]string, lines)
	for i := 0; i < lines; i++ {
		line := ""
		for j := 0; j < elementsPerLine; j++ {
			index := i*elementsPerLine + j
			if index < length {
				line += string(runes[index])
			} else {
				line += " "
			}
		}
		result[i] = line
	}
	return result
}
