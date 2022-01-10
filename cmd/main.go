package main

import (
	"fmt"
	"syscall"

	"github.com/sockon-script/packet"
)

func readFromSocket(fd int) {

	data := make([]byte, 65535)
	//for {
	typeByte := 0
	var frame packet.EthernetIIFrame
	for typeByte != 4 {
		syscall.Recvfrom(fd, data, 0)
		// fmt.Println(data)
		// fmt.Println()
		frame = packet.DataToFrame(data)
		typeByte = frame.GetPacketType()
		// fmt.Println(frame.ToHexString())
		// fmt.Println()
		// fmt.Println(frame.ToString())
		// fmt.Println()
		// fmt.Println(hex.EncodeToString(data))
		//}
	}
	fmt.Println(frame.ToHexString())
	fmt.Println()
	fmt.Println(frame.ToString())
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
	err = syscall.BindToDevice(fd, "wlp115s0")
	if err != nil {
		panic(err)
	}
	readFromSocket(fd)

}
