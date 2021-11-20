package main

import (
	DISYS "DISYS_Mini_Project_3/gRPC"
	"context"
	//"errors"
	"google.golang.org/grpc"
	"log"
	//"net"
	//"os"
	//"strconv"
	//"time"
)

func main() {
	conn, err := grpc.Dial("localhost:8100", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client failed to connect to server")
	}
	defer conn.Close()

	c := DISYS.NewBidAuctionClientFEClient(conn)
	log.Printf("Client is now connected to server")

	response, err := c.SendBidRequest(context.Background(), &DISYS.BidRequest{Amount: 1, RequestID: 1, ClientID: "My"})
	if err != nil {
		log.Fatalf("Error when calling BidRequest: %s", err)
	}
	log.Printf("test if we get here")
	log.Printf("Response from server: %s", response.Success)

	resultReq, err := c.SendResultRequest(context.Background(), &DISYS.ResultRequest{RequestID: 2, ClientID: "My"})
	if err != nil {
		log.Fatalf("Error when calling ResultRequest: %v", err)
	}

	log.Printf("Result response from server: %s", resultReq.Result)

}

func bid() {
	panic("implement me")
}

func result() {
	panic("implement me")
}
