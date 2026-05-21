package networking

import (
	"net"
)

func GetInterfaces(interfaces *[]net.Interface) map[string][]string {
	ips := make(map[string][]string)
	if interfaces == nil {
		return ips
	}

	for _, iface := range *interfaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ips[iface.Name] = append(ips[iface.Name], addr.String())
		}
	}

	return ips
}
