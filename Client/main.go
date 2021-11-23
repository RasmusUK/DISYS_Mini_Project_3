package main

import (
	DISYS "DISYS_Mini_Project_3/gRPC"
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
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
var clientID string

func main() {
	clientID = uuid.New().String()
	log.Printf("Welcome to the Auction house")
	readInput()

}

func bid(amount int32) {
	conn1, err := grpc.Dial("localhost:8100", grpc.WithInsecure())
	if err != nil {
		log.Printf("client failed to connect to server")
	}
	defer conn1.Close()

	c1 := DISYS.NewBidAuctionClientFEClient(conn1)

	conn2, err := grpc.Dial("localhost:8200", grpc.WithInsecure())
	if err != nil {
		log.Printf("client failed to connect to server")
	}
	defer conn2.Close()

	c2 := DISYS.NewBidAuctionClientFEClient(conn2)

	conn3, err := grpc.Dial("localhost:8300", grpc.WithInsecure())
	if err != nil {
		log.Printf("client failed to connect to server")
	}
	defer conn3.Close()

	c3 := DISYS.NewBidAuctionClientFEClient(conn3)

	response1, err := c1.SendBidRequest(context.Background(), &DISYS.BidRequest{Amount: amount, RequestID: logicalClock, ClientID: clientID})
	if err != nil {
		log.Fatalf("Error when calling BidRequest: %s", err)
	}

	response2, err := c2.SendBidRequest(context.Background(), &DISYS.BidRequest{Amount: amount, RequestID: logicalClock, ClientID: clientID})
	if err != nil {
		log.Fatalf("Error when calling BidRequest: %s", err)
	}

	response3, err := c3.SendBidRequest(context.Background(), &DISYS.BidRequest{Amount: amount, RequestID: logicalClock, ClientID: clientID})
	if err != nil {
		log.Fatalf("Error when calling BidRequest: %s", err)
	}

	var agreement bool
	agreement = response1.Success == response2.Success && response1.Success == response3.Success
	log.Printf("Agreement: ", fmt.Sprint(agreement))

	var corrupt1 bool
	corrupt1 = response2.Success == response3.Success && response2.Success != response1.Success
	log.Printf("Corrupt1: ", fmt.Sprint(corrupt1))

	var corrupt2 bool
	corrupt2 = response1.Success == response3.Success && response1.Success != response2.Success
	log.Printf("Corrupt2: ", fmt.Sprint(corrupt2))

	var corrupt3 bool
	corrupt1 = response1.Success == response2.Success && response1.Success != response3.Success
	log.Printf("Corrupt3: ", fmt.Sprint(corrupt3))

	if agreement {
		if response1.Success {
			log.Printf("Bid was accepted")
		} else {
			log.Printf("Bid too low - try again")
		}
	} else if corrupt1 {
		if response2.Success {
			log.Printf("Bid was accepted")
		} else {
			log.Printf("Bid too low - try again")
		}
	} else if corrupt2 {
		if response1.Success {
			log.Printf("Bid was accepted")
		} else {
			log.Printf("Bid too low - try again")
		}
	} else if corrupt3 {
		if response1.Success {
			log.Printf("Bid was accepted")
		} else {
			log.Printf("Bid too low - try again")
		}
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

	resultReq, err := c.SendResultRequest(context.Background(), &DISYS.ResultRequest{RequestID: logicalClock, ClientID: clientID})
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
