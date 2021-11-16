package main

import (
	gRPC "DISYS_Mini_Project_3/gRPC"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

var highestBid int32 = 0
var highestBidder string
var biddersLogicalTime = make(map[string]int32)
var auctionIsActive bool = false

func main() {
	initServer()
	time.Sleep(time.Hour)
}

type server struct {
	gRPC.UnsafeBidAuctionClientFEServer
}

func (s server) Ping(ctx context.Context, empty *gRPC.Empty) (*gRPC.Empty, error) {
	return &gRPC.Empty{}, nil
}

func (s server) SendBidRequest(ctx context.Context, request *gRPC.BidRequest) (*gRPC.BidResponse, error) {
	if highestBid == 0 {
		auctionIsActive = true
		go auctionTime()
	}
	waitForYourTurn(request.ClientID, request.RequestID)
	biddersLogicalTime[request.ClientID] = request.RequestID

	if !auctionIsActive {
		return &gRPC.BidResponse{Success: false}, errors.New("auction is over")
	}

	if request.Amount > highestBid {
		highestBid = request.Amount
		highestBidder = request.ClientID
		return &gRPC.BidResponse{Success: true}, nil
	}
	return &gRPC.BidResponse{Success: false}, nil
}

func (s server) SendResultRequest(ctx context.Context, request *gRPC.ResultRequest) (*gRPC.ResultResponse, error) {
	waitForYourTurn(request.ClientID, request.RequestID)
	biddersLogicalTime[request.ClientID] = request.RequestID

	if highestBid == 0 {
		return nil, errors.New("no bids has been made")
	}

	name := highestBidder[:3]
	result := "Client " + name + " amount: " + strconv.Itoa(int(highestBid))

	return &gRPC.ResultResponse{
		Result: result,
		Active: auctionIsActive,
	}, nil
}

func initServer() {
	var lis net.Listener
	err := errors.New("not initiated yet")
	baseString := "localhost:80"
	counter := 10
	for err != nil && counter < 100 {
		connectionString := baseString + strconv.Itoa(counter)
		lis, err = net.Listen("tcp", connectionString)
		counter++
	}

	s := grpc.NewServer()
	gRPC.RegisterBidAuctionClientFEServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func auctionTime() {
	fmt.Println("Auction has started")
	time.Sleep(time.Minute * 1)
	auctionIsActive = false
	fmt.Println("Auction is done")
}

func waitForYourTurn(clientID string, requestID int32) {
	for biddersLogicalTime[clientID] != requestID-1 {
		time.Sleep(time.Millisecond * 500)
	}
}
