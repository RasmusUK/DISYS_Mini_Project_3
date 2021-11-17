package main

import (
	"DISYS_Mini_Project_3/gRPC"
	"bufio"
	"context"
	"errors"
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
	fmt.Println("Your username is: ", ID[:3])
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
	log.Println("Pinging server:", address)
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection")
		}
	}(conn)

	c := gRPC.NewBidAuctionClientFEClient(conn)

	_, err = c.Ping(ctx, &gRPC.Empty{})

	if err == nil {
		serverAddresses = append(serverAddresses, address)
	}

	log.Println("Received ping response from:", address)
}

func readInputForever() {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, _ := reader.ReadLine()
		if string(line) == "r" {
			sendResultRequestToAll()
		} else if number, err := strconv.Atoi(string(line)); err == nil {
			sendBidRequestToAll(int32(number))
		} else {
			fmt.Println("Invalid input")
		}
	}

}

func SendResultRequest(wg *sync.WaitGroup, address string) {
	log.Println("Send result request to:", address)
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(address, "did not respond - deleting from server addresses")
		removeAddressFromAddresses(address)
		return
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection")
		}
	}(conn)

	c := gRPC.NewBidAuctionClientFEClient(conn)

	response, err := c.SendResultRequest(ctx, &gRPC.ResultRequest{
		RequestID: requestNumber,
		ClientID:  ID,
	})

	log.Println("Received result response from:", address)

	if err != nil {
		sendToChannelIfNotFull("No bids have been made")
	} else if response.Active {
		sendToChannelIfNotFull(response.Result)
	} else {
		sendToChannelIfNotFull("Auction is over:\n" + response.Result)
	}
}

func SendBidRequest(wg *sync.WaitGroup, address string, bid int32) {
	log.Println("Send bid request to:", address)

	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println(address, "did not respond - deleting from server addresses")
		removeAddressFromAddresses(address)
		return
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection")
		}
	}(conn)

	c := gRPC.NewBidAuctionClientFEClient(conn)

	if err != nil {
		removeAddressFromAddresses(address)
		return
	}

	response, err := c.SendBidRequest(ctx, &gRPC.BidRequest{
		Amount:    bid,
		RequestID: requestNumber,
		ClientID:  ID,
	})

	log.Println("Received bid response from:", address)

	if err != nil {
		sendToChannelIfNotFull("Action is closed")
	} else if response.Success {
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

func removeAddressFromAddresses(address string) {
	index, err := findIndexOfAddress(address)
	if err != nil {
		return
	}
	serverAddresses = append(serverAddresses[:index], serverAddresses[index+1:]...)
}

func findIndexOfAddress(address string) (int, error) {
	for i := 0; i < len(serverAddresses); i++ {
		if serverAddresses[i] == address {
			return i, nil
		}
	}
	return 0, errors.New("could not find address")
}
