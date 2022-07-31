/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package network

import (
	"io/ioutil"
	"net/http"
	"strings"
)

//GetExternalIPFromIPSB get external ip to ip.sb
func GetExternalIPFromIPSB() string {
	resp, err := http.Get("https://api.ip.sb/ip")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	ip := strings.TrimRight(string(content), "\n")
	return ip
}
