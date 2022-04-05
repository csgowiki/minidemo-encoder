package parser

import (
	"math"

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

const Pi = 3.14159265358979323846

func radian2degree(radian float64) float64 {
	return radian * 180 / Pi
}

func normalizeDegree(degree float64) float64 {
	if degree < 0.0 {
		degree = degree + 360.0
	}
	return degree
}

var (
	lastX float32
	lastY float32
	lastZ float32
)

func parsePlayerFrame(player *common.Player, addonButton int32, tickrate float64, fullsnap bool) {
	if !player.IsAlive() {
		return
	}
	iFrameInfo := new(encoder.FrameInfo)
	iFrameInfo.PredictedVelocity[0] = 0.0
	iFrameInfo.PredictedVelocity[1] = 0.0
	iFrameInfo.PredictedVelocity[2] = 0.0
	iFrameInfo.ActualVelocity[0] = float32(player.Velocity().X)
	iFrameInfo.ActualVelocity[1] = float32(player.Velocity().Y)
	iFrameInfo.ActualVelocity[2] = float32(player.Velocity().Z)
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
	if fullsnap {
		iFrameInfo.AdditionalFields |= encoder.FIELDS_ORIGIN
		iFrameInfo.AtOrigin[0] = float32(player.Position().X)
		iFrameInfo.AtOrigin[1] = float32(player.Position().Y)
		iFrameInfo.AtOrigin[2] = float32(player.Position().Z)
	}
	// test
	if player.Name == "b1t" {
		deltaX := float32(player.Position().X) - lastX
		deltaY := float32(player.Position().Y) - lastY
		deltaZ := float32(player.Position().Z) - lastZ
		lastX = float32(player.Position().X)
		lastY = float32(player.Position().Y)
		lastZ = float32(player.Position().Z)

		// velocity in Z direction need to be recorded specially
		iFrameInfo.ActualVelocity[0] = deltaX * float32(tickrate)
		iFrameInfo.ActualVelocity[1] = deltaY * float32(tickrate)
		iFrameInfo.ActualVelocity[2] = deltaZ * float32(tickrate)

		// Since I don't know how to get player's button bits in a tick frame,
		// I have to use *actual vels* and *angles* to generate *predict vels* approximately
		// This will cause some error, but it's not a big deal
		var velAngle float64 = 0.0
		if iFrameInfo.ActualVelocity[0] == 0.0 {
			if iFrameInfo.ActualVelocity[1] < 0.0 {
				velAngle = 270.0
			} else {
				velAngle = 90.0
			}
		} else {
			velAngle = radian2degree(math.Atan2(float64(iFrameInfo.ActualVelocity[1]), float64(iFrameInfo.ActualVelocity[0])))
			velAngle = normalizeDegree(velAngle)
		}
		faceFront := normalizeDegree(float64(iFrameInfo.PredictedAngles[1]))
		deltaAngle := normalizeDegree(velAngle - faceFront)
		const threshold = 30.0
		if 0.0+threshold < deltaAngle && deltaAngle < 180.0-threshold {
			iFrameInfo.PredictedVelocity[1] = -450.0 // left
		}
		if 90.0+threshold < deltaAngle && deltaAngle < 270.0-threshold {
			iFrameInfo.PredictedVelocity[0] = -450.0 // back
		}
		if 180.0+threshold < deltaAngle && deltaAngle < 360.0-threshold {
			iFrameInfo.PredictedVelocity[1] = 450.0 // right
		}
		if 270.0+threshold < deltaAngle || deltaAngle < 90.0-threshold {
			iFrameInfo.PredictedVelocity[0] = 450.0 // front
		}
	}

	encoder.PlayerFramesMap[player.Name] = append(encoder.PlayerFramesMap[player.Name], *iFrameInfo)
}

func saveToRecFile(player *common.Player, roundNum int32) {
	encoder.WriteToRecFile(player.Name, roundNum)
}
