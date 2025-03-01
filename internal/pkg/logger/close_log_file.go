package logger

import (
	"fmt"
	"os"
)

func CloseLog(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Printf("Error closing log file: %v", err)
	}
}
