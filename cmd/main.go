package main

import (
	"fmt"
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

func main() {
	printBanner()
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

	fmt.Printf("\t%s • %s • %s\n\n",MAINTAINER, VERSION, buildTime)
}
