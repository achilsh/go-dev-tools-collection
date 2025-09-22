package main

import (
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	log.Println("This is a test log")
	log.Println("This is a test log")
	log.Printf("This is a test log with format: %s", "1111")
	log.Printf("This is a test log with format: %s", "2222")
	log.Fatalln("This is a fatal log")
}
