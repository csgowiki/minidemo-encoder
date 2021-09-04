package encoder

import (
	"bytes"
	"os"
	"time"

	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
	iparser "github.com/hx-w/minidemo-encoder/internal/parser"
)

var __MAGIC__ int32 = -559038737
var __FORMAT_VERSION__ int8 = 2
var __FIELDS_ORIGIN__ int32 = 1 << 0
var __FIELDS_ANGLES__ int32 = 1 << 1
var __FIELDS_VELOCITY__ int32 = 1 << 2

var bufMap map[string]*bytes.Buffer = make(map[string]*bytes.Buffer)

func init() {
	saveDir := "./output"
	if ok, _ := PathExists(saveDir); !ok {
		os.Mkdir(saveDir, os.ModePerm)
		ilog.InfoLogger.Println("未找到保存目录，已创建：", saveDir)
	} else {
		ilog.InfoLogger.Println("保存目录存在：", saveDir)
	}
}

func InitPlayerRecFile(initFrame iparser.FrameInitInfo) {
	if bufMap[initFrame.PlayerName] == nil {
		bufMap[initFrame.PlayerName] = new(bytes.Buffer)
	} else {
		bufMap[initFrame.PlayerName].Reset()
	}
	// step.1 MAGIC NUMBER
	WriteToBuf(initFrame.PlayerName, __MAGIC__)

	// step.2 VERSION
	WriteToBuf(initFrame.PlayerName, __FORMAT_VERSION__)

	// step.3 timestamp
	WriteToBuf(initFrame.PlayerName, int32(time.Now().Unix()))

	// step.4 name length
	WriteToBuf(initFrame.PlayerName, int8(len(initFrame.PlayerName)))

	// step.5 name
	WriteToBuf(initFrame.PlayerName, []byte(initFrame.PlayerName))

	// step.6 initial position
	for idx := 0; idx < 3; idx++ {
		WriteToBuf(initFrame.PlayerName, float32(initFrame.Position[idx]))
	}

	// step.7 initial angle
	for idx := 0; idx < 2; idx++ {
		WriteToBuf(initFrame.PlayerName, initFrame.Angles[idx])
	}
	ilog.InfoLogger.Println("初始化成功: ", initFrame.PlayerName)
}

func WriteToRecFile(playerName string, roundNum int32) {
	fileName := "./output/" + string(roundNum) + "_" + playerName + ".rec"
	file, err := os.Create(fileName) // 创建文件, "binbin"是文件名字
	if err != nil {
		ilog.ErrorLogger.Println("文件创建失败", err.Error())
		return
	}
	defer file.Close()

	// step.8 tick count
	var tickCount int32 = int32(len(iparser.PlayerFramesMap[playerName]))
	WriteToBuf(playerName, tickCount)

	// step.9 bookmark count
	WriteToBuf(playerName, int32(0))

	// step.10 all bookmark
	// ignore

	// step.11 all tick frame
	for _, frame := range iparser.PlayerFramesMap[playerName] {
		WriteToBuf(playerName, frame.PlayerButtons)
		WriteToBuf(playerName, frame.PlayerImpulse)
		for idx := 0; idx < 3; idx++ {
			WriteToBuf(playerName, frame.ActualVelocity[idx])
		}
		for idx := 0; idx < 3; idx++ {
			WriteToBuf(playerName, frame.PredictedVelocity[idx])
		}
		for idx := 0; idx < 2; idx++ {
			WriteToBuf(playerName, frame.PredictedAngles[idx])
		}
		WriteToBuf(playerName, frame.CSWeaponID)
		WriteToBuf(playerName, frame.PlayerSubtype)
		WriteToBuf(playerName, frame.PlayerSeed)
		WriteToBuf(playerName, frame.AdditionalFields)
		// 附加信息
		if frame.AdditionalFields|__FIELDS_ORIGIN__ != 0 {
			for idx := 0; idx < 3; idx++ {
				WriteToBuf(playerName, frame.AtOrigin[idx])
			}
		}
		if frame.AdditionalFields|__FIELDS_ANGLES__ != 0 {
			for idx := 0; idx < 3; idx++ {
				WriteToBuf(playerName, frame.AtAngles[idx])
			}
		}
		if frame.AdditionalFields|__FIELDS_VELOCITY__ != 0 {
			for idx := 0; idx < 3; idx++ {
				WriteToBuf(playerName, frame.AtVelocity[idx])
			}
		}
	}

	file.Write(bufMap[playerName].Bytes())
	ilog.InfoLogger.Printf("[%d]初始化成功: %s\n", roundNum, playerName)
}
