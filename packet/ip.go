package packet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

/*
Constants that store various byte values of significance.
*/
const (
	IPv4Version byte = 0x04
	IPv6Version byte = 0x06
	TCPProt     byte = 0x06
	UDPProt     byte = 0x11
	ICMPProt    byte = 0x01
)

/*
Struct that helps describe IPv4 packet structure for both
looking at IPv4 packets and constructing IPv4 packets.
*/
type IPv4Packet struct {
	VerAndHeadLen  byte // First 4 bits Version, last 4 IHL
	DSCPandECN     byte // First 6 bits DSCP, last 2 ECN
	Length         [2]byte
	Identification [2]byte
	FlagAndFrag    [2]byte
	TTL            byte
	Protocol       byte
	Checksum       [2]byte
	SourceAddr     [4]byte
	DestAddr       [4]byte
	Data           []byte
}

/*
Struct that helps describe IPv6 packet structure for both
looking at IPv6 packets and constucting IPv6 packets.
*/
type IPv6Packet struct {
	Version       byte // Actually half a byte (4 bits 0110 = 6 meaning IPv6)
	TrafficClass  byte
	FlowLabel     [3]byte // Actually 2 and a half bytes
	PayloadLength [2]byte
	NextHeader    byte
	HopLimit      byte
	SourceAddr    [16]byte
	DestAddr      [16]byte
	Data          []byte
}

/*
Creates, populates, and returns an instance of the IPv4Packet
struct based on the slice of bytes passed into the function.
*/
func DataToIPv4Packet(dataArr []byte) IPv4Packet { // TODO: Look at encoding/gib
	verAndHeadLen := dataArr[0]
	dSCPandECN := dataArr[1]
	length := [2]byte{dataArr[2], dataArr[3]}
	identification := [2]byte{dataArr[4], dataArr[5]}
	flagAndFrag := [2]byte{dataArr[6], dataArr[7]}
	tTL := dataArr[8]
	proto := dataArr[9]
	check := [2]byte{dataArr[10], dataArr[11]}
	sourceAddr := [4]byte{dataArr[12], dataArr[13], dataArr[14], dataArr[15]}
	destAddr := [4]byte{dataArr[16], dataArr[17], dataArr[18], dataArr[19]}

	dataSize := binary.BigEndian.Uint16(length[:])
	// fmt.Printf("%v\n", length)
	// fmt.Printf("%d\n", dataSize)
	// fmt.Println()
	data := dataArr[20:dataSize]

	// fmt.Println("IPv4")
	// fmt.Println(dataArr[20:dataSize])
	// fmt.Println(dataArr[20:])

	return IPv4Packet{
		VerAndHeadLen:  verAndHeadLen,
		DSCPandECN:     dSCPandECN,
		Length:         length,
		Identification: identification,
		FlagAndFrag:    flagAndFrag,
		TTL:            tTL,
		Protocol:       proto,
		Checksum:       check,
		SourceAddr:     sourceAddr,
		DestAddr:       destAddr,
		Data:           data,
	}
}

/*
Creates, populates, and returns an instance of the IPv6Packet
struct based on the slice of bytes passed into the function.
*/
func DataToIPv6Packet(dataArr []byte) IPv6Packet {
	version := byte(0x06)
	tempStr := ""
	for _, n := range dataArr[0:2] {
		tempStr += fmt.Sprintf("%0.8b", n)
	}
	tempStr = tempStr[4:]
	tempStr2 := "0000" + tempStr[len(tempStr)-4:]
	tempStr = tempStr[:len(tempStr)-4]

	tempInt := 0
	tempInt2 := 0
	for i := 0; i < 8; i++ {
		if tempStr[i] == '1' {
			tempInt += 2 ^ (7 - i)
		}
		if tempStr2[i] == '1' {
			tempInt2 += 2 ^ (7 - i)
		}
	}

	trafficClass := byte(tempInt)
	flowLabel := [3]byte{byte(tempInt2), dataArr[2], dataArr[3]}
	payloadLength := [2]byte{dataArr[4], dataArr[5]}
	nextHeader := dataArr[6]
	hopLimit := dataArr[7]
	var sourceAddr [16]byte
	copy(sourceAddr[:], dataArr[8:24])
	var destAddr [16]byte
	copy(destAddr[:], dataArr[24:40])

	dataSize := binary.BigEndian.Uint16(payloadLength[:])

	data := dataArr[40 : dataSize+40]

	return IPv6Packet{
		Version:       version,
		TrafficClass:  trafficClass,
		FlowLabel:     flowLabel,
		PayloadLength: payloadLength,
		NextHeader:    nextHeader,
		HopLimit:      hopLimit,
		SourceAddr:    sourceAddr,
		DestAddr:      destAddr,
		Data:          data,
	}
}

/*
Returns a string containing the byte representation of the data
in the IPv4Packet struct instance.
*/
func (packet IPv4Packet) ToHexString() string {
	var hexString string
	hexString += hex.EncodeToString([]byte{packet.VerAndHeadLen, packet.DSCPandECN})
	hexString += hex.EncodeToString(packet.Length[:])
	hexString += hex.EncodeToString(packet.Identification[:])
	hexString += hex.EncodeToString(packet.FlagAndFrag[:])
	hexString += hex.EncodeToString([]byte{packet.TTL, packet.Protocol})
	hexString += hex.EncodeToString(packet.Checksum[:])
	hexString += hex.EncodeToString(packet.SourceAddr[:])
	hexString += hex.EncodeToString(packet.DestAddr[:])
	hexString += hex.EncodeToString(packet.Data)

	return hexString
}

/*
Returns a string containing the byte representation of the data
in the IPv6Packet struct instance.
*/
func (packet IPv6Packet) ToHexString() string {
	var hexString string
	hexString += hex.EncodeToString([]byte{packet.Version})[2:]
	hexString += hex.EncodeToString([]byte{packet.TrafficClass})
	hexString += hex.EncodeToString(packet.FlowLabel[:])[2:]
	hexString += hex.EncodeToString(packet.PayloadLength[:])
	hexString += hex.EncodeToString([]byte{packet.NextHeader, packet.HopLimit})
	hexString += hex.EncodeToString(packet.SourceAddr[:])
	hexString += hex.EncodeToString(packet.DestAddr[:])
	hexString += hex.EncodeToString(packet.Data)

	return hexString
}

/*
Returns a string that is a more human-readable representation
of the data in the IPv4Packet struct instance.
*/
func (packet IPv4Packet) ToString() string {
	var returnStr string

	returnStr += "        VerAndHeadLength: 0x"
	returnStr += insertHexFormat([]byte{packet.VerAndHeadLen}, "")

	returnStr += "\n        DSCP and ECN: 0x"
	returnStr += insertHexFormat([]byte{packet.DSCPandECN}, "")

	returnStr += "\n        Total Length: "
	returnStr += fmt.Sprintf("%d", binary.BigEndian.Uint16(packet.Length[:]))

	returnStr += "\n        Identification: 0x"
	returnStr += insertHexFormat(packet.Identification[:], "")

	returnStr += "\n        Flags: 0x"
	returnStr += insertHexFormat(packet.FlagAndFrag[:], "")

	returnStr += "\n        Time to Live: "
	returnStr += fmt.Sprintf("%d", int(packet.TTL))

	returnStr += "\n        Protocol: "
	returnStr += fmt.Sprintf("%d", int(packet.Protocol))

	returnStr += "\n        Checksum: 0x"
	returnStr += insertHexFormat(packet.Checksum[:], "")

	returnStr += "\n        Source IP: "
	returnStr += insertDecimalFormat(packet.SourceAddr[:], ".")

	returnStr += "\n        Destination IP: "
	returnStr += insertDecimalFormat(packet.DestAddr[:], ".")

	return returnStr
}

/*
Returns a string that is a more human-readable representation of the
data in the IPv6Packet struct instance.
*/
func (packet IPv6Packet) ToString() string {
	var returnStr string

	returnStr += "        Version: "
	returnStr += fmt.Sprintf("%d", int(packet.Version))

	returnStr += "\n        Traffic Class: 0x"
	returnStr += insertHexFormat([]byte{packet.TrafficClass}, "")

	returnStr += "\n        Flow Label: 0x"
	returnStr += insertHexFormat(packet.FlowLabel[:], "")

	returnStr += "\n        Payload Length: "
	returnStr += fmt.Sprintf("%d", binary.BigEndian.Uint16(packet.PayloadLength[:]))

	returnStr += "\n        Next Header: "
	returnStr += fmt.Sprintf("%d", int(packet.NextHeader))

	returnStr += "\n        Hop Limit: "
	returnStr += fmt.Sprintf("%d", int(packet.HopLimit))

	returnStr += "\n        Source Address: "
	srcAddr := insertHexFormat(packet.SourceAddr[:], "")
	for i := len(srcAddr) - 4; i > 0; i -= 4 {
		srcAddr = srcAddr[:i] + ":" + srcAddr[i:]
	}
	returnStr += srcAddr

	returnStr += "\n        Destination Address: "
	destAddr := insertHexFormat(packet.DestAddr[:], "")
	for i := len(destAddr) - 4; i > 0; i -= 4 {
		destAddr = destAddr[:i] + ":" + destAddr[i:]
	}
	returnStr += destAddr

	return returnStr
}

/*
Returns a byte slice of all the data in the IPv4Packet struct.
*/
func (packet IPv4Packet) ToData() []byte {
	var data []byte
	data = append(data, packet.VerAndHeadLen, packet.DSCPandECN)
	data = append(data, packet.Length[:]...)
	data = append(data, packet.Identification[:]...)
	data = append(data, packet.FlagAndFrag[:]...)
	data = append(data, packet.TTL, packet.Protocol)
	data = append(data, packet.Checksum[:]...)
	data = append(data, packet.SourceAddr[:]...)
	data = append(data, packet.DestAddr[:]...)
	data = append(data, packet.Data...)

	return data
}

/*
Returns a byte slice of all the data in the IPv6Packet struct. *In progress*
*/
func (packet IPv6Packet) ToData() []byte { //TODO
	var data []byte

	return data
}
