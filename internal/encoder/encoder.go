package encoder

import (
	"os"

	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
)

func init() {
	saveDir := "./output"
	if ok, _ := PathExists(saveDir); !ok {
		os.Mkdir(saveDir, os.ModePerm)
		ilog.InfoLogger.Println("未找到保存目录，已创建：", saveDir)
	} else {
		ilog.InfoLogger.Println("保存目录存在：", saveDir)
	}
}
