package network

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	MaskRegex = "^(254|252|248|240|224|192|128|0)\\.0\\.0\\.0|255\\.(254|252|248|240|224|192|128|0)\\.0\\.0|255\\.255\\.(254|252|248|240|224|192|128|0)\\.0|255\\.255\\.255\\.(255|254|252|248|240|224|192|128|0)$"
)

func ValidateIPAndMask(ipAndMask string) error {
	ipAndMaskInfo := strings.Split(ipAndMask, "/")
	if len(ipAndMaskInfo) != 2 {
		return fmt.Errorf("fail to parse identify match")
	}
	currIP := ipAndMaskInfo[0]
	ipv4Addr := net.ParseIP(currIP)
	if ipv4Addr == nil {
		return fmt.Errorf("ip:%s wrong format", ipAndMask)
	}
	currMask := ipAndMaskInfo[1]
	maskRegex, err := regexp.Compile(MaskRegex)
	if err != nil {
		return fmt.Errorf("create mask:%s regex err:%s", MaskRegex, err)
	}
	if !maskRegex.MatchString(currMask) {
		_, _, err = net.ParseCIDR(ipAndMask)
		if err != nil {
			return fmt.Errorf("parse CIDR err:%s", err)
		}
	}
	return nil
}
