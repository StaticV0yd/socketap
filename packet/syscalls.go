package packet

import (
	"net"
	"syscall"
)

/*
Returns an int representing the socket created, a string containing the name
of the bound interface, and a byte array of length 6 containing the hardware
(MAC) address of the bound interface.
*/
func CreateSocket(iface net.Interface) (int, error) {

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0300)
	if err != nil {
		return -1, err
	}

	err = syscall.BindToDevice(fd, iface.Name)
	if err != nil {
		return -1, err
	}

	return fd, nil
}

/*
Closes a socket.
*/
func CloseSocket(fd int) {
	syscall.Close(fd)
}

/*
Reads data from a socket and returns the data as an EthernetIIFrame.
*/
func ReadFromSocket(fd int) EthernetIIFrame {

	data := make([]byte, 65535)
	syscall.Recvfrom(fd, data, 0)
	frame := DataToFrame(data)

	return frame
}

/*
Sends data from a socket.
*/
func SendFromSocket(fd int, ifIndex int, frame EthernetIIFrame) error { //TODO
	var addr = &syscall.SockaddrLinklayer{
		Protocol: 0x0300,
		Ifindex:  ifIndex,
	}
	addr.Protocol = syscall.ETH_P_ALL
	frameBytes := make([]byte, 0)
	frameBytes = append(frameBytes, frame.DestinationMac[:]...)
	frameBytes = append(frameBytes, frame.SourceMac[:]...)
	frameBytes = append(frameBytes, frame.EtherType[:]...)
	frameBytes = append(frameBytes, frame.DataIPv4.ToData()...)

	err := syscall.Sendto(fd, frameBytes, 0, addr)

	return err
	// addr.Halen = 6
	// addr.Addr = [8]byte{
	// 	destMac[0],
	// 	destMac[1],
	// 	destMac[2],
	// 	destMac[3],
	// 	destMac[4],
	// 	destMac[5],
	// }

}
