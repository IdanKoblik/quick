package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"

	"quick/internal/logging"
	"quick/internal/networking"
	"quick/pkg/types"
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

	PURPLE         = "\033[38;2;87;87;232m"
	RESET          = "\033[0m"
	ALT_SCREEN_ON  = "\033[?1049h"
	ALT_SCREEN_OFF = "\033[?1049l"
)

var BUILD_TIME string
var VERSION string

func main() {
	closer, err := logging.SetupLogger()
	if err != nil {
		panic(err)
	}

	defer closer.Close()
	printBanner()

	interfaces, err := net.Interfaces()
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	logging.Log.Debug("Finding system network interfaces")
	ifaces := networking.GetInterfaces(&interfaces)

	var selectedIface, selectedIP string
	var selectedMode types.ConnMode

	ifaceOptions := make([]huh.Option[string], 0, len(ifaces))
	for iface := range ifaces {
		ifaceOptions = append(ifaceOptions, huh.NewOption(iface, iface))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a network interface").
				Options(ifaceOptions...).
				Value(&selectedIface),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an IP address").
				OptionsFunc(func() []huh.Option[string] {
					addrs := ifaces[selectedIface]
					opts := make([]huh.Option[string], 0, len(addrs))
					for _, addr := range addrs {
						opts = append(opts, huh.NewOption(addr, addr))
					}
					return opts
				}, &selectedIface).
				Value(&selectedIP),
		),
		huh.NewGroup(
			huh.NewSelect[types.ConnMode]().
				Title("Select a connection mode").
				Options(
					huh.NewOption("Direct connection", types.DIRECT),
					huh.NewOption("P2P", types.P2P),
				).
				Value(&selectedMode),
		),
	)

	fmt.Print(ALT_SCREEN_ON)
	err = form.Run()
	fmt.Print(ALT_SCREEN_OFF)
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	logging.Log.Infof("Selected interface: %s", selectedIface)
	logging.Log.Infof("Selected IP: %s", selectedIP)
	logging.Log.Infof("Selected mode: %s", selectedMode.String())

	peer, err := networking.GenerateIdentity(fmt.Sprintf("%s:%d", selectedIP, networking.Port), selectedMode)
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	logging.Log.Infof("Your address: %s", peer.Addr)
	if peer.Code != "" {
		logging.Log.Infof("Your code: %s", peer.Code)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	switch selectedMode {
	case types.DIRECT:
		err = runDirect(ctx, peer)
	default:
		err = networking.StartServer(ctx, peer)
	}
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}
}

func runDirect(ctx context.Context, peer *types.Identity) error {
	go func() {
		if err := networking.StartServer(ctx, peer); err != nil {
			logging.Log.Errorf("server: %v", err)
		}
	}()

	var peerAddr string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Peer address (ip:port)").
				Description("Leave empty to wait for an incoming connection").
				Value(&peerAddr),
		),
	)

	fmt.Print(ALT_SCREEN_ON)
	err := form.Run()
	fmt.Print(ALT_SCREEN_OFF)
	if err != nil {
		return err
	}

	if strings.TrimSpace(peerAddr) == "" {
		logging.Log.Info("waiting for incoming connection, press Ctrl+C to quit")
		<-ctx.Done()
		return nil
	}

	return networking.Connect(ctx, peer, strings.TrimSpace(peerAddr))
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
