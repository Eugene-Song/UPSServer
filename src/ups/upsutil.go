package ups

import (
	pb "UPSServer/pb"
	"crypto/rand"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
)

type PackageMetaData struct {
	// packageID or shipId
	PackageId int64
	// destination coordinate x
	DestX int32
	// destination coordinate y
	DestY int32
	// warehouse id
	whID int32
	// truck id
	TruckId int32
	// pickup coordinate x
	pickupX int32
	// pickup coordinate y
	pickupY int32
	// username
	username string
	// Status
	Status string
	// current X
	currX int32
	// current Y
	currY int32
	// item details
	itemDetails string
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

	u.PackageMutex.Lock()
	defer u.PackageMutex.Unlock()
	log.Printf("Successfully acquire PackageMutex lock")
	for _, pickUpRequest := range pickUpRequests {
		seqNum := RandomInt64()
		uGoPickup := &pb.UGoPickup{
			Truckid: &truckId,
			Whid:    pickUpRequest.WarehouseId,
			Seqnum:  &(seqNum),
		}
		ucommands.Pickups = append(ucommands.Pickups, uGoPickup)
		u.UnAckedPickupMutex.Lock()
		u.UnAckedPickup[seqNum] = uGoPickup
		u.UnAckedPickupMutex.Unlock()

		// fill Package mapping
		packageMeta := &PackageMetaData{
			PackageId:   pickUpRequest.GetShipId(),
			DestX:       pickUpRequest.GetDestinationX(),
			DestY:       pickUpRequest.GetDestinationY(),
			whID:        pickUpRequest.GetWarehouseId(),
			TruckId:     truckId,
			pickupX:     pickUpRequest.GetX(),
			pickupY:     pickUpRequest.GetY(),
			username:    pickUpRequest.GetUpsName(),
			Status:      "truck en route to warehouse",
			currX:       pickUpRequest.GetX(),
			currY:       pickUpRequest.GetY(),
			itemDetails: pickUpRequest.GetItems(),
		}
		u.Package[*pickUpRequest.ShipId] = packageMeta

		u.updatePackageTable(packageMeta)
		log.Printf("Successfully update package table in ConstructUCommandsPick.")

		u.MapTruckShipMutex.Lock()
		shipIds, ok := u.MapTruckShip[truckId]
		if ok {
			shipIds = append(shipIds, packageMeta.PackageId)
		} else {
			u.MapTruckShip[truckId] = []int64{packageMeta.PackageId}
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
	truckId = packageMetaData.TruckId

	// update truck status
	u.TruckMutex.Lock()
	u.Truck[truckId] = "delivering"
	u.TruckMutex.Unlock()

	// construct UCommands
	seqNum := RandomInt64()
	uGoDeliver := &pb.UGoDeliver{
		Truckid:  &truckId,
		Packages: []*pb.UDeliveryLocation{},
		Seqnum:   &(seqNum),
	}

	u.PackageMutex.Lock()
	for _, deliverRequest := range deliverRequests {
		shipID := *deliverRequest.ShipId
		packageData := u.Package[shipID]

		uDeliverLocation := &pb.UDeliveryLocation{
			Packageid: &shipID,
			X:         &packageData.DestX,
			Y:         &packageData.DestY,
		}
		uGoDeliver.Packages = append(uGoDeliver.Packages, uDeliverLocation)

		// update package status
		packageData.Status = "out for delivery"
		u.updatePackageTable(packageData)
	}
	u.PackageMutex.Unlock()

	u.UnAckedDeliverMutex.Lock()
	u.UnAckedDeliver[seqNum] = uGoDeliver
	u.UnAckedDeliverMutex.Unlock()

	return ucommands
}

func SendAmazonACK(acks []int64, connA net.Conn) {
	uaCommand := &pb.UACommand{
		Acks: acks,
	}
	// while send request to world
	marshaledUCommands, _ := proto.Marshal(uaCommand)
	connectBytes := prefixVarintLength(marshaledUCommands)

	log.Printf("Sending Amazon ACK: %v", uaCommand)
	// Send the UConnect message
	_, err := connA.Write(connectBytes)
	if err != nil {
		log.Printf("Failed to send ACK message to Amazon: %v", err)
	}
}

func sendWorldACK(acks []int64, connW net.Conn) {
	uCommands := &pb.UCommands{
		Acks: acks,
	}
	// while send request to world
	marshaledUCommands, _ := proto.Marshal(uCommands)
	connectBytes := prefixVarintLength(marshaledUCommands)

	//log.Printf("Sending world ACK: %v", uCommands)
	log.Printf("Sending world ACK!")
	// Send the UConnect message
	_, err := connW.Write(connectBytes)
	if err != nil {
		log.Printf("Failed to send ACK message to World: %v", err)
	}
}

func (u *UPS) sendAmazonLoadReq(shipIds []int64, truckId int32, connA net.Conn) {
	uaCommand := &pb.UACommand{
		LoadRequests: []*pb.UALoadRequest{},
	}

	for _, shipId := range shipIds {
		seqNum := RandomInt64()
		uaLoadRequest := &pb.UALoadRequest{
			SeqNum:  &seqNum,
			TruckId: &truckId,
			ShipId:  &shipId,
		}
		u.UnAckedLoadMutex.Lock()
		u.UnAckedLoad[seqNum] = uaLoadRequest
		u.UnAckedLoadMutex.Unlock()

		uaCommand.LoadRequests = append(uaCommand.GetLoadRequests(), uaLoadRequest)
	}
	marshaledUCommands, _ := proto.Marshal(uaCommand)
	connectBytes := prefixVarintLength(marshaledUCommands)
	log.Printf("Sending Amazon UACommand Load Request: %v", uaCommand)
	// Send the UConnect message
	_, err := connA.Write(connectBytes)
	if err != nil {
		log.Printf("Failed to send UACommand Load Request message: %v", err)
	}
}

// Helper functions
func RandomInt64() int64 {
	var randomInt64 int64
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randomInt64)
	return randomInt64
}

func (u *UPS) updatePackLoadStatus(shipIds []int64) {
	u.PackageMutex.Lock()
	for _, shipId := range shipIds {
		packageMeta := u.Package[shipId]
		packageMeta.Status = "truck waiting for package"
		u.updatePackageTable(packageMeta)
	}
	u.PackageMutex.Unlock()

}
