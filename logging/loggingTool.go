package main

import (
	"fmt"
)

// func loggingTool() {
// 	// Open a log file
// 	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	// Set the log output to the file
// 	log.SetOutput(file)

// 	// Use the log package to write log messages
// 	log.Println("This is a log message")
// 	log.Printf("This is a formatted log message with a value: %d", 42)
// }

func LoggingTool() {
	fmt.Println("TEST")
}
