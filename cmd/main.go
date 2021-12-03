package main

import (
	"fmt"
	"syscall"

	"github.com/sockon-script/packet"
)

func readFromSocket() {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0300)
	if err != nil {
		syscall.Close(fd)
		panic(err)
	}
	defer syscall.Close(fd)
	err = syscall.BindToDevice(fd, "wlp115s0")
	if err != nil {
		panic(err)
	}
	data := make([]byte, 1024)
	//for {
	syscall.Recvfrom(fd, data, 0)
	// fmt.Println(data)
	// fmt.Println()
	frame := packet.DataToFrame(data)
	fmt.Println(frame.ToHexString())
	// fmt.Println()
	// fmt.Println(hex.EncodeToString(data))
	//}
}

func main() {
	//fmt.Println(packet.IPv4)
	//test := packet.EthernetFrameFromString()
	// fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Obtained fd ", fd)
	// defer syscall.Close(fd)
	readFromSocket()
}
