// Package nethelp provides assistance to network related tasks
package nethelp

import "net"

// GetLocalAddrs returns all available local addresses.
func GetLocalAddrs() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var list []net.IP

	for _, addr := range addrs {
		v := addr.(*net.IPNet)
		if v.IP.To4() != nil {
			list = append(list, v.IP)
		}
	}

	return list, nil
}
