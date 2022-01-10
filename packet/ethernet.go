package packet

import (
	"encoding/hex"
	"fmt"
)

var IPv4 [2]byte = [2]byte{0x08, 0x00}
var IPv6 [2]byte = [2]byte{0x86, 0xdd}

type EthernetIIFrame struct {
	DestinationMac [6]byte
	SourceMac      [6]byte
	EtherType      [2]byte
	DataIPv4       IPv4Packet
	DataIPv6       IPv6Packet

	PacketType int

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
	ipType := DataPacketType(dataArr[14:])
	ipTypeInt := int(ipType)
	dataIPv4 := IPv4Packet{}
	dataIPv6 := IPv6Packet{}
	if ipTypeInt == 4 {
		dataIPv4 = DataToIPv4Packet(dataArr[14:])
	} else if ipTypeInt == 6 {
		dataIPv6 = DataToIPv6Packet(dataArr[14:])
	}

	return EthernetIIFrame{
		DestinationMac: destMac,
		SourceMac:      srcMac,
		EtherType:      etherType,
		DataIPv4:       dataIPv4,
		DataIPv6:       dataIPv6,
		PacketType:     ipTypeInt,
	}
}

func (frame EthernetIIFrame) GetPacketType() int {
	return frame.PacketType
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
	hexString += frame.DataIPv4.IPv4ToHexString()
	hexString += frame.DataIPv6.IPv6ToHexString()

	return hexString
}

func (frame EthernetIIFrame) ToString() string {
	var returnStr string
	returnStr += "Ethernet II: {"

	returnStr += "\n    Destination MAC: "
	returnStr += insertHexFormat(frame.DestinationMac[:], ":")

	returnStr += "\n    Source MAC: "
	returnStr += insertHexFormat(frame.SourceMac[:], ":")

	returnStr += "\n    Type: 0x"
	returnStr += insertHexFormat(frame.EtherType[:], "")
	if frame.EtherType == IPv4 {
		returnStr += " (IPv4)"
	} else if frame.EtherType == IPv6 {
		returnStr += " (IPv6)"
	}

	returnStr += "\n\n    IPv4: {"
	returnStr += "\n" + frame.DataIPv4.ToString()
	returnStr += "\n    }"

	returnStr += "\n}\n"
	return returnStr
}

func insertHexFormat(byteArr []byte, delimiter string) string {
	var returnStr string
	for i := 0; i < len(byteArr); i++ {
		returnStr += hex.EncodeToString([]byte{byteArr[i]})
		if i+1 != len(byteArr) {
			returnStr += delimiter
		}
	}
	return returnStr
}

func insertDecimalFormat(byteArr []byte, delimiter string) string {
	var returnStr string
	for i := 0; i < len(byteArr); i++ {
		returnStr += fmt.Sprint(byteArr[i])
		if i+1 != len(byteArr) {
			returnStr += "."
		}
	}

	return returnStr
}
