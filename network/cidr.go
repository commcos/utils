package network

import (
	"math/big"
	"net"
)

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}

// Converts IP mask to 16 bit int.
func MasktoInt(mask net.IPMask) int64 {
	maskInt := big.NewInt(0)
	maskInt.SetBytes(mask)
	return maskInt.Int64()
}

func CheckCidrContain(cidr1, cidr2 string) (contain bool, err error) {
	_, ipv4Net1, err := net.ParseCIDR(cidr1)
	if err != nil {
		return false, err
	}
	_, ipv4Net2, err := net.ParseCIDR(cidr2)
	if err != nil {
		return false, err
	}
	mask := MasktoInt(ipv4Net1.Mask) & MasktoInt(ipv4Net2.Mask)
	ip1Int := IP4toInt(ipv4Net1.IP)
	ip2Int := IP4toInt(ipv4Net2.IP)
	return (ip1Int&mask == ip2Int&mask), nil
}
