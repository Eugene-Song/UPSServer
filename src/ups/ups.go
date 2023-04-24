package ups

import (
	pb "UPSServer/pb"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"sync"
	"time"
)

type UPS struct {
	SimSpeed uint32

	SeqNum int64

	PackageMutex sync.Mutex
	// mapping between packageId and packageMetaData
	Package map[int64]*PackageMetaData

	// mapping between seqNum and NOT-ACKED pickUpCommand
	UnAckedPickup map[int64]*pb.UGoPickup

	// mapping between seqNum and NOT-ACKED deliverCommand
	UnAckedDeliver map[int64]*pb.UGoDeliver

	// mapping of trucks
	Truck      map[int32]string
	TruckMutex sync.Mutex
}

func (u *UPS) HandlePickupRequest(pickUpRequests []*pb.AUPickupRequest, connW net.Conn) {

	ucommands := &pb.UCommands{
		Pickups:  []*pb.UGoPickup{},
		Simspeed: &u.SimSpeed,
	}
	var truckId int32

	u.TruckMutex.Lock()
	for k, v := range u.Truck {
		if v == "idle" || v == "arrive warehouse" || v == "delivering" {
			//find a truck
			u.Truck[k] = "traveling"
			truckId = k
		}
	}
	u.TruckMutex.Unlock()

	// construct UCommands
	for _, pickUpRequest := range pickUpRequests {
		seqNum := randomInt64()
		uGoPickup := &pb.UGoPickup{
			Truckid: &truckId,
			Whid:    pickUpRequest.WarehouseId,
			Seqnum:  &(seqNum),
		}
		ucommands.Pickups = append(ucommands.Pickups, uGoPickup)
		u.UnAckedPickup[seqNum] = uGoPickup

		// fill Package mapping
		packageMeta := &PackageMetaData{
			packageId: *pickUpRequest.ShipId,
			destX:     *pickUpRequest.DestinationX,
			destY:     *pickUpRequest.DestinationY,
			whID:      *pickUpRequest.WarehouseId,
			truckId:   truckId,
			pickupX:   *pickUpRequest.X,
			pickupY:   *pickUpRequest.Y,
		}
		u.PackageMutex.Lock()
		u.Package[*pickUpRequest.ShipId] = packageMeta
		u.PackageMutex.Unlock()
	}
	// while send request to world

	marshaledUCommands, _ := proto.Marshal(ucommands)
	connectBytes := prefixVarintLength(marshaledUCommands)

	// Send the UConnect message
	_, err := connW.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UConnect message: %v", err)
	}
}

func (u *UPS) HandleUFinished(uFinishedResponse []*pb.UFinished, connA net.Conn) {
	// update truck mapping 
	
}

// function to delete acked command
func (u *UPS) DeleteAckedCommand(acks []int64) {
	for _, ack := range acks {
		_, keyExists := u.UnAckedPickup[ack]
		if keyExists {
			delete(u.UnAckedPickup, ack)
		}
		_, keyExists = u.UnAckedPickup[ack]
		if keyExists {
			delete(u.UnAckedDeliver, ack)
		}
	}
}

func (u *UPS) LoopSendUnAcked(conn net.Conn) {
	for true {
		ucommands := &pb.UCommands{
			Pickups:    []*pb.UGoPickup{},
			Deliveries: []*pb.UGoDeliver{},
			Simspeed:   &u.SimSpeed,
		}
		for _, v := range u.UnAckedPickup {
			ucommands.Pickups = append(ucommands.Pickups, v)
		}
		for _, v := range u.UnAckedDeliver {
			ucommands.Deliveries = append(ucommands.Deliveries, v)
		}
		if len(ucommands.Pickups) != 0 || len(ucommands.Deliveries) != 0 {
			marshaledUCommands, _ := proto.Marshal(ucommands)
			connectBytes := prefixVarintLength(marshaledUCommands)

			// Send the UConnect message
			_, err := conn.Write(connectBytes)
			if err != nil {
				log.Fatalf("Failed to send UConnect message: %v", err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func randomInt64() int64 {
	var randomInt64 int64
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randomInt64)
	fmt.Println("Generated random int64:", randomInt64)
	return randomInt64
}
