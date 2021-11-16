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
	"sync"
	"time"
)

var highestBid int32 = 0
var highestBidder string
var biddersLogicalTime = make(map[string]int32)
var auctionIsActive bool = false

var biddersLogicalTimeLock sync.Mutex
var highestBidLock sync.Mutex
var highestBidderLock sync.Mutex

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
	if getHighestBid() == 0 {
		auctionIsActive = true
		go auctionTime()
	}
	waitForYourTurn(request.ClientID, request.RequestID)
	setBiddersLogicalTime(request.ClientID, request.RequestID)

	if !auctionIsActive {
		return &gRPC.BidResponse{Success: false}, errors.New("auction is over")
	}

	if setHighestBid(request.Amount) {
		setHighestBidder(request.ClientID)
		return &gRPC.BidResponse{Success: true}, nil
	}

	return &gRPC.BidResponse{Success: false}, nil
}

func (s server) SendResultRequest(ctx context.Context, request *gRPC.ResultRequest) (*gRPC.ResultResponse, error) {
	waitForYourTurn(request.ClientID, request.RequestID)
	setBiddersLogicalTime(request.ClientID, request.RequestID)

	if getHighestBid() == 0 {
		return nil, errors.New("no bids has been made")
	}

	name := getHighestBidder()[:3]
	result := "Client " + name + " amount: " + strconv.Itoa(int(getHighestBid()))

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
	for {
		if getBiddersLogicalTime(clientID) == requestID-1 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func setBiddersLogicalTime(clientID string, requestID int32) {
	biddersLogicalTimeLock.Lock()
	defer biddersLogicalTimeLock.Unlock()
	biddersLogicalTime[clientID] = requestID
}

func getBiddersLogicalTime(clientID string) (requestID int32) {
	biddersLogicalTimeLock.Lock()
	defer biddersLogicalTimeLock.Unlock()
	return biddersLogicalTime[clientID]
}

func getHighestBid() (highest int32) {
	highestBidLock.Lock()
	defer highestBidLock.Unlock()
	return highestBid
}

func setHighestBid(input int32) (success bool) {
	highestBidLock.Lock()
	defer highestBidLock.Unlock()
	if input > highestBid {
		highestBid = input
		return true
	}
	return false
}

func getHighestBidder() (name string) {
	highestBidLock.Lock()
	defer highestBidLock.Unlock()
	return highestBidder
}

func setHighestBidder(name string) {
	highestBidderLock.Lock()
	defer highestBidderLock.Unlock()
	highestBidder = name
}
