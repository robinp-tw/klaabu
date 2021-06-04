package iputil

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

// CompareIPs Useful for sorting. CompareIPs(10.0.0.0/8, 10.0.0.1/8) == -1
func CompareIPs(x, y net.IP) (int, error) {
	if len(x) != len(y) {
		return 0, fmt.Errorf("length of IP %v is different from %v , meaning two invalid IPs", x, y)
	}

	for i := 0; i < len(x); i++ {
		if x[i] < y[i] {
			return -1, nil
		} else if x[i] > y[i] {
			return 1, nil
		}
	}

	return 0, nil
}

// MinMaxIP first and last IP of the CIDR range
func MinMaxIP(cidr string) (net.IP, net.IP, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, fmt.Errorf("error while parsing your CIDR %v with error: %s", cidr, err)
	}

	min := NetworkIP(ipNet)
	max := BroadcastIP(ipNet)

	return min, max, nil
}

func NetworkIP(ipNet *net.IPNet) net.IP {
	res := make([]byte, len(ipNet.IP))
	for i := range ipNet.IP {
		res[i] = ipNet.Mask[i] & ipNet.IP[i]
	}
	return res
}

func BroadcastIP(ipNet *net.IPNet) net.IP {
	res := make([]byte, len(ipNet.IP))
	for i := range ipNet.IP {
		res[i] = ipNet.Mask[i]&ipNet.IP[i] | ^ipNet.Mask[i]
	}
	return res
}

func NextIP(ip net.IP) (net.IP, error) {
	result := CloneIP(ip)

	for i := len(result) - 1; i >= 0; i-- {
		if i == 0 && result[i] == 255 {
			return nil, errors.New("unable to increment max IP")
		}

		if result[i] < 255 {
			result[i]++
			break
		} else {
			result[i] = 0
		}
	}

	return result, nil
}

func PreviousIP(ip net.IP) (net.IP, error) {
	result := CloneIP(ip)

	for i := len(result) - 1; i >= 0; i-- {
		if i == 0 && result[i] == 0 {
			return nil, errors.New("unable to decrement min IP")
		}

		if result[i] > 0 {
			// no underflow, terminate
			result[i]--
			break
		} else {
			result[i] = 255
		}
	}

	return result, nil
}

func CloneIP(ip net.IP) net.IP {
	result := net.IP(make([]byte, len(ip)))
	for index, b := range ip {
		result[index] = b
	}

	return result
}

// NetContainsNet is true if a contains b completely.
func NetContainsNet(a, b *net.IPNet) bool {
	// Based on https://github.com/google/ipaddr-py/blob/master/ipaddr.py#L648
	return bytes.Compare(NetworkIP(a), NetworkIP(b)) <= 0 && bytes.Compare(BroadcastIP(a), BroadcastIP(b)) >= 0
}
