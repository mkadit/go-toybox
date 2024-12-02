package rpc_server

import (
	"fmt"
	"net"
	"net/rpc"

	logfile "github.com/mkadit/go-toybox/internal/logger"
	"github.com/mkadit/go-toybox/internal/models"
)

func Start(config models.RPCConfiguration) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	defer listener.Close()

	logfile.LogEvent("RPC server listening on port 1234")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}

}
