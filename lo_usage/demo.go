package main

import (
	"log"

	"github.com/samber/lo"
)

func demoSlices() {
	srcList := []string{"aa", "bb", "cc", "aa"}
	dst := lo.Uniq(srcList)
	log.Printf("data: %v", dst)
}

func main() {
	demoSlices()

}
