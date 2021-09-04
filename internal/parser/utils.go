package parser

import (
	ilog "github.com/hx-w/minidemo-encoder/internal/logger"
)

// Function to handle errors
func checkError(err error) {
	if err != nil {
		ilog.ErrorLogger.Println(err.Error())
	}
}
