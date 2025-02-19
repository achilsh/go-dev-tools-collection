package main

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestFunOne(t *testing.T) {
	// Test code here
	defer goleak.VerifyNone(t)
	// t.Run()
	go func() {
		fmt.Println("call demo...")
		time.Sleep(100*time.Hour)
	}()
}