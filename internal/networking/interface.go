package networking

import (
	"net"
	"strings"

	"quick/internal/logging"
)

func GetInterfaces(interfaces *[]net.Interface) map[string][]string {
	ips := make(map[string][]string)
	if interfaces == nil {
		return ips
	}

	for _, iface := range *interfaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ips[iface.Name] = append(ips[iface.Name], ipNet.IP.String())
		}
	}

	for iface, addrs := range ips {
		logging.Log.Debugf("interface: %s ips[%s]", iface, strings.Join(addrs, ", "))
	}

	return ips
}
