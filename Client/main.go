package main

import (
	DISYS "DISYS_Mini_Project_3/gRPC"
	"bufio"
	"context"
	"fmt"
	"os"

	//"errors"
	"google.golang.org/grpc"
	"log"
	//"net"
	//"os"
	"strconv"
	//"time"
)

var logicalClock int32 = 1

func main() {

	log.Printf("Welcome to the Auction house")
	readInput()

}

func bid(amount int32) {
	conn, err := grpc.Dial("localhost:8100", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client failed to connect to server")
	}
	defer conn.Close()

	c := DISYS.NewBidAuctionClientFEClient(conn)

	response, err := c.SendBidRequest(context.Background(), &DISYS.BidRequest{Amount: amount, RequestID: logicalClock, ClientID: "My"})
	if err != nil {
		log.Fatalf("Error when calling BidRequest: %s", err)
	}

	if response.Success {
		log.Printf("Bid was accepted")
	} else {
		log.Printf("Bid too low - try again")
	}

	logicalClock++
}

func result() {
	conn, err := grpc.Dial("localhost:8100", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client failed to connect to server")
	}
	defer conn.Close()

	c := DISYS.NewBidAuctionClientFEClient(conn)

	resultReq, err := c.SendResultRequest(context.Background(), &DISYS.ResultRequest{RequestID: logicalClock, ClientID: "My"})
	if err != nil {
		log.Fatalf("Error when calling ResultRequest: %v", err)
	}

	if resultReq.Active {
		log.Printf("Highest bidder: %s", resultReq.Result)
	} else {
		log.Printf("Auction is over, the winner is: %s", resultReq.Result)
	}

	logicalClock++
}

func readInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := reader.ReadLine()
		if string(line) == "bid" {
			fmt.Println("Enter bid amount:")
			var amount string
			fmt.Scanln(&amount)
			temp, _ := strconv.ParseInt(amount, 10, 32)

			bid(int32(temp))
		} else if string(line) == "result" {
			result()
		} else {
			log.Printf("Invalid input.")
		}

	}
}
