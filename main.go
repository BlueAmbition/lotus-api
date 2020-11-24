package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
)

func main() {
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.wUt4O2OahC6-1ZgBVnSG2qTJO-zjWRv9pbvm7pa1qLM"
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	addr := "47.242.15.74:2345"

	var api apistruct.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	// Now you can call any API you're interested in.
	tipset, err := api.ChainHead(context.Background())
	if err != nil {
		log.Fatalf("calling chain head: %s", err)
	}
	fmt.Printf("Current chain head is: %s", tipset.String())
}
