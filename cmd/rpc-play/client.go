package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/mkadit/go-toybox/internal/models"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer client.Close()

	// Test RPCDropMigration
	var reply string
	err = client.Call("DbAdapter.RPCDropMigration", models.RPCMigrationArgs{}, &reply)
	if err != nil {
		log.Fatalf("Error calling RPCDropMigration: %v", err)
	}
	fmt.Println("RPCDropMigration:", reply)

	// Test RPCMigrateDatabase
	// 	err = client.Call("DbAdapter.RPCMigrateDatabase", models.RPCMigrationArgs{}, &reply)
	// 	if err != nil {
	// 		log.Fatalf("Error calling RPCMigrateDatabase: %v", err)
	// 	}
	// 	fmt.Println("RPCMigrateDatabase:", reply)
}
