package main

import (
	"fmt"
	"syscall"

	packet "github.com/StaticV0yd/socketap/socketap-packet"
)

func readFromSocket(fd int) packet.EthernetIIFrame {

	data := make([]byte, 65535)
	//for {
	// typeByte := 0
	// var frame packet.EthernetIIFrame
	syscall.Recvfrom(fd, data, 0)
	// fmt.Println(data)
	// fmt.Println()
	frame := packet.DataToFrame(data)
	// typeByte = frame.GetPacketType()
	// fmt.Println(frame.ToHexString())
	// fmt.Println()
	// fmt.Println(frame.ToString())
	// fmt.Println()
	// fmt.Println(hex.EncodeToString(data))
	//}
	// fmt.Println(frame.ToHexString())
	// fmt.Println()
	// fmt.Println(frame.ToString())

	return frame
}

func sendFromSocket(fd int) {
	var addr syscall.SockaddrLinklayer
	addr.Protocol = syscall.ETH_P_ALL

}

func main() {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0300)
	defer syscall.Close(fd)
	if err != nil {
		panic(err)
	}
	frame := packet.EthernetIIFrame{}
	//go pwnboard.RecurringUpdate()
	for true {
		frame = readFromSocket(fd)
		fmt.Println(frame.ToString())
		fmt.Println(frame.ToHexString())
		if frame.DataIPv4.Protocol == byte(0x01) && frame.DataIPv4.SourceIP == [4]byte{8, 8, 8, 8} {
			break
		}
	}

	fmt.Println(frame.ToString())
	fmt.Println(frame.ToHexString())

}
