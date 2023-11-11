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

func buildPacket(sourceIP [4]byte, remoteHostIP [4]byte) packet.IPv4Packet {
	var p packet.IPv4Packet
	p.VerAndHeadLen = byte(5)
	p.DSCPandECN = byte(0)
	//p.Length = [2]byte{0, 0}
	p.Identification = [2]byte{0, 0}
	p.FlagAndFrag = [2]byte{0, 0}
	p.TTL = byte(20) //scanner := bufio.NewScanner(file)
	p.Protocol = byte(0)
	p.Checksum = [2]byte{0, 0}
	p.SourceAddr = sourceIP
	p.DestAddr = remoteHostIP
	dataBytes := []byte("this is a test!!!!!")
	p.Data = append([]byte{9, 9, 9, 9, 9, 9, 9, 9}, dataBytes...)

	// Get total length as type uint16 and manually convert
	//	to [2]byte (technically in big-endian?)
	var length uint16 = uint16(20 + len(p.Data))
	p.Length = [2]byte{byte(length >> 8), byte(length)}

	return p
}

// func buildFrame(gatewayMac packet.MacAddress, dataIPv4 packet.IPv4Packet) packet.EthernetIIFrame {
// 	var f packet.EthernetIIFrame

// 	return f
// }

func main() {
	var args []string = os.Args
	var remoteHostIP [4]byte
	if len(args) == 1 || args[1] == "-h" {
		help()
		return
	} else if len(args) >= 2 && (args[1] == "-r" || args[1] == "--remote-host") {
		remoteHostIP = [4]byte(net.ParseIP(args[2]).To4())
	}

	// Get list of interfaces on machine
	ifaces, err := net.Interfaces() // TODO: Specify interface being used as a command line arg
	var iface net.Interface
	if err != nil {
		panic(err)
	}

	// Assuming the wanted interface is the second interface in the list,
	//	attempts to get the IPv4 address associated with the interface as a
	//	byte array with a length of 4
	iface = ifaces[1]
	ifaceAddrs, _ := iface.Addrs()
	var ifaceIP [4]byte
	if ifaceAddrs[0].(*net.IPNet).IP.To4() != nil {
		ifaceIP = [4]byte(ifaceAddrs[0].(*net.IPNet).IP.To4())
	} else {
		panic("Error: Could not find the IPv4 address of the interface being used.")
	}

	//Get IP potential address(es) of default gateway
	addresses, _ := net.LookupHost("_gateway")
	fmt.Println(addresses)

	// Janky, but my hope is that if theres no gateway, the
	//	frame will just be broadcast and reach it's destination
	var gatewayMacStr string = "ff:ff:ff:ff:ff:ff"

	// Use the addresses found for the gateway to find a corresponding
	//	mac address for the gateway in /proc/net/arp that will replace
	//	the broadcast address if found
	for i := 0; i < len(addresses); i++ {
		file, err := os.Open("/proc/net/arp")
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			splitLine := strings.Fields(scanner.Text())
			if addresses[i] == splitLine[0] {
				gatewayMacStr = splitLine[3]
				i = len(addresses)
				break
			}
		}
	}

	//Convert mac string into mac address
	netMac, _ := net.ParseMAC(gatewayMacStr)
	var gatewayMac packet.MacAddress = *((*packet.MacAddress)(netMac))
	fmt.Println(gatewayMac)

	//Build packet to send with command
	var p4 packet.IPv4Packet = buildPacket(ifaceIP, remoteHostIP) //TODO: Implement IPv6
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
