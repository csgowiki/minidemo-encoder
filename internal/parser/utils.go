package parser

import (
	encoder "github.com/hx-w/minidemo-encoder/internal/encoder"
	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

// Function to handle errors
func checkError(err error) {
	if err != nil {
		ilog.ErrorLogger.Println(err.Error())
	}
}

func parsePlayerInitFrame(player *common.Player) {
	iFrameInit := encoder.FrameInitInfo{
		PlayerName: player.Name,
	}
	iFrameInit.Position[0] = float32(player.Position().X)
	iFrameInit.Position[1] = float32(player.Position().Y)
	iFrameInit.Position[2] = float32(player.Position().Z)
	// 注意XY，需要测试
	iFrameInit.Angles[0] = float32(player.ViewDirectionY())
	iFrameInit.Angles[1] = float32(player.ViewDirectionX())

	encoder.InitPlayer(iFrameInit)
}


func parsePlayerFrame(player *common.Player, isAttack bool) {
	if !player.IsAlive() {
		return
	}
	iFrameInfo := new(encoder.FrameInfo)
	iFrameInfo.PredictedVelocity[0] = float32(player.Velocity().X)
	iFrameInfo.PredictedVelocity[1] = float32(player.Velocity().Y)
	iFrameInfo.PredictedVelocity[2] = float32(player.Velocity().Z)
	iFrameInfo.ActualVelocity[0] = float32(player.Velocity().X)
	iFrameInfo.ActualVelocity[1] = float32(player.Velocity().Y)
	iFrameInfo.ActualVelocity[2] = float32(player.Velocity().Z)
	iFrameInfo.PredictedAngles[0] = player.ViewDirectionY()
	iFrameInfo.PredictedAngles[1] = player.ViewDirectionX()
	// ----- button encode
	iFrameInfo.PlayerButtons = 0
	iFrameInfo.PlayerImpulse = 0
	iFrameInfo.PlayerSeed = 0
	iFrameInfo.PlayerSubtype = 0
	if len(encoder.PlayerFramesMap[player.Name]) == 0 {
		iFrameInfo.CSWeaponID = 2
	} else {
		iFrameInfo.CSWeaponID = 0 // glock
	}
	iFrameInfo.AdditionalFields |= encoder.FIELDS_ORIGIN
	iFrameInfo.AtOrigin[0] = float32(player.Position().X)
	iFrameInfo.AtOrigin[1] = float32(player.Position().Y)
	iFrameInfo.AtOrigin[2] = float32(player.Position().Z)
	encoder.PlayerFramesMap[player.Name] = append(encoder.PlayerFramesMap[player.Name], *iFrameInfo)
}

func saveToRecFile(player *common.Player, roundNum int32) {
	encoder.WriteToRecFile(player.Name, roundNum)
}
