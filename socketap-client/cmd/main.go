package cmd

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("Usage: socketap-client [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	fmt.Println("  -h, --help          show this help message")
	fmt.Println("  -r, --remote-host   the ip of the remote host running socketap-server")

}

func main() {
	var args []string = os.Args
	if len(args) == 1 {
		help()
	}

}
