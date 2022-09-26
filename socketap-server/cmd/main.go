package main

import (
	"fmt"
	"net"

	"github.com/StaticV0yd/socketap/packet"
)

func main() {

	ifaces, err := net.Interfaces()
	var iface net.Interface
	if err == nil {
		iface = ifaces[1]
	}

	fmt.Println("Attempting to bind to interface", iface.Name+"...")
	fmt.Println(iface.HardwareAddr.String())

	fd, err := packet.CreateSocket(iface.Name)
	if err != nil {
		panic(err)
	}

	defer packet.CloseSocket(fd)
	frame := packet.EthernetIIFrame{}
	//go pwnboard.RecurringUpdate()
	for {
		frame = packet.ReadFromSocket(fd)

		fmt.Println(frame.ToString())
		fmt.Println(frame.ToHexString())

		if frame.DataIPv4.Protocol == byte(0x01) && frame.DataIPv4.SourceIP == [4]byte{8, 8, 8, 8} {
			break
		}
	}

}
