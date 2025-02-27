package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"go.uber.org/dig"

	"dig_usage_demo/common"
)

func main() {
	//d := common.NewContainer()
	d := common.ExampleCall()
	var digGraphData bytes.Buffer
	if err := dig.Visualize(d, &digGraphData); err != nil {
		fmt.Println("get dig graph data fail, err: ", err)
		return
	}
	g, err := graphviz.New(context.Background())
	if err != nil {
		fmt.Println("new graphviz obj fail, err: ", err)
		return
	}
	defer g.Close()
	
	fmt.Println(digGraphData.String())
	graph, err := graphviz.ParseBytes(digGraphData.Bytes())
	if err != nil {
		fmt.Println("parse graph data to graphviz fail, err: ", err)
		return
	}
	
	graph.SetRankDir(cgraph.LRRank)
	err = g.RenderFilename(context.Background(), graph, graphviz.SVG, "demo.graph"+".svg")
	if err != nil {
		fmt.Println("write demo graph to file fail, err: ", err)
		return
	}
}
