package main

import (
	pb "UPSServer/pb"
	ua "UPSServer/pb"
	"UPSServer/ups"
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// connect to mysql db
	db, err := sql.Open("mysql", "root:Wadqq3.23@tcp(localhost:3306)/upsdb")
	if err != nil {
		panic(err.Error())
	}
	// Create a channel to listen for signals
	signalChannel := make(chan os.Signal, 1)

	// Notify the signalChannel when an interrupt signal (Ctrl+C) is received
	signal.Notify(signalChannel, syscall.SIGINT)

	// connect to World
	numTruck := int32(50)
	connW, worldID := initConnectWorld(numTruck)
	// connect to Amazon, send worldId
	connA := initAmazon(worldID)

	trucks := make(map[int32]string, numTruck)
	for i := int32(0); i < numTruck; i++ {
		trucks[i] = "idle"
	}

	upsServer := &ups.UPS{
		SimSpeed:       100,
		Package:        make(map[int64]*ups.PackageMetaData),
		UnAckedPickup:  make(map[int64]*pb.UGoPickup),
		UnAckedDeliver: make(map[int64]*pb.UGoDeliver),
		Truck:          trucks,
		MapTruckShip:   make(map[int32][]int64),
		DB:             db,
	}

	go upsServer.LoopSendUnAcked(connW)
	go recvAmazon(connA, connW, upsServer)
	go recvWorld(connA, connW, upsServer)

	// Create a goroutine to wait for the interrupt signal
	go func() {
		<-signalChannel // Block until an interrupt signal is received
		upsServer.PrintUPS()
		// Perform any cleanup or additional actions before exiting
		fmt.Println("\nReceived Ctrl+C. Performing cleanup...")

		// Exit the program with a success status code (0)
		os.Exit(0)
	}()

	for true {
		uCommands := &pb.UCommands{
			Queries: []*pb.UQuery{},
		}
		for i := int32(0); i < numTruck; i++ {
			truckId := i
			seqNum := ups.RandomInt64()
			uQuery := &pb.UQuery{
				Truckid: &truckId,
				Seqnum:  &seqNum,
			}
			uCommands.Queries = append(uCommands.Queries, uQuery)
		}
		// while send request to world
		marshaledUCommands, _ := proto.Marshal(uCommands)
		connectBytes := prefixVarintLength(marshaledUCommands)

		log.Printf("Send UCommand Deliver Request to World")
		// Send the UConnect message
		_, err := connW.Write(connectBytes)
		if err != nil {
			log.Fatalf("Failed to send UConnect message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
}

func recvWorld(connA net.Conn, connW net.Conn, ups *ups.UPS) {
	for {
		decodedBytes, err := recvConn(connW)
		if err != nil {
			continue
		}
		// Unmarshal the UConnected message from the response
		uResponses := &pb.UResponses{}
		err = proto.Unmarshal(decodedBytes, uResponses)
		if err != nil {
			continue
		}

		// handle acked
		acks := uResponses.Acks
		if acks != nil {
			ups.DeleteAckedCommand(acks)
		}

		// 1. error field
		//errs := uResponses.Error

		//// 2. handle finished
		if uResponses.Finished != nil && *uResponses.Finished {
			connW.Close()
			connA.Close()
		}
		// 3. handle completions
		completions := uResponses.Completions
		if completions != nil {
			ups.HandleUFinished(completions, connA, connW)
		}
		// 4. handle delivered
		delivered := uResponses.Delivered
		if delivered != nil {
			ups.HandleUDeliverMade(delivered, connA, connW)
		}

		//5. handle truckstatus
		truckStatus := uResponses.Truckstatus
		if truckStatus != nil {
			ups.HandleTruckStatus(truckStatus, connW)
		}
	}
}

// used for recieve AUcommand from Amazon
func recvAmazon(connA net.Conn, connW net.Conn, ups *ups.UPS) {
	// read from connection
	// unmarshal
	// send to World
	for {
		decodedBytes, err := recvConn(connA)
		if err != nil {
			continue
		}
		// Unmarshal the UConnected message from the response
		auCommand := &pb.AUCommand{}
		err = proto.Unmarshal(decodedBytes, auCommand)
		if err != nil {
			continue
		}
		// 1. error field
		//errs := auCommand.Error
		// 2. handle pickup request
		pickups := auCommand.PickupRequests
		if pickups != nil {
			ups.HandlePickupRequest(pickups, connW, connA)
		}

		// 3. handle delivery request
		deliveries := auCommand.DeliverRequests
		if deliveries != nil {
			ups.HandleDeliverRequest(deliveries, connW, connA)
		}

		//deliveries := auCommand.DeliverRequests
	}
}

// used for init connection with simulation world
func initConnectWorld(numTruck int32) (net.Conn, int64) {
	// Construct the UConnect message
	isAmazon := false
	connect := &pb.UConnect{
		Trucks:   []*pb.UInitTruck{},
		IsAmazon: &isAmazon,
	}
	for i := int32(0); i < numTruck; i++ {
		id := i
		x := int32(0)
		y := int32(0)
		// Construct the Truck message
		truck := &pb.UInitTruck{
			Id: &id,
			X:  &x,
			Y:  &y,
		}
		connect.Trucks = append(connect.Trucks, truck)
	}

	log.Printf("Sending Trcuks message: %v", connect.Trucks)
	// Serialize the UConnect message
	marshaledConnect, err := proto.Marshal(connect)
	if err != nil {
		log.Fatalf("Failed to marshal UConnect message: %v", err)
	}

	// Connect to the server
	connW, err := net.Dial("tcp", "vcm-32169.vm.duke.edu:12345")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	decodedBytes := sendAndRecv(marshaledConnect, connW)

	// Unmarshal the UConnected message from the response
	connected := &pb.UConnected{}
	err = proto.Unmarshal(decodedBytes, connected)
	if err != nil {
		log.Fatalf("Failed to unmarshal UConnected message: %v", err)
	}
	log.Printf("Received UConnected message: %v", connected)
	worldId := *connected.Worldid

	return connW, worldId
}

// for sending worldId to Amazon
func initAmazon(worldId int64) net.Conn {
	seqNum := int64(1)
	connect := &ua.UAstart{
		Worldid: &worldId,
		Seqnum:  &seqNum,
	}
	// Serialize the UAstart message
	marshaledConnect, err := proto.Marshal(connect)
	if err != nil {
		log.Fatalf("Failed to marshal UAstart message: %v", err)
	}

	// Connect to the server
	// TODO: change to Amazon server
	connA, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to Amazon: %v", err)
	}

	// Encode the length of the message as a varint and prepend it to the message
	connectBytes := prefixVarintLength(marshaledConnect)

	// Send the UAstart message
	_, err = connA.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UAstart message: %v", err)
	}

	return connA
}

func sendAndRecv(serialized []byte, conn net.Conn) []byte {
	// Encode the length of the message as a varint and prepend it to the message
	connectBytes := prefixVarintLength(serialized)

	// Send the UConnect message
	_, err := conn.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UConnect message: %v", err)
	}

	// receive from connection
	decodedBytes, err := recvConn(conn)
	return decodedBytes
}

// used for encode after marshal
func prefixVarintLength(data []byte) []byte {
	messageLen := uint64(len(data))
	varintBytes := make([]byte, binary.MaxVarintLen64)
	varintLen := binary.PutUvarint(varintBytes, messageLen)
	return append(varintBytes[:varintLen], data...)
}

// decode length and read data from connection
func recvConn(conn net.Conn) ([]byte, error) {
	// Read the length of the data as a varint
	var messageLen uint64
	var bytesRead int
	var buf [binary.MaxVarintLen64]byte

	for messageLen == 0 {
		n, err := conn.Read(buf[bytesRead : bytesRead+1])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			continue
		}

		bytesRead += n
		messageLen, _ = binary.Uvarint(buf[:bytesRead])
	}

	// Read the data itself
	data := make([]byte, messageLen)
	_, err := io.ReadFull(conn, data)
	return data, err
}
