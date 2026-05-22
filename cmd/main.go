package main

import (
	"fmt"
	"net"
	"crypto/rand"
	"math/big"

	"quick/internal/logging"
	"quick/internal/networking"
)

const BANNER = `
		 ██████╗ ██╗   ██╗██╗ ██████╗██╗  ██╗
		██╔═══██╗██║   ██║██║██╔════╝██║ ██╔╝
		██║   ██║██║   ██║██║██║     █████╔╝
		██║▄▄ ██║██║   ██║██║██║     ██╔═██╗
		╚██████╔╝╚██████╔╝██║╚██████╗██║  ██╗
 	 	 ╚══▀▀═╝  ╚═════╝ ╚═╝ ╚═════╝╚═╝  ╚═╝
`

const (
	MAINTAINER = "Idan Koblik"

	PURPLE = "\033[38;2;87;87;232m"
	RESET = "\033[0m"
)

var BUILD_TIME string
var VERSION string

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

func main() {
	fmt.Println(GenerateRandomString(8))
	closer, err := logging.SetupLogger()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	printBanner()

	interfaces, err := net.Interfaces()
	if err != nil {
		logging.Log.Error(err)
		panic(err)
	}

	logging.Log.Debug("Finding system network interfaces")
	ifaces := networking.GetInterfaces(&interfaces)
	for iface, ips := range ifaces {
		str := fmt.Sprintf("Interface: %s [", iface)

		for i, addr := range ips {
			str += addr
			if i != len(ips)-1 {
				str += " ,"
			}
		}

		str += "]"
		logging.Log.Debug(str)
	}
}

func printBanner() {
	fmt.Print(PURPLE)
	fmt.Print(BANNER)
	fmt.Print(RESET)
	fmt.Println()

	buildTime := BUILD_TIME
	if buildTime == "" {
		buildTime = "unknown"
	}

	fmt.Printf("\t\t%s • %s • %s\n\n",MAINTAINER, VERSION, buildTime)
}
