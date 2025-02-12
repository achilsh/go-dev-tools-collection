package main

import (
	"fmt"

	"github.com/google/uuid"
)

func callGoogleUUID() {
	id := uuid.New()
	fmt.Println(id.ID())
	//对于基于时间的 UUID，可以提取时间戳：
	id , err := uuid.NewUUID()
	if err != nil {
		fmt.Println("call new uuid fail, err:", err)
		return 
	}
	fmt.Println("uuid: ", id, ", node: ", id.NodeID(), ", timeseq: ", id.ClockSequence())
	//
	srcId := "0fe8df98-e91c-11ef-bcd9-00155de954ed"
	dstId, err := uuid.Parse(srcId)
	if err != nil {
		fmt.Println("parse fail, err: ", err)
		return 
	}
	fmt.Println(dstId.String())
	fmt.Println(dstId.ID())
}
