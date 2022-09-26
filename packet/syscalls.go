package packet

import (
	"syscall"
)

/*
Returns an int representing the socket created, a string containing the name
of the bound interface, and a byte array of length 6 containing the hardware
(MAC) address of the bound interface.
*/
func CreateSocket(iface string) (int, error) {

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0x0300)
	if err != nil {
		return -1, err
	}

	err = syscall.BindToDevice(fd, iface)
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
func SendFromSocket(fd int, ifIndex int, frame EthernetIIFrame) {
	var addr syscall.SockaddrLinklayer
	addr.Protocol = syscall.ETH_P_ALL

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
