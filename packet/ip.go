package packet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	IPv4Version byte = 0x04
	IPv6Version byte = 0x06
	TCPProt     byte = 0x06
	UDPProt     byte = 0x11
	ICMPProt    byte = 0x01
)

// var IPv4Version byte = 0x04
// var IPv6Version byte  = 0x06
// var TCPProt byte = 0x06
// var UDPProt byte = 0x11
// var ICMPProt byte = 0x01

type IPv4Packet struct {
	VerAndHeadLen  byte // First 4 bits Version, last 4 IHL
	DSCPandECN     byte // First 6 bits DSCP, last 2 ECN
	Length         [2]byte
	Identification [2]byte
	FlagAndFrag    [2]byte
	TTL            byte
	Protocol       byte
	Checksum       [2]byte
	SourceIP       [4]byte
	DestIP         [4]byte
	Data           []byte
}

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

func DataPacketType(dataArr []byte) byte {
	if insertHexFormat([]byte{dataArr[0]}, "")[0] == '4' {
		return 0x04
	} else if insertHexFormat([]byte{dataArr[0]}, "")[0] == '6' {
		return 0x06
	}

	return 0x00
}

func DataToIPv4Packet(dataArr []byte) IPv4Packet { // TODO: Look at encoding/gib
	verAndHeadLen := dataArr[0]
	dSCPandECN := dataArr[1]
	length := [2]byte{dataArr[2], dataArr[3]}
	identification := [2]byte{dataArr[4], dataArr[5]}
	flagAndFrag := [2]byte{dataArr[6], dataArr[7]}
	tTL := dataArr[8]
	proto := dataArr[9]
	check := [2]byte{dataArr[10], dataArr[11]}
	sourceIP := [4]byte{dataArr[12], dataArr[13], dataArr[14], dataArr[15]}
	destIP := [4]byte{dataArr[16], dataArr[17], dataArr[18], dataArr[19]}

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
		SourceIP:       sourceIP,
		DestIP:         destIP,
		Data:           data,
	}
}

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
	//fmt.Println(dataSize)
	// fmt.Println("IPv6")

	// fmt.Println(dataArr[40 : dataSize+40])
	// fmt.Println(dataArr[40:])
	// fmt.Printf("%v\n", payloadLength)
	// fmt.Printf("%d\n", dataSize)
	// fmt.Println()

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

func (packet IPv4Packet) IPv4ToHexString() string {
	var hexString string
	hexString += hex.EncodeToString([]byte{packet.VerAndHeadLen, packet.DSCPandECN})
	hexString += hex.EncodeToString(packet.Length[:])
	hexString += hex.EncodeToString(packet.Identification[:])
	hexString += hex.EncodeToString(packet.FlagAndFrag[:])
	hexString += hex.EncodeToString([]byte{packet.TTL, packet.Protocol})
	hexString += hex.EncodeToString(packet.Checksum[:])
	hexString += hex.EncodeToString(packet.SourceIP[:])
	hexString += hex.EncodeToString(packet.DestIP[:])
	hexString += hex.EncodeToString(packet.Data)

	return hexString
}

func (packet IPv6Packet) IPv6ToHexString() string {
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

func (packet IPv4Packet) ToString() string {
	var returnStr string

	returnStr += "        VerAndHeadLength: "
	returnStr += insertHexFormat([]byte{packet.VerAndHeadLen}, "")

	returnStr += "\n        Source IP: "
	returnStr += insertDecimalFormat(packet.SourceIP[:], ".")

	returnStr += "\n        Destination IP: "
	returnStr += insertDecimalFormat(packet.DestIP[:], ".")

	return returnStr
}
