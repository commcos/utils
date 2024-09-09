/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package factory

import (
	"github.com/commcos/utils/logapi"
	loggerconf "github.com/commcos/utils/logapi/loggerbackend"
	backendlogrus "github.com/commcos/utils/logapi/logrus"
)

func newLogrusbackend(conf *loggerconf.Config) logapi.Interface {
	return backendlogrus.NewLogger(conf.CallerDepth)
}
