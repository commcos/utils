/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package factory

import (
	"github.com/xsbull/utils/logapi"
	loggerconf "github.com/xsbull/utils/logapi/loggerbackend"
)

func CreateLogHandle(conf *loggerconf.Config) logapi.Interface {
	return newLogrusbackend(conf)
}
