package ups

import "encoding/binary"

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
