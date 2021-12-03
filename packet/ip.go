package packet

import "encoding/hex"

var IPv4Version [4]bool = [4]bool{false, true, false, false}
var TCPProt byte = 0x06
var UDPProt byte = 0x11
var ICMPProt byte = 0x01

type IPv4Packet struct {
	VerAndHeadLen  byte
	DSCPandECN     byte
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

func DataToPacket(dataArr []byte) IPv4Packet { // TODO: Look at encoding/gib
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
	data := dataArr[20:]

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

func (packet IPv4Packet) ToHexString() string {
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

// type IPv6Packet struct {
// 	Version       [4]bool
// 	TrafficClass  [1]byte
// 	FlowLabel     [20]bool
// 	PayloadLength [2]byte
// 	NextHeader    [1]byte
// 	HopLimit      [1]byte
// 	SourceAddr    [16]byte
// 	DestAddr      [16]byte
// }
