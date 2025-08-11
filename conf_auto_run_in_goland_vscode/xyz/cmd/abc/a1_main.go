package main

import (
	"flag"
	"fmt"
)

var port = flag.Int("port", 8080, "server port")
var debug = flag.Bool("debug", false, "enable debug mode")

func main() {
	flag.Parse()
	fmt.Println("port: ", *port, "debug: ", *debug)

}
