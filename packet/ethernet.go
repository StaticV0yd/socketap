package packet

import (
	"encoding/hex"
	"fmt"
)

type EtherType [2]byte
type MacAddress [6]byte

/*
Variable storing the byte values that designate IPv4
in the EtherType field.
*/
var IPv4 EtherType = [2]byte{0x08, 0x00}

/*
Variable storing the byte values that designate IPv6
in the EtherType field.
*/
var IPv6 EtherType = [2]byte{0x86, 0xdd}

/*
Struct that helps describe Ethernet II frames for both
looking at incoming ethernet frames and for constructing
new/custom ethernet frames.
*/
type EthernetIIFrame struct {
	DestinationMac MacAddress
	SourceMac      MacAddress
	EtherType      EtherType
	DataIPv4       IPv4Packet
	DataIPv6       IPv6Packet
}

/*
Function that takes in a byte slice and organizes it into
the EthernetIIFrame struct.
*/
func DataToFrame(dataArr []byte) EthernetIIFrame { // TODO: Look at encoding/gib
	var destMac MacAddress
	for i := 0; i < 6; i++ {
		destMac[i] = dataArr[i]
	}
	var srcMac MacAddress
	for i := 0; i < 6; i++ {
		srcMac[i] = dataArr[6+i]
	}
	etherType := EtherType{dataArr[12], dataArr[13]}
	// ipType := DataPacketType(dataArr[14:])
	// ipTypeInt := int(ipType)
	var ipTypeInt int
	if etherType == IPv4 {
		ipTypeInt = 4
	} else if etherType == IPv6 {
		ipTypeInt = 6
	} else {
		ipTypeInt = -1
	}
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
	}
}

/*
Returns an int value of 4 or 6 if the IP packet in the ethernet
frame uses IPv4 or IPv6 respectively, and otherwise returns -1.
*/
func (frame EthernetIIFrame) GetPacketType() int {
	if frame.EtherType == IPv4 {
		return 4
	}
	if frame.EtherType == IPv6 {
		return 6
	}
	return -1
}

/*
Reads data throughout the EthernetIIFrame struct and reconstructs
the order of bytes represented by a string in hex format and
then returns the string.
*/
func (frame EthernetIIFrame) ToHexString() string {
	var hexString string
	hexString += hex.EncodeToString(frame.DestinationMac[:])
	hexString += hex.EncodeToString(frame.SourceMac[:])
	hexString += hex.EncodeToString(frame.EtherType[:])
	if frame.EtherType == IPv4 {
		hexString += frame.DataIPv4.ToHexString()
	}
	if frame.EtherType == IPv6 {
		hexString += frame.DataIPv6.ToHexString()
	}

	return hexString
}

/*
Returns a string that is a more human-readable representation
of the data contained in the EthernetIIFrame struct.
*/
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
		returnStr += "\n\n    IPv4: {"
		returnStr += "\n" + frame.DataIPv4.ToString()
		returnStr += "\n    }"
	} else if frame.EtherType == IPv6 {
		returnStr += " (IPv6)"
		returnStr += "\n\n    IPv6: {"
		returnStr += "\n" + frame.DataIPv6.ToString()
		returnStr += "\n    }"
	}

	returnStr += "\n}\n"
	return returnStr
}

/*
Returns a string with bytes represented by their hex values
in a string with those hex values separated by the supplied
delimiter.
*/
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

/*
Returns a string with bytes represented by their decimal values
in a string with those decimal values separated by the supplied
delimiter.
*/
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
