package main

import (
	pb "UPSServer/pb"
	ua "UPSServer/pb"
	"UPSServer/ups"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

//func main() {
//	// Construct the UConnect message
//	var i int32 = 1
//	var x int32 = 0
//	var y int32 = 0
//	truck := &pb.UInitTruck{
//		Id: &i,
//		X:  &x,
//		Y:  &y,
//	}
//
//	var i2 int32 = 2
//	var x2 int32 = 0
//	var y2 int32 = 0
//	truck2 := &pb.UInitTruck{
//		Id: &i2,
//		X:  &x2,
//		Y:  &y2,
//	}
//
//	isAmazon := false
//	connect := &pb.UConnect{
//		Trucks:   []*pb.UInitTruck{truck, truck2},
//		IsAmazon: &isAmazon,
//	}
//
//	// Serialize the UConnect message
//	marshaledBytes, err := proto.Marshal(connect)
//	if err != nil {
//		log.Fatalf("Failed to marshal UConnect message: %v", err)
//	}
//
//	// Encode the length of the message as a varint and prepend it to the message
//	connectBytes := prefixVarintLength(marshaledBytes)
//
//	// Connect to the server
//	conn, err := net.Dial("tcp", "vcm-32169.vm.duke.edu:12345")
//	if err != nil {
//		log.Fatalf("Failed to connect to server: %v", err)
//	}
//	defer conn.Close()
//
//	// Send the UConnect message
//	_, err = conn.Write(connectBytes)
//	if err != nil {
//		log.Fatalf("Failed to send UConnect message: %v", err)
//	}
//
//	_, err = conn.Write(connectBytes)
//	if err != nil {
//		log.Fatalf("Failed to send UConnect message: %v", err)
//	}
//
//	fmt.Println("Starting...")
//	time.Sleep(5 * time.Second) // Sleep for 3 seconds
//	fmt.Println("...Finished")
//
//	//// Read the response
//	//buf := make([]byte, 1024)
//	//n, err := conn.Read(buf)
//	//if err != nil {
//	//	log.Fatalf("Failed to read response: %v", err)
//	//}
//
//	decodedBytes, err := decodeVarintPrefixed(conn)
//	// Unmarshal the UConnected message from the response
//	connected := &pb.UConnected{}
//	err = proto.Unmarshal(decodedBytes, connected)
//
//	log.Printf("Received UConnected message: %v", connected)
//	//log.Printf(string(buf))
//	if err != nil {
//		log.Fatalf("Failed to unmarshal UConnected message: %v", err)
//	}
//}
//
//// used for encode after marshal
//func prefixVarintLength(data []byte) []byte {
//	messageLen := uint64(len(data))
//	varintBytes := make([]byte, binary.MaxVarintLen64)
//	varintLen := binary.PutUvarint(varintBytes, messageLen)
//	return append(varintBytes[:varintLen], data...)
//}

// used for decode before unmarshal
//func decodeVarintPrefixed(data []byte) ([]byte, error) {
//	messageLen, varintLen := binary.Uvarint(data)
//	if varintLen <= 0 {
//		return nil, errors.New("invalid varint encoding")
//	}
//	if uint64(len(data)) < messageLen+uint64(varintLen) {
//		return nil, errors.New("insufficient data")
//	}
//
//	return data[varintLen : varintLen+int(messageLen)], nil
//}

//func decodeVarintPrefixed(conn net.Conn) ([]byte, error) {
//	// Read the length of the data as a varint
//	var messageLen uint64
//	var bytesRead int
//	var buf [binary.MaxVarintLen64]byte
//
//	for messageLen == 0 {
//		n, err := conn.Read(buf[bytesRead : bytesRead+1])
//		if err != nil {
//			return nil, err
//		}
//		if n == 0 {
//			continue
//		}
//
//		bytesRead += n
//		messageLen, _ = binary.Uvarint(buf[:bytesRead])
//	}
//
//	// Read the data itself
//	data := make([]byte, messageLen)
//	_, err := io.ReadFull(conn, data)
//	return data, err
//}

//package main
//
//import (
//	pb "UPSServer/pb"
//	ua "UPSServer/pb"
//	"encoding/binary"
//	"errors"
//	"log"
//	"net"
//
//	"github.com/golang/protobuf/proto"
//)

//
//func main() {
//	// Construct the UConnect message
//	var i int32 = 1
//	var x int32 = 0
//	var y int32 = 0
//	truck := &pb.UInitTruck{
//		Id: &i,
//		X:  &x,
//		Y:  &y,
//	}
//
//	var i2 int32 = 2
//	var x2 int32 = 0
//	var y2 int32 = 0
//	truck2 := &pb.UInitTruck{
//		Id: &i2,
//		X:  &x2,
//		Y:  &y2,
//	}
//
//	worldId := int64(1)
//	isAmazon := false
//	connect := &pb.UConnect{
//		Worldid:  &worldId,
//		Trucks:   []*pb.UInitTruck{truck, truck2},
//		IsAmazon: &isAmazon,
//	}
//
//	// Serialize the UConnect message
//	marshaledBytess, err := proto.Marshal(connect)
//	if err != nil {
//		log.Fatalf("Failed to marshal UConnect message: %v", err)
//	}
//
//	// Connect to the server
//	conn, err := net.Dial("tcp", "vcm-32169.vm.duke.edu:12345")
//	if err != nil {
//		log.Fatalf("Failed to connect to server: %v", err)
//	}
//	defer conn.Close()
//
//	decodedBytes := sendAndRecv(marshaledBytess, conn)
//
//	// Unmarshal the UConnected message from the response
//	connected := &pb.UConnected{}
//	err = proto.Unmarshal(decodedBytes, connected)
//
//	if err != nil {
//		log.Fatalf("Failed to unmarshal UConnected message: %v", err)
//	}
//
//	log.Printf("Received UConnected message: %v", connected)
//
//	// go pick up
//	var whid int32 = 1
//	var seqNum int64 = 1
//	var truckId int32 = 1
//	pickup := &pb.UGoPickup{
//		Truckid: &truckId,
//		Whid:    &whid,
//		Seqnum:  &seqNum,
//	}
//
//	var seqNum2 int64 = 2
//	var truckId2 int32 = 2
//	pickup2 := &pb.UGoPickup{
//		Truckid: &truckId2,
//		Whid:    &whid,
//		Seqnum:  &seqNum2,
//	}
//
//	command := &pb.UCommands{
//		Pickups: []*pb.UGoPickup{pickup, pickup2},
//	}
//
//	// Serialize the UCommand message
//	marshaledBytes, err := proto.Marshal(command)
//	if err != nil {
//		log.Fatalf("Failed to marshal UConnect message: %v", err)
//	}
//	decodedBytes = sendAndRecv(marshaledBytes, conn)
//	commandACK := &pb.UResponses{}
//	err = proto.Unmarshal(decodedBytes, commandACK)
//
//	if err != nil {
//		log.Fatalf("Failed to unmarshal UCommands ACK: %v", err)
//	}
//
//	log.Printf("Received UConnected message: %v", commandACK)
//
//}

func main() {
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
