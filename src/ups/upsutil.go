package ups

import (
	pb "UPSServer/pb"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

type PackageMetaData struct {
	// packageID or shipId
	packageId int64
	// destination coordinate x
	destX int32
	// destination coordinate y
	destY int32
	// warehouse id
	whID int32
	// truck id
	truckId int32
	// pickup coordinate x
	pickupX int32
	// pickup coordinate y
	pickupY int32
}

// used for encode after marshal
func prefixVarintLength(data []byte) []byte {
	messageLen := uint64(len(data))
	varintBytes := make([]byte, binary.MaxVarintLen64)
	varintLen := binary.PutUvarint(varintBytes, messageLen)
	return append(varintBytes[:varintLen], data...)
}

func (u *UPS) ConstructUCommandsPick(pickUpRequests []*pb.AUPickupRequest) *pb.UCommands {
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
			break
		}
	}
	u.TruckMutex.Unlock()
	log.Printf("Find a truck %d", truckId)
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

		u.MapTruckShipMutex.Lock()
		shipIds, ok := u.MapTruckShip[truckId]
		if ok {
			shipIds = append(shipIds, packageMeta.packageId)
		} else {
			u.MapTruckShip[truckId] = []int64{packageMeta.packageId}
		}
		u.MapTruckShipMutex.Unlock()
	}

	return ucommands
}

func (u *UPS) ConstructUCommandsDeliver(deliverRequests []*pb.AUDeliverRequest) *pb.UCommands {

	ucommands := &pb.UCommands{
		Deliveries: []*pb.UGoDeliver{},
		Simspeed:   &u.SimSpeed,
	}
	// get truckid
	var truckId int32
	shipId := *deliverRequests[0].ShipId
	packageMetaData := u.Package[shipId]
	truckId = packageMetaData.truckId

	// update truck status
	u.TruckMutex.Lock()
	u.Truck[truckId] = "delivering"
	u.TruckMutex.Unlock()

	// construct UCommands
	seqNum := randomInt64()
	uGoDeliver := &pb.UGoDeliver{
		Truckid:  &truckId,
		Packages: []*pb.UDeliveryLocation{},
		Seqnum:   &(seqNum),
	}
	for _, deliverRequest := range deliverRequests {
		shipID := *deliverRequest.ShipId
		packageData := u.Package[shipID]
		uDeliverLocation := &pb.UDeliveryLocation{
			Packageid: &shipID,
			X:         &packageData.destX,
			Y:         &packageData.destY,
		}
		uGoDeliver.Packages = append(uGoDeliver.Packages, uDeliverLocation)

	}
	u.UnAckedDeliver[seqNum] = uGoDeliver

	return ucommands
}

func sendAmazonACK(acks []int64, connA net.Conn) {
	uaCommand := &pb.UACommand{
		Acks: acks,
	}
	// while send request to world
	marshaledUCommands, _ := proto.Marshal(uaCommand)
	log.Printf("Sending Amazon ACK: %v", uaCommand)
	connectBytes := prefixVarintLength(marshaledUCommands)

	// Send the UConnect message
	_, err := connA.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UConnect message: %v", err)
	}
}

func sendWorldACK(acks []int64, connW net.Conn) {
	uCommands := &pb.UCommands{
		Acks: acks,
	}
	// while send request to world
	marshaledUCommands, _ := proto.Marshal(uCommands)
	connectBytes := prefixVarintLength(marshaledUCommands)

	// Send the UConnect message
	_, err := connW.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UConnect message: %v", err)
	}
	log.Printf("Sending world ACK: %v", uCommands)
}

func sendAmazonLoadReq(shipIds []int64, truckId int32, connA net.Conn) {
	uaCommand := &pb.UACommand{
		LoadRequests: []*pb.UALoadRequest{},
	}

	for _, shipId := range shipIds {
		seqNum := randomInt64()
		uaLoadRequest := &pb.UALoadRequest{
			SeqNum:  &seqNum,
			TruckId: &truckId,
			ShipId:  &shipId,
		}
		uaCommand.LoadRequests = append(uaCommand.GetLoadRequests(), uaLoadRequest)
	}
	marshaledUCommands, _ := proto.Marshal(uaCommand)
	connectBytes := prefixVarintLength(marshaledUCommands)
	log.Printf("Sending Amazon Load Request: %v", uaCommand)
	// Send the UConnect message
	_, err := connA.Write(connectBytes)
	if err != nil {
		log.Fatalf("Failed to send UACommand message: %v", err)
	}
}

// Helper functions
func randomInt64() int64 {
	var randomInt64 int64
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randomInt64)
	fmt.Println("Generated random int64:", randomInt64)
	return randomInt64
}
