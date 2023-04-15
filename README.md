# mahjong-go [![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](LICENSE)

English | [中文](README-CN.md)

mahjong-go is a Riichi Mahjong game logic package implemented in Go, based on the [tempai-core](https://github.com/dnovikoff/tempai-core) tenpai, winning hand, and other algorithms to implement the basic logic of four-player Riichi Mahjong. It includes:
1. Basic logic for a four-player Mahjong game, including dealing, chi, pon, kan, and winning hands
2. Customizable rules for four-player Mahjong, such as whether to allow exposed tiles with a terminal tile
3. JSON serialization and deserialization of game events
4. Saving and restoring a game session

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Documentation](#documentation)
- [License](#license)

## Installation
```bash
go get github.com/hphphp123321/mahjong-go
```

## Usage
Here is a simple example:
```go
package main

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"math/rand"
)

func main() {
	var seed int64 = 1
	players := make([]*mahjong.Player, 4)
	posCalls := make(map[mahjong.Wind]mahjong.Calls, 4)
	posCall := make(map[mahjong.Wind]*mahjong.Call, 4)
	for i := 0; i < 4; i++ {
		players[i] = mahjong.NewMahjongPlayer()
	}
	game := mahjong.NewMahjongGame(seed, nil)
	posCalls = game.Reset(players, nil)
	var flag = mahjong.EndTypeNone
	for flag != mahjong.EndTypeGame {
		for wind, calls := range posCalls {
			posCall[wind] = calls[rand.Intn(len(calls))]
		}
		posCalls, flag = game.Step(posCall)
		posCall = make(map[mahjong.Wind]*mahjong.Call, 4)
	}
}
```

## Examples
There are some examples using mahjong-go in the [tests](tests) directory as references.
- [game_test.go](tests/game_test.go): A simple example of using the package, including tests for single-game and multi-game sessions.
- [boardstate_test.go](tests/boardstate_test.go): An example of converting a certain game state into a BoardState structure, which can be used to quickly convert game states into all information visible to a player.
- [reconstruct_test.go](tests/reconstruct_test.go): An example of reconstructing game states by restoring all game events, which can be used for playback.
- [rule_test.go](tests/rule_test.go): An example of customizing Mahjong rules, which can be used for special Mahjong rules.

## Documentation
The documentation for the latest version of mahjong-go is available at [godoc.org](https://godoc.org/github.com/hphphp123321/mahjong-go).

## License
mahjong-go is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.





