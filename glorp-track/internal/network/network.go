package network

import (
	"net"
)

func GetIP(ifaceName string) (net.IP, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil { return nil, err }

	addrs, err := iface.Addrs()
	if err != nil { return nil, err }

	if len(addrs) == 0 { return nil, nil }

	ipv4Addr, _, err := net.ParseCIDR(addrs[0].String())
	if err != nil { return nil, err } 

	return ipv4Addr, nil
}

func IfaceExists(ifaceName string) bool {
	iface, _ := net.InterfaceByName(ifaceName)
	if iface != nil { return true }
	return false
}
