# rpc_tutorial

On top of developing security analysis and prevention tools, the Prevoty Engineering team is always looking for better technologies to manage our cloud and on-premise service-oriented architectures. 

A package of the Go standard library that we use extensively is `net/rpc` - [this particular package](a href="http://golang.org/pkg/net/rpc/) simplifies the approach and LOC when it comes to developing your own RPC. If you're unfamiliar with Go, then you should probably [take the tour](http://tour.golang.org/). If the acronym _RPC_ sounds foreign, then perhaps you you should read [a little more](http://en.wikipedia.org/wiki/Remote_procedure_call). 

This project features a primitive key/value cache, complete with client and server implementations that communicate over TCP. Unit tests show how this overall system works. No external dependencies are required.
