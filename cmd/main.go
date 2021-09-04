package main

import (
	encoder "github.com/hx-w/minidemo-encoder/internal/encoder"
	iparser "github.com/hx-w/minidemo-encoder/internal/parser"
)

func main() {
	frameInit := iparser.FrameInitInfo{
		PlayerName: "apEX",
	}
	frameInit.Position[0] = 1.0
	frameInit.Position[1] = 1.0
	frameInit.Position[2] = 1.0
	frameInit.Angles[0] = 1.0
	frameInit.Angles[1] = 1.0
	encoder.InitPlayerRecFile(frameInit)
	// miniParser.Start()
}
