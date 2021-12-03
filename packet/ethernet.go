package packet

import "encoding/hex"

var IPv4 [2]byte = [2]byte{0x08, 0x00}
var IPv6 [2]byte = [2]byte{0x08, 0x06}

type EthernetIIFrame struct {
	DestinationMac [6]byte
	SourceMac      [6]byte
	EtherType      [2]byte
	Data           IPv4Packet
	//CRCChecksum    [4]byte
}

func DataToFrame(dataArr []byte) EthernetIIFrame { // TODO: Look at encoding/gib
	var destMac [6]byte
	for i := 0; i < 6; i++ {
		destMac[i] = dataArr[i]
	}
	var srcMac [6]byte
	for i := 0; i < 6; i++ {
		srcMac[i] = dataArr[6+i]
	}
	etherType := [2]byte{dataArr[12], dataArr[13]}
	data := DataToPacket(dataArr[14:])

	return EthernetIIFrame{
		DestinationMac: destMac,
		SourceMac:      srcMac,
		EtherType:      etherType,
		Data:           data,
	}
}

// func NewEthernetFrame(destMac [6]byte, srcMac [6]byte,
// 	ethType [2]byte, data []byte /*checksum [4]byte*/) EthernetIIFrame {

// 	frame := EthernetIIFrame{destMac, srcMac, ethType, data} //, checksum}
// 	return frame
// }

// func EthernetFrameFromString(frame string) EthernetIIFrame {
// 	return EthernetIIFrame{}
// }

func (frame EthernetIIFrame) ToHexString() string {
	var hexString string
	hexString += hex.EncodeToString(frame.DestinationMac[:])
	hexString += hex.EncodeToString(frame.SourceMac[:])
	hexString += hex.EncodeToString(frame.EtherType[:])
	hexString += frame.Data.ToHexString()

	return hexString
}
