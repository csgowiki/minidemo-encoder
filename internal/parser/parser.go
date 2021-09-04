package parser

import (
	"fmt"
	"os"

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
	attackTickMap := make(map[int][]events.WeaponFire)

	iParser.RegisterEventHandler(func(e events.FrameDone) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()

		// 解析WeaponAttack事件
		if attackEvent, ok := attackTickMap[currentTick]; ok {
			fmt.Println(attackEvent)
			delete(attackTickMap, currentTick)
		}

		
	})

	iParser.RegisterEventHandler(func(e events.WeaponFire) {
		gs := iParser.GameState()
		currentTick := gs.IngameTick()
		attackTickMap[currentTick] = append(attackTickMap[currentTick], e)
	})

	err = iParser.ParseToEnd()
	checkError(err)
}
