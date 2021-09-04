package parser

import (
	"fmt"
)

// Function to handle errors
func checkError(err error) {
	if err != nil {
		fmt.Println("ERROR", err)
		// ErrorLogger.Println("DEMO STREAM ERROR")
		// WarningLogger.Println("Demo stream errors can still write output, check for JSON file")
		// ErrorLogger.Println(err.Error())
	}
}
