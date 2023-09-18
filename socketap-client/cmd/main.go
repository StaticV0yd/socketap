package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/StaticV0yd/socketap/packet"
)

func help() {
	fmt.Println("Usage: socketap-client [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	fmt.Println("  -h, --help          show this help message")
	fmt.Println("  -r, --remote-host   the ip of the remote host running socketap-server")
}

func buildPacket(gatewayIP string) packet.IPv4Packet {
	var p packet.IPv4Packet

	return p
}

func buildFrame(gatewayMac packet.MacAddress, dataIPv4 packet.IPv4Packet) packet.EthernetIIFrame {
	var f packet.EthernetIIFrame

	return f
}

func main() {
	var args []string = os.Args
	if len(args) == 1 || args[1] == "-h" {
		help()
		return
	}

	ifaces, err := net.Interfaces()
	var iface net.Interface
	if err != nil {
		panic(err)
	}

	iface = ifaces[1]

	//Get IP address of the default gateway
	addresses, _ := net.LookupHost("_gateway")
	fmt.Println(addresses)
	gatewayIP := addresses[0]

	//Look for entry in arp table using gateway IP to get
	//	gateway mac address
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		panic(err)
	}

	// Janky, but my hope is that if theres no gateway, the
	//	frame will just be broadcast and reach it's destination
	var gatewayMacStr string = "ff:ff:ff:ff:ff:ff"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitLine := strings.Fields(scanner.Text())
		if gatewayIP == splitLine[0] {
			fmt.Println(splitLine)
			gatewayMacStr = splitLine[3]
		}
	}
	fmt.Println(gatewayMacStr)

	//Convert mac string into mac address
	netMac, _ := net.ParseMAC(gatewayMacStr)
	var gatewayMac packet.MacAddress = *((*packet.MacAddress)(netMac))
	fmt.Println(gatewayMac)

	//Build packet to send with command
	var p4 packet.IPv4Packet = buildPacket(gatewayIP) //TODO: Implement IPv6
	var p6 packet.IPv6Packet

	//Build frame to send
	var f packet.EthernetIIFrame = packet.CreateFrame(gatewayMac, iface, packet.IPv4, p4, p6)

	//Open raw socket

	fd, err := packet.CreateSocket(iface)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	defer packet.CloseSocket(fd)

	//Send frame
	err = packet.SendFromSocket(fd, iface.Index, f)
	if err != nil {
		fmt.Println("Error sending Ethernet frame:", err)
	} else {
		fmt.Println("Ethernet frame sent successfully.")
	}
}
