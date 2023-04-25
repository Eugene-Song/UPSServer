package main

import (
	pb "UPSServer/pb"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	// Create a TCP server listening on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	connU, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}

	var wordlId int64
	for {
		decodedBytes, err := recvConn(connU)
		if err != nil {
			continue
		}
		// Unmarshal the UConnected message from the response
		uaStart := &pb.UAstart{}
		err = proto.Unmarshal(decodedBytes, uaStart)
		if err != nil {
			continue
		}
		log.Printf("UAStart: %v", uaStart)
		wordlId = *uaStart.Worldid
		break
	}

	// Construct the UConnect message
	var i int32 = 1
	var x int32 = 10
	var y int32 = 10
	wh := &pb.AInitWarehouse{
		Id: &i,
		X:  &x,
		Y:  &y,
	}

	isAmazon := true
	connect := &pb.AConnect{
		Worldid:  &wordlId,
		Initwh:   []*pb.AInitWarehouse{wh},
		IsAmazon: &isAmazon,
	}

	// Connect to the world
	connW, err := net.Dial("tcp", "vcm-32169.vm.duke.edu:23456")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	// Serialize the UConnect message
	marshaledBytes, err := proto.Marshal(connect)
	if err != nil {
		log.Fatalf("Failed to marshal UConnect message: %v", err)
	}
	decodedBytes := sendAndRecv(marshaledBytes, connW)

	// Unmarshal the UConnected message from the response
	connected := &pb.AConnected{}
	err = proto.Unmarshal(decodedBytes, connected)

	if err != nil {
		log.Fatalf("Failed to unmarshal UConnected message: %v", err)
	}

	log.Printf("Received UConnected message: %v", connected)
	log.Printf("World Id: %v", connected.Worldid)

	//// Connect to the ups
	//connU, err := net.Dial("tcp", "vcm-32169.vm.duke.edu:12345")
	//if err != nil {
	//	log.Fatalf("Failed to connect to server: %v", err)
	//}

	var seqNum int64 = 1
	var shipId int64 = 1
	var whId1 int32 = 1
	var x1 int32 = 10
	var y1 int32 = 10
	var destX int32 = 1
	var destY int32 = 1
	pickUpReq1 := &pb.AUPickupRequest{
		SeqNum:       &seqNum,
		ShipId:       &shipId,
		WarehouseId:  &whId1,
		X:            &x1,
		Y:            &y1,
		DestinationX: &destX,
		DestinationY: &destY,
	}

	auCommand := &pb.AUCommand{
		PickupRequests: []*pb.AUPickupRequest{pickUpReq1},
	}

	// Serialize the UConnect message
	auCommandBytes, err := proto.Marshal(auCommand)
	if err != nil {
		log.Fatalf("Failed to marshal UConnect message: %v", err)
	}

	// Encode the length of the message as a varint and prepend it to the message
	connectBytes := prefixVarintLength(auCommandBytes)

	log.Printf("Send pickup aucommand request")
	// Send the UConnect message
	_, err = connU.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UConnect message: %v", err)
	}
	log.Printf("Start recv from UPS")
	go recvUPS(connU, connW)

	for true {

	}
}

// used for recieve UAcommand from Amazon
func recvUPS(connU net.Conn, connW net.Conn) {
	// read from connection
	// unmarshal
	// send to World
	for {
		decodedBytes, err := recvConn(connU)
		if err != nil {
			continue
		}
		// Unmarshal the UConnected message from the response
		uaCommand := &pb.UACommand{}
		err = proto.Unmarshal(decodedBytes, uaCommand)
		if err != nil {
			continue
		}
		log.Printf("UACommand: %v", uaCommand)
	}
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
