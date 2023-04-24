//package main
//
//import (
//	pb "UPSServer/src/pb"
//	"encoding/binary"
//	"errors"
//	"github.com/golang/protobuf/proto"
//	"log"
//	"net"
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
//	// Read the response
//	buf := make([]byte, 1024)
//	n, err := conn.Read(buf)
//	if err != nil {
//		log.Fatalf("Failed to read response: %v", err)
//	}
//
//	decodedBytes, err := decodeVarintPrefixed(buf[:n])
//	// Unmarshal the UConnected message from the response
//	connected := &pb.UConnected{}
//	err = proto.Unmarshal(decodedBytes, connected)
//
//	log.Printf(string(buf))
//	if err != nil {
//		log.Fatalf("Failed to unmarshal UConnected message: %v", err)
//	}
//
//	log.Printf("Received UConnected message: %v", connected)
//}
//
//// used for encode after marshal
//func prefixVarintLength(data []byte) []byte {
//	messageLen := uint64(len(data))
//	varintBytes := make([]byte, binary.MaxVarintLen64)
//	varintLen := binary.PutUvarint(varintBytes, messageLen)
//	return append(varintBytes[:varintLen], data...)
//}
//
//// used for decode before unmarshal
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

package main

import (
	pb "UPSServer/src/pb"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

func main() {
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
		Initwh:   []*pb.AInitWarehouse{wh},
		IsAmazon: &isAmazon,
	}

	// Serialize the UConnect message
	marshaledBytes, err := proto.Marshal(connect)
	if err != nil {
		log.Fatalf("Failed to marshal UConnect message: %v", err)
	}
	decodedBytes := sendAndRecv(marshaledBytes)

	// Unmarshal the UConnected message from the response
	connected := &pb.AConnected{}
	err = proto.Unmarshal(decodedBytes, connected)

	if err != nil {
		log.Fatalf("Failed to unmarshal UConnected message: %v", err)
	}

	log.Printf("Received UConnected message: %v", connected)
	log.Printf("World Id: %v", connected.Worldid)

}

func sendAndRecv(serialized []byte) []byte {
	// Encode the length of the message as a varint and prepend it to the message
	connectBytes := prefixVarintLength(serialized)

	// Connect to the server
	conn, err := net.Dial("tcp", "vcm-32169.vm.duke.edu:23456")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send the UConnect message
	_, err = conn.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UConnect message: %v", err)
	}

	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	decodedBytes, err := decodeVarintPrefixed(buf[:n])
	return decodedBytes
}

// used for encode after marshal
func prefixVarintLength(data []byte) []byte {
	messageLen := uint64(len(data))
	varintBytes := make([]byte, binary.MaxVarintLen64)
	varintLen := binary.PutUvarint(varintBytes, messageLen)
	return append(varintBytes[:varintLen], data...)
}

// used for decode before unmarshal
func decodeVarintPrefixed(data []byte) ([]byte, error) {
	messageLen, varintLen := binary.Uvarint(data)
	if varintLen <= 0 {
		return nil, errors.New("invalid varint encoding")
	}
	if uint64(len(data)) < messageLen+uint64(varintLen) {
		return nil, errors.New("insufficient data")
	}

	return data[varintLen : varintLen+int(messageLen)], nil
}
