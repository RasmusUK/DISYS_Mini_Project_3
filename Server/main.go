package main

import (
	gRPC "DISYS_Mini_Project_3/gRPC"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var highestBid int32 = 0
var highestBidder string
var biddersLogicalTime map[string]int32
var auctionIsActive bool = false

func main() {
	input := readArgs()
	initServer(input)
}

func readArgs() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	} else {
		return "localhost:8100"
	}
}

type server struct {
	gRPC.UnsafeBidAuctionClientFEServer
}

func (s server) SendBidRequest(ctx context.Context, request *gRPC.BidRequest) (*gRPC.BidResponse, error) {
	if highestBid == 0 {
		auctionIsActive = true
		go auctionTime()
	}

	waitForYourTurn(request.ClientID, request.RequestID)

	if !auctionIsActive {
		return &gRPC.BidResponse{Success: false}, errors.New("auction is over")
	}

	biddersLogicalTime[request.ClientID] = request.RequestID

	if request.Amount > highestBid {
		highestBid = request.Amount
		highestBidder = request.ClientID
		return &gRPC.BidResponse{Success: true}, nil
	}
	return &gRPC.BidResponse{Success: false}, nil
}

func (s server) SendResultRequest(ctx context.Context, request *gRPC.ResultRequest) (*gRPC.ResultResponse, error) {
	waitForYourTurn(request.ClientID, request.RequestID)

	if highestBid == 0 {
		return nil, errors.New("no bids has been made")
	}

	index := getIndexOfBidder(highestBidder)
	result := "Client " + strconv.Itoa(index) + "amount: " + strconv.Itoa(int(highestBid))

	return &gRPC.ResultResponse{
		Result: result,
		Active: auctionIsActive,
	}, nil
}

func initServer(connectionString string) {
	lis, err := net.Listen("tcp", connectionString)
	if err != nil {
		log.Fatalf("failed to listen: %v\nPlease try another port", err)
	}
	s := grpc.NewServer()
	gRPC.RegisterBidAuctionClientFEServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func auctionTime() {
	time.Sleep(time.Minute * 2)
	auctionIsActive = false
}

func waitForYourTurn(clientID string, requestID int32) {
	for biddersLogicalTime[clientID] != requestID-1 {
		//wait for sequential consistency
	}
}

func getIndexOfBidder(bidderID string) int {
	index := 1
	for i, _ := range biddersLogicalTime {
		if i == bidderID {
			return index
		}
		index++
	}
	return 0
}
