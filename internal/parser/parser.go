package parser

import (
	"os"

	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

func Start() {
	filePath := "./demofiles/faze-vs-vitality-m1-mirage.dem"
	iFile, err := os.Open(filePath)
	checkError(err)

	iParser := dem.NewParser(iFile)
	defer iParser.Close()

	// 用来记录某一Tick下WeaponAttack事件，在FrameDone中处理
	var attackTickMap map[int][]events.WeaponFire = make(map[int][]events.WeaponFire)
	// flags
	var (
		roundStarted      = 0
		roundInFreezetime = 0
		roundNum          = 0
		currentFrameIdx   = 0
	)

	iParser.RegisterEventHandler(func(e events.FrameDone) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()

		if (roundInFreezetime == 0) && (currentFrameIdx == 0) {
			tPlayers := gs.TeamTerrorists().Members()
			ctPlayers := gs.TeamCounterTerrorists().Members()
			Players := append(tPlayers, ctPlayers...)
			for _, player := range Players {
				if player != nil {
					// 解析WeaponAttack事件
					var isAttacking bool = false
					if attackEvent, ok := attackTickMap[currentTick]; ok {
						for _, atEvent := range attackEvent {
							if atEvent.Shooter.SteamID64 == player.SteamID64 {
								isAttacking = true
								break
							}
						}
					}
					parsePlayerFrame(player, isAttacking)
				}
			}
			delete(attackTickMap, currentTick)
		} else {
			if currentFrameIdx == 0 {
				currentFrameIdx = 0
			} else {
				currentFrameIdx = currentFrameIdx + 1
			}
		}

	})

	iParser.RegisterEventHandler(func(e events.WeaponFire) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()
		attackTickMap[currentTick] = append(attackTickMap[currentTick], e)
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

	// 正式结束，包括自由活动时间
	iParser.RegisterEventHandler(func(e events.RoundEndOfficial) {
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

	// 回合结束，不包括自由活动时间
	iParser.RegisterEventHandler(func(e events.RoundEnd) {
		if roundStarted == 0 {
			roundStarted = 1
			roundNum = 0
		}
	})
	err = iParser.ParseToEnd()
	checkError(err)
}
