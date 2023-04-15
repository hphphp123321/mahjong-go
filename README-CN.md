# mahjong-go [![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](LICENSE)

[English](README.md) | 中文

mahjong-go是一个用Go实现的立直麻将游戏逻辑包，基于[tempai-core](https://github.com/dnovikoff/tempai-core)的听牌、和牌等算法实现了四人立直麻将的基本逻辑。其中包括：
1. 4人麻将的基本逻辑，包括发牌、吃、碰、杠、和牌等
2. 自定义4人麻将规则，包括断幺是否可以有副露等
3. 游戏事件的json序列化与反序列化
4. 一局游戏的保存与恢复

## 目录
- [安装](#安装)
- [使用](#使用)
- [示例](#示例)
- [文档](#文档)
- [许可证](#许可证)

## 安装
```bash
go get github.com/hphphp123321/mahjong-go
```

## 使用
下面是一个简单的使用例子
```go
package main

import (
	mahjong "github.com/hphphp123321/mahjong-go"
)

func main() {
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

## 示例
在tests目录下有一些使用mahjong-go的例子可以作为参考。
- [game_test.go](tests/game_test.go): 简单的使用例子，其中包括一局游戏和多局游戏的测试。
- [boardstate_test.go](tests/boardstate_test.go): 把游戏的某个状态给转变为[BoardState](https://github.com/hphphp123321/mahjong-go/mahjong/boardstate.go)结构体的例子，可以用于将游戏状态快速转变为某个玩家能够看到的所有信息。
- [reconstruct_test.go](tests/reconstruct_test.go): 通过将游戏的所有事件进行复原，还原出游戏的状态的例子，可以用于录像回放。
- [rule_test.go](tests/rule_test.go): 自定义麻将规则的例子，可以用于一些特殊的麻将规则。


## 文档
- [go文档](https://pkg.go.dev/github.com/hphphp123321/mahjong-go)

## 许可证
mahjong-go基于Apache V2许可证发布，详情请参考[LICENSE](LICENSE)文件。