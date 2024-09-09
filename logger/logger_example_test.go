/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package logger_test

import (
	"github.com/commcos/utils/logger"
)

func ExampleLog() {
	logger.SetLogLevel(logger.TraceLevel)
	logger.Log(logger.InfoLevel, "url", "test")

	// Output:
	//
}
