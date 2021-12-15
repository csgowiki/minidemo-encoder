package parser

import (
	encoder "github.com/hx-w/minidemo-encoder/internal/encoder"
	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

var bufWeaponMap map[string]int32 = make(map[string]int32)

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
	iFrameInit.Angles[0] = float32(player.ViewDirectionY())
	iFrameInit.Angles[1] = float32(player.ViewDirectionX())

	encoder.InitPlayer(iFrameInit)
	delete(bufWeaponMap, player.Name)
}

func parsePlayerFrame(player *common.Player, addonButton int32) {
	if !player.IsAlive() {
		return
	}
	iFrameInfo := new(encoder.FrameInfo)
	iFrameInfo.PredictedVelocity[0] = 0.0
	iFrameInfo.PredictedVelocity[1] = 0.0
	iFrameInfo.PredictedVelocity[2] = 0.0
	iFrameInfo.ActualVelocity[0] = float32(player.Velocity().X)
	iFrameInfo.ActualVelocity[1] = float32(player.Velocity().Y)
	iFrameInfo.ActualVelocity[2] = float32(-250) // debug
	iFrameInfo.PredictedAngles[0] = player.ViewDirectionY()
	iFrameInfo.PredictedAngles[1] = player.ViewDirectionX()
	iFrameInfo.PlayerImpulse = 0
	iFrameInfo.PlayerSeed = 0
	iFrameInfo.PlayerSubtype = 0
	// ----- button encode
	iFrameInfo.PlayerButtons = ButtonConvert(player, addonButton)

	// ---- weapon encode
	var currWeaponID int32 = int32(WeaponStr2ID(player.ActiveWeapon().String()))
	if len(encoder.PlayerFramesMap[player.Name]) == 0 {
		iFrameInfo.CSWeaponID = currWeaponID
		bufWeaponMap[player.Name] = currWeaponID
	} else if currWeaponID == bufWeaponMap[player.Name] {
		iFrameInfo.CSWeaponID = int32(CSWeapon_NONE)
	} else {
		iFrameInfo.CSWeaponID = currWeaponID
		bufWeaponMap[player.Name] = currWeaponID
	}

	// 附加项
	// iFrameInfo.AdditionalFields |= encoder.FIELDS_ORIGIN
	// iFrameInfo.AtOrigin[0] = float32(player.Position().X)
	// iFrameInfo.AtOrigin[1] = float32(player.Position().Y)
	// iFrameInfo.AtOrigin[2] = float32(player.Position().Z)
	// iFrameInfo.AdditionalFields |= encoder.FIELDS_ANGLES
	// iFrameInfo.AtAngles[0] = float32(player.ViewDirectionY())
	// iFrameInfo.AtAngles[1] = float32(player.ViewDirectionX())
	// iFrameInfo.AtAngles[2] = 0
	// iFrameInfo.AdditionalFields |= encoder.FIELDS_VELOCITY
	// iFrameInfo.AtVelocity[0] = float32(player.Velocity().X)
	// iFrameInfo.AtVelocity[1] = float32(player.Velocity().Y)
	// iFrameInfo.AtVelocity[2] = float32(player.Velocity().Z)
	encoder.PlayerFramesMap[player.Name] = append(encoder.PlayerFramesMap[player.Name], *iFrameInfo)
}

func saveToRecFile(player *common.Player, roundNum int32) {
	encoder.WriteToRecFile(player.Name, roundNum)
}
