package ups

import (
	pb "UPSServer/pb"
	"database/sql"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"sync"
	"time"
)

type UPS struct {
	SimSpeed uint32

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

	// mapping between seqNum and shipId for pickUpCommand
	MapTruckShip      map[int32][]int64
	MapTruckShipMutex sync.Mutex

	DB *sql.DB
}

func (u *UPS) HandlePickupRequest(pickUpRequests []*pb.AUPickupRequest, connW net.Conn, connA net.Conn) {
	log.Printf("Enter HandlePickupRequest function")
	ucommands := u.ConstructUCommandsPick(pickUpRequests)
	log.Printf("Successfully construct pickup-UCommand: %v", ucommands)
	// while send request to world
	marshaledUCommands, _ := proto.Marshal(ucommands)
	connectBytes := prefixVarintLength(marshaledUCommands)

	log.Printf("Send UCommand Pickup Request to World")
	// Send the UConnect message
	_, err := connW.Write(connectBytes)
	if err != nil {
		log.Printf("Failed to send Pickup UCommand message: %v", err)
	}

	// Send Amazon ACKS
	acks := make([]int64, len(pickUpRequests))
	for i, each := range pickUpRequests {
		acks[i] = *each.SeqNum
	}

	sendAmazonACK(acks, connA)
}

func (u *UPS) HandleDeliverRequest(deliverRequests []*pb.AUDeliverRequest, connW net.Conn, connA net.Conn) {

	log.Printf("Enter HandleDeliverRequest function")
	ucommands := u.ConstructUCommandsDeliver(deliverRequests)
	log.Printf("Successfully construct Deliver-UCommands: %v", ucommands)
	// while send request to world
	marshaledUCommands, _ := proto.Marshal(ucommands)
	connectBytes := prefixVarintLength(marshaledUCommands)

	log.Printf("Send UCommand Deliver Request to World")
	// Send the UConnect message
	_, err := connW.Write(connectBytes)
	if err != nil {
		log.Printf("Failed to send Deliver UCommand message: %v", err)
	}

	// Send Amazon ACKS
	acks := make([]int64, len(deliverRequests))
	for i, each := range deliverRequests {
		acks[i] = *each.SeqNum
	}
	sendAmazonACK(acks, connA)
}

// function for handling UFinished from World
func (u *UPS) HandleUFinished(uFinishedResponses []*pb.UFinished, connA net.Conn, connW net.Conn) {
	log.Printf("Enter HandleUFinished function, uFinished responses: %v", uFinishedResponses)
	// update truck mapping
	// 1. get shipId that matches with truckId
	// 2. change truck status
	// 3. send load req to amazon
	// 4. send ACKs back to world
	shipIds := []int64{}
	u.MapTruckShipMutex.Lock()
	u.TruckMutex.Lock()
	for _, uFinishedResponse := range uFinishedResponses {
		truckId := *uFinishedResponse.Truckid
		// If the truck has finished all the deliveries and idle
		if *uFinishedResponse.Status == "idle" {
			log.Printf("Truck complete all deliveries: %v", truckId)
			u.Truck[truckId] = "idle"
			continue
		}

		// if the truck has associated packages on it
		if len(u.MapTruckShip[truckId]) != 0 {
			shipIds = u.MapTruckShip[truckId]
			u.MapTruckShip[truckId] = []int64{}
			u.Truck[truckId] = "arrive warehouse"

			// change status for all the packages or shipids
			u.updatePackLoadStatus(shipIds)

			sendAmazonLoadReq(shipIds, truckId, connA)
		} else {
			continue
		}
	}
	u.TruckMutex.Unlock()
	u.MapTruckShipMutex.Unlock()

	// Send World ACKS
	acks := make([]int64, len(uFinishedResponses))
	for i, each := range uFinishedResponses {
		acks[i] = *each.Seqnum
	}
	sendWorldACK(acks, connW)
}

func (u *UPS) HandleUDeliverMade(uDeliverMadeResponses []*pb.UDeliveryMade, connA net.Conn, connW net.Conn) {

	uaCommand := &pb.UACommand{
		Delivered: []*pb.UADelivered{},
	}
	for _, uDeliverMadeResponse := range uDeliverMadeResponses {
		seqNum := RandomInt64()
		shipId := *uDeliverMadeResponse.Packageid
		uaDelivered := &pb.UADelivered{
			SeqNum: &seqNum,
			ShipId: &shipId,
		}
		uaCommand.Delivered = append(uaCommand.Delivered, uaDelivered)

		// update package status to delivered
		u.Package[shipId].status = "delivered"
		u.Package[shipId].currX = u.Package[shipId].destX
		u.Package[shipId].currY = u.Package[shipId].destY
	}

	// send delivered to Amazon
	marshaledUCommands, _ := proto.Marshal(uaCommand)
	connectBytes := prefixVarintLength(marshaledUCommands)

	log.Printf("Sending Amazon UADelivered: %v", uaCommand)
	// Send the UConnect message
	_, err := connA.Write(connectBytes)
	if err != nil {
		log.Printf("Failed to send UADelivered-UACommand message: %v", err)
	}

	// Send World ACKS
	acks := make([]int64, len(uDeliverMadeResponses))
	for i, each := range uDeliverMadeResponses {
		acks[i] = *each.Seqnum
	}
	sendWorldACK(acks, connW)
}

// handle truck status queris
func (u *UPS) HandleTruckStatus(truckStatuses []*pb.UTruck, connW net.Conn) {

	for _, truckStatus := range truckStatuses {
		truckID := truckStatus.GetTruckid()
		truckX := truckStatus.GetX()
		truckY := truckStatus.GetY()

		u.PackageMutex.Lock()
		for _, v := range u.Package {
			if v.truckId == truckID && v.status == "out for delivery" {
				v.currX = truckX
				v.currY = truckY
			}
		}
		u.PackageMutex.Unlock()
	}

	// Send World ACKS
	acks := make([]int64, len(truckStatuses))
	for i, each := range truckStatuses {
		acks[i] = *each.Seqnum
	}
	sendWorldACK(acks, connW)
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
		log.Printf("One loop start send unacked!!!")
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
		log.Printf("send unacked, enter next loop!!!")
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()
}
