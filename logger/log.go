package logger

import (
	"encoding/json"
	"log"
	"os"
)

func Debug(obj interface{}) {
	f, err := os.OpenFile("/tmp/network-monitor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()

	log.SetOutput(f)
	jsonBytes, _ := json.MarshalIndent(&obj, "", " ")
	log.Println(string(jsonBytes))
}

func Println(msg string) {
	f, err := os.OpenFile("/tmp/network-monitor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(msg)
}

func Printf(msg string, v ...interface{}) {
	f, err := os.OpenFile("/tmp/network-monitor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf(msg, v)
}