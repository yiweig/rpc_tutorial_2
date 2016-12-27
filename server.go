package main

import (
	"encoding/gob"
	"log"
	"net"
	"net/rpc"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	rpc.Register(CreateNewRPC())
	rpc.Register(CreateNewRPCForREST())

	l, e := net.Listen("tcp", ":9876")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	rpc.Accept(l)
}
