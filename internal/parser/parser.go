package parser

import (
	"os"

	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

type TickPlayer struct {
	tick int
	steamid uint64
}

func Start() {
	filePath := "./demofiles/faze-vs-vitality-m1-mirage.dem"
	iFile, err := os.Open(filePath)
	checkError(err)

	iParser := dem.NewParser(iFile)
	defer iParser.Close()

	// 处理特殊event构成的button表示
	var buttonTickMap map[TickPlayer]int32 = make(map[TickPlayer]int32)
	var (
		roundStarted      = 0
		roundInFreezetime = 0
		roundNum          = 0
	)

	iParser.RegisterEventHandler(func(e events.FrameDone) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()

		if roundInFreezetime == 0 {
			tPlayers := gs.TeamTerrorists().Members()
			ctPlayers := gs.TeamCounterTerrorists().Members()
			Players := append(tPlayers, ctPlayers...)
			for _, player := range Players {
				if player != nil {
					var addonButton int32 = 0
					key := TickPlayer{currentTick, player.SteamID64}
					if val, ok := buttonTickMap[key]; ok {
						addonButton = val
						delete(buttonTickMap, key)
					}
					parsePlayerFrame(player, addonButton, false)
				}
			}
		}
	})

	iParser.RegisterEventHandler(func(e events.WeaponFire) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()
		key := TickPlayer{currentTick, e.Shooter.SteamID64}
		if _, ok := buttonTickMap[key]; ok {
			buttonTickMap[key] |= IN_ATTACK
		} else {
			buttonTickMap[key] = IN_ATTACK
		}
	})

	iParser.RegisterEventHandler(func(e events.PlayerJump) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()
		key := TickPlayer{currentTick, e.Player.SteamID64}
		if _, ok := buttonTickMap[key]; ok {
			buttonTickMap[key] |= IN_JUMP
		} else {
			buttonTickMap[key] = IN_JUMP
		}
	})


	// 包括开局准备时间
	iParser.RegisterEventHandler(func(e events.RoundStart) {
		roundStarted = 1
		roundInFreezetime = 1
	})

	// 准备时间结束，正式开始
	iParser.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		roundInFreezetime = 0
		roundNum += 1
		ilog.InfoLogger.Println("回合开始：", roundNum)
		// 初始化录像文件
		// 写入所有选手的初始位置和角度
		gs := iParser.GameState()
		tPlayers := gs.TeamTerrorists().Members()
		ctPlayers := gs.TeamCounterTerrorists().Members()
		Players := append(tPlayers, ctPlayers...)
		for _, player := range Players {
			if player != nil {
				// parse player
				parsePlayerInitFrame(player)
			}
		}
	})

	// 回合结束，不包括自由活动时间
	iParser.RegisterEventHandler(func(e events.RoundEnd) {
		if roundStarted == 0 {
			roundStarted = 1
			roundNum = 0
		}
		ilog.InfoLogger.Println("回合结束：", roundNum)
		// 结束录像文件
		gs := iParser.GameState()
		tPlayers := gs.TeamTerrorists().Members()
		ctPlayers := gs.TeamCounterTerrorists().Members()
		Players := append(tPlayers, ctPlayers...)
		for _, player := range Players {
			if player != nil {
				// save to rec file
				saveToRecFile(player, int32(roundNum))
			}
		}
	})
	err = iParser.ParseToEnd()
	checkError(err)
}
