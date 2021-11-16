package main

import (
	gRPC "DISYS_Mini_Project_3/gRPC"
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var ID string
var serverAddresses = make([]string, 0)
var requestNumber int32 = 0
var messageChannel = make(chan string, 1)

func main() {
	ID = uuid.New().String()
	findServerAddresses()
	fmt.Println("Addresses found:")
	for _, address := range serverAddresses {
		fmt.Println(address)
	}
	fmt.Println("Welcome to the action\nEnter an integer and press enter to make a bid")
	fmt.Println("Get the status of option by entering r and then press enter")
	readInputForever()
}

func findServerAddresses() {
	var wg sync.WaitGroup
	wg.Add(90)

	baseString := "localhost:80"
	for i := 10; i < 100; i++ {
		go pingServer(&wg, baseString+strconv.Itoa(i))
	}
	wg.Wait()
}

func pingServer(wg *sync.WaitGroup, address string) {
	defer wg.Done()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))

	if err != nil {
		return
	}

	defer conn.Close()

	c := gRPC.NewBidAuctionClientFEClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Ping(ctx, &gRPC.Empty{})

	if err == nil {
		serverAddresses = append(serverAddresses, address)
	}
}

func readInputForever() {

	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, _ := reader.ReadLine()
		if string(line) == "r" {
			sendResultRequestToAll()
		} else if number, err := strconv.Atoi(string(line)); err == nil {
			sendBidRequestToAll(int32(number))
		}
	}

}

func SendResultRequest(wg *sync.WaitGroup, address string) {
	defer wg.Done()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := gRPC.NewBidAuctionClientFEClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	response, err := c.SendResultRequest(ctx, &gRPC.ResultRequest{
		RequestID: requestNumber,
		ClientID:  ID,
	})

	if err != nil {
		sendToChannelIfNotFull("No bids have been made")
	}

	if response.Active {
		sendToChannelIfNotFull(response.Result)
	} else {
		sendToChannelIfNotFull("Auction is over:\n" + response.Result)
	}
}

func SendBidRequest(wg *sync.WaitGroup, address string, bid int32) {
	defer wg.Done()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := gRPC.NewBidAuctionClientFEClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	response, err := c.SendBidRequest(ctx, &gRPC.BidRequest{
		Amount:    bid,
		RequestID: requestNumber,
		ClientID:  ID,
	})

	if err != nil {
		sendToChannelIfNotFull("Action is closed")
	}

	if response.Success {
		sendToChannelIfNotFull("Successful bid")
	} else {
		sendToChannelIfNotFull("Your bid is too low")
	}
}

func sendResultRequestToAll() {
	requestNumber++
	var wg sync.WaitGroup
	wg.Add(len(serverAddresses))
	for _, serverAddress := range serverAddresses {
		go SendResultRequest(&wg, serverAddress)
	}
	wg.Wait()
	fmt.Println(<-messageChannel)
}

func sendBidRequestToAll(bid int32) {
	requestNumber++
	var wg sync.WaitGroup
	wg.Add(len(serverAddresses))
	for _, serverAddress := range serverAddresses {
		go SendBidRequest(&wg, serverAddress, bid)
	}
	wg.Wait()
	fmt.Println(<-messageChannel)
}

func sendToChannelIfNotFull(message string) {
	if len(messageChannel) != 1 {
		messageChannel <- message
	}
}
