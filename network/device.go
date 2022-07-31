/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package network

import (
	"runtime"
	"strings"

	"github.com/google/gopacket/pcap"
)

var deviceAnySupported = runtime.GOOS == "linux"

//DeviceInfo device info list all network interface on host
type DeviceInfo struct {
	DeviceName        string
	DeviceDescription string
	IPAddr            []string
}

// ListDeviceNames returns the list of adapters available for sniffing on
// this computer. If the withDescription parameter is set to true, a human
// readable version of the adapter name is added. If the withIP parameter
// is set to true, IP address of the adapter is added.
func ListDeviceNames() (devList []DeviceInfo, err error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return
	}

	for _, dev := range devices {
		if strings.HasPrefix(dev.Name, "bluetooth") ||
			strings.HasPrefix(dev.Name, "dbus") ||
			strings.HasPrefix(dev.Name, "nf") || strings.HasPrefix(dev.Name, "usb") {
			continue
		}

		item := &DeviceInfo{
			DeviceName: dev.Name,
		}

		if len(dev.Description) > 0 {
			item.DeviceDescription = dev.Description
		}

		if len(dev.Addresses) > 0 {

			for _, address := range []pcap.InterfaceAddress(dev.Addresses) {
				item.IPAddr = append(item.IPAddr, address.IP.String())
			}
		}
		devList = append(devList, *item)
	}
	return
}
